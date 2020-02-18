package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type message struct {
	To      int64
	Message string
	Silent  bool
}

func main() {

	appAuthToken := os.Getenv("APP_TOKEN")

	telegramToken := os.Getenv("TELEGRAM_TOKEN")

	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Panic(err)
	}

	if _, ok := os.LookupEnv("DEBUG"); ok {
		log.Println("Debug logging enabled")
		bot.Debug = true
	}

	log.Printf("Authorized on Telegram as %s", bot.Self.UserName)

	http.HandleFunc("/msg", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			var authHeaderValue = r.Header.Get("X-Auth-Token")

			if authHeaderValue != appAuthToken {
				log.Println("Request denied due to missing or invalid auth token")
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintln(w, "401 - Missing or invalid auth token in header X-Auth-Token")
				return
			}

			var m message
			err := json.NewDecoder(r.Body).Decode(&m)

			if err != nil {
				panic(err)
			}

			log.Printf("Message to %d with msg %q, silent: %t\n", m.To, m.Message, m.Silent)

			msg := tgbotapi.NewMessage(m.To, m.Message)
			if m.Silent {
				msg.DisableNotification = true
			}
			bot.Send(msg)

			w.WriteHeader(http.StatusInternalServerError)

			fmt.Fprintf(w, "Sent message %q to %d (silent=%t)\n", m.Message, m.To, m.Silent)

		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
