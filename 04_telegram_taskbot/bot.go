package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	tgbotapi "github.com/skinass/telegram-bot-api/v5"
)

var (
	BotToken = "5976082281:AAGpCxuZ_B3lZIPtfqKnsJ4a2JxWIHa04iI"

	WebhookURL = "https://6f08-178-167-55-173.eu.ngrok.io"

	Port = 8081
)

type Response struct {
	ChatID  int64
	Message string
}

func startTaskBot(ctx context.Context) error {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		return fmt.Errorf("NewBotApi failed: %w", err)
	}

	// bot.Debug = true
	log.Printf("Authorized on account %s\n", bot.Self.UserName)

	wh, err := tgbotapi.NewWebhook(WebhookURL)
	if err != nil {
		return fmt.Errorf("NewWebhook failed: %w", err)
	}

	_, err = bot.Request(wh)
	if err != nil {
		return fmt.Errorf("SetWebhook failed: %w", err)
	}

	updates := bot.ListenForWebhook("/")

	server := &http.Server{Addr: ":" + strconv.Itoa(Port), Handler: nil}
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe failed: %v\n", err)
		}
	}()
	log.Println("Listening on port:", Port)

	handler := TaskHandler{storage: NewTaskStorage()}

	for {
		select {
		case <-ctx.Done():
			updates.Clear()
			log.Println("Server is shutting down")
			if err := server.Shutdown(context.Background()); err != nil {
				return fmt.Errorf("HTTP server shutdown error: %w", err)
			}
			log.Println("Server shut down")
			return nil
		case update := <-updates:
			if update.Message == nil {
				continue
			}
			responses := handler.HandleUpdate(&update)
			err := SendResponses(bot, responses)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func SendResponses(bot *tgbotapi.BotAPI, responses []Response) error {
	for _, r := range responses {
		msg := tgbotapi.NewMessage(r.ChatID, r.Message)
		_, err := bot.Send(msg)
		if err != nil {
			return fmt.Errorf("bot.Send failed: %w", err)
		}
	}
	return nil
}
