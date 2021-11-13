package main

import (
	"os"
	"fmt"
	"log"
	"context"

	pb "app/proto"

	"google.golang.org/grpc"
	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"

	_ "app/client/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Message struct {
	Text string `json:"text"`
	Priority string `json:"priority"`
}

var c pb.SendMessageServiceClient

// @Summary MessageSender
// @Description Send Message
// @Tags sender
// @ID message_sender
// @Accept json
// @Produce json
// @Param input body Message true "message info"
// @Success 200
// @Failure 400
// @Router /send [post]
func SendMessageController(ctx *gin.Context) {

	var newMessage pb.Request

	err := ctx.BindJSON(&newMessage) 

	if err != nil {

		log.Fatalf("Failed to accept message: %v", err)
	}

	res, err := c.SendMessage(context.Background(), &pb.Request{
		Text:  newMessage.Text,
		Priority: newMessage.Priority,
	})

	fmt.Println(res)

	if err != nil {

		log.Fatalf("Failed to call Sender RPC: %v", err)
	}
}

// @title Message Sender Bot
// @version 1.0
// @description Telegram Bot which sends messages to channels and groups
// @host localhost:8080
// @BasePath /

func main() {

	err := godotenv.Load()

    if err != nil {
        log.Fatalf("Problem with loading env: %v", err)
    }

	log.Println("Welcome to the client!")

    port := os.Getenv("GRPC_SERVER_PORT")

	conn, err := grpc.Dial(fmt.Sprintf("localhost%s", port), grpc.WithInsecure())
	
	if err != nil {
		log.Printf("Problem with connecting client to server: %s", err)
		return
	}

	defer conn.Close()

	c = pb.NewSendMessageServiceClient(conn)

	http_port := os.Getenv("HTTP_PORT")

	router := gin.Default()

	router.POST("/send", SendMessageController)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	
	log.Println("Listening and serving HTTP on localhost", http_port)

	err = router.Run(http_port)

	if err != nil {
		log.Fatalf("Cannot connecting to the http_port: %v", err)
	}
}