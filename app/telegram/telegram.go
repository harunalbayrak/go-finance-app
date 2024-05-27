package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendMessage(text string) (bool, error) {
	// Global variables
	var err error
	var response *http.Response

	fmt.Println("ChatID:", os.Getenv("CHAT_ID"))

	// Send the message
	tokenUrl := fmt.Sprintf("https://api.telegram.org/bot%s", os.Getenv("TELEGRAM_TOKEN"))
	url := fmt.Sprintf("%s/sendMessage", tokenUrl)
	body, _ := json.Marshal(map[string]string{
		"chat_id": os.Getenv("CHAT_ID"),
		// "message_thread_id": os.Getenv("MESSAGE_THREAD_ID"),
		"text": text,
		// "parse_mode":        "MarkdownV2",
	})
	response, err = http.Post(
		url,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return false, err
	}

	// Close the request at the end
	defer response.Body.Close()

	// Body
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return false, err
	}

	// Log
	fmt.Printf("Message\n%s\n", text)
	fmt.Println("ResponseJSON", string(body))

	// Return
	return true, nil
}
func SendPhoto(photoPath string) error {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		return err
	}

	chatId, err := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewPhoto(chatId, tgbotapi.FilePath(photoPath))
	msg.Caption = "EKIZ Grafik"
	_, err = bot.Send(msg)
	if err != nil {
		return err
	}

	newMsg := tgbotapi.NewMessage(chatId, "asdas")
	_, err = bot.Send(newMsg)
	if err != nil {
		return err
	}

	return nil
}
