package main

import (
	"os"	
	"log"
	"net"
	"time"
	"context"

	pb "app/proto"
	bot "app/sendTelegram"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedSendMessageServiceServer
}

var High []pb.Request
var Medium []pb.Request
var Low []pb.Request



func (s *server) SendMessage(ctx context.Context, req *pb.Request) (*pb.Response, error) {

	var message pb.Request

	message.Text = req.GetText()
	message.Priority = req.GetPriority()

	if message.Priority == "high" {

		High = append(High, message)
	
	} else if message.Priority == "medium" {

		Medium = append(Medium, message)
	
	} else if message.Priority == "low" {

		Low = append(Low, message)
	
	} else {

		log.Println("Unknown Priority")
	}

	return &pb.Response{}, nil
}

func main() {

	err := godotenv.Load()
    
    if err != nil {
        log.Fatalf("Problem with loading env: %v", err)
    }	

    log.Println("Welcome to the server!")
    
    port := os.Getenv("GRPC_SERVER_PORT")

    listener, err := net.Listen("tcp", port)

    if err != nil {
		log.Fatalf("Failed to connecting tcp port: %v", err)
	}

	go FIFO()

	s := grpc.NewServer()

	pb.RegisterSendMessageServiceServer(s, &server{})

	err = s.Serve(listener)

	if err != nil {
		log.Fatalf("Failed connecting to the Server: %v", err)
	}
}

func FIFO() {

	var sended bool

	for {
		sended = false
		for i := range High {
			err := bot.SendText(High[i].Text)
			if err != nil {
				log.Fatalf("Failed sending to the bot: %v", err)
			} else {
				sended = true
				High = Remove(High, i)
				time.Sleep(5 * time.Second)

				break
			}
		} 

		if sended {
			continue
		}

		for i := range Medium {
			err := bot.SendText(Medium[i].Text)
			if err != nil {
				log.Fatalf("Failed sending to the bot: %v", err)
			} else {
				sended = true
				Medium = Remove(Medium, i)
				time.Sleep(5 * time.Second)

				break
			}
		}

		if sended {
			continue
		}

		for i := range Low {
			err := bot.SendText(Low[i].Text)
			if err != nil {
				log.Fatalf("Failed sending to the bot: %v", err)
			} else {
				sended = true
				Low = Remove(Low, i)
				time.Sleep(5 * time.Second)

				break
			}
		}


	}
}

func Remove(s []pb.Request, i int) []pb.Request {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}