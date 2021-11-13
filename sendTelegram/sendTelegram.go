package sendTelegram

import (
    "os"
	"log"
    "strconv"

    "github.com/joho/godotenv"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func SendText(text string) error {

    err := godotenv.Load()
    
    if err != nil {
        log.Fatalf("Problem with loading env: %v", err)
        return err
    }

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
    
    if err != nil {
        log.Fatalf("Cannot connecting to the bot: %v", err)
        return err
    }

    chatId := os.Getenv("CHAT_ID")

    id, err := strconv.ParseInt(chatId, 10, 64)
    if err != nil {
        log.Fatalf("Cannot converting string to int64: %v", err)
    }
    
    bot.Debug = true
    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    msg := tgbotapi.NewMessage(id, text)

    log.Printf("%T, %+v ", msg, msg)

    if _, err := bot.Send(msg); err != nil {
        log.Fatalf("Problem with sending message: %v",err)
        return err
    }

    return err
}