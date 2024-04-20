package gemini

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
	"google.golang.org/api/option"
)

func init() {
	functions.HTTP("gemini", callback)
}

func gemini(s string) string {
	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.0-pro-latest")
	resp, err := model.GenerateContent(ctx, genai.Text(s))
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range resp.Candidates {
		for _, p := range c.Content.Parts {
			return fmt.Sprintf("%v", p)
		}
	}
	return ""
}

func callback(w http.ResponseWriter, r *http.Request) {
	bot, err := messaging_api.NewMessagingApiAPI(
		os.Getenv("LINE_CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	cb, err := webhook.ParseRequest(os.Getenv("LINE_CHANNEL_SECRET"), r)
	if err != nil {
		log.Fatal(err)
	}
	for _, event := range cb.Events {
		switch e := event.(type) {
		case webhook.MessageEvent:
			switch message := e.Message.(type) {
			case webhook.TextMessageContent:
				switch source := e.Source.(type) {
				case webhook.UserSource:
					if _, err = bot.ShowLoadingAnimation(
						&messaging_api.ShowLoadingAnimationRequest{
							ChatId:         source.UserId,
							LoadingSeconds: 30,
						},
					); err != nil {
						log.Print(err)
					}
				}
				if _, err = bot.ReplyMessage(
					&messaging_api.ReplyMessageRequest{
						ReplyToken: e.ReplyToken,
						Messages: []messaging_api.MessageInterface{
							messaging_api.TextMessage{
								Text: gemini(message.Text),
							},
						},
					},
				); err != nil {
					log.Print(err)
				} else {
					log.Println("Sent text reply.")
				}
			}
		}
	}
	w.WriteHeader(200)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/callback", callback)
	http.ListenAndServe(":8080", nil)
}
