package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"

	core_postgres "github.com/toxagpark/HoneyGame/internal/core/postgres"
	users_pg_repo "github.com/toxagpark/HoneyGame/internal/feature/users/repository/postgres"
	users_service "github.com/toxagpark/HoneyGame/internal/feature/users/service"
	users_tg_transport "github.com/toxagpark/HoneyGame/internal/feature/users/transport/tg"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	pool, err := core_postgres.NewPool(ctx)
	if err != nil {
		panic(err)
	}
	defer core_postgres.ClosePool(pool)

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		panic("TELEGRAM_BOT_TOKEN is not set")
	}

	bot, err := telego.NewBot(token)
	if err != nil {
		panic(err)
	}

	updates, err := bot.UpdatesViaLongPolling(ctx, nil)
	if err != nil {
		panic(err)
	}

	repo := users_pg_repo.NewRepository(pool)
	service := users_service.NewService(repo)
	transport := users_tg_transport.NewHandler(bot, service)

	bh, err := th.NewBotHandler(bot, updates)
	if err != nil {
		panic(err)
	}

	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		transport.HandleCreateOrUpdateUser(ctx, &update)
		return nil
	}, th.CommandEqual("start"))

	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		transport.HandleGetUser(ctx, &update)
		return nil
	}, th.CommandEqual("getUser"))

	go func() {
		if err := bh.Start(); err != nil {
			log.Printf("bot handler: %v", err)
		}
	}()

	<-ctx.Done()

	if err := bh.Stop(); err != nil {
		log.Printf("bot handler stop: %v", err)
	}
	fmt.Println("Shutting down...")
}
