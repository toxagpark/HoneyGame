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
	fight_pg_repo "github.com/toxagpark/HoneyGame/internal/feature/fight/repository/postgres"
	fight_service "github.com/toxagpark/HoneyGame/internal/feature/fight/service"
	fight_tg_transport "github.com/toxagpark/HoneyGame/internal/feature/fight/transport/tg"
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

	fightRepo := fight_pg_repo.NewRepository(pool)
	fightService := fight_service.NewService(fightRepo, service)
	fightTransport := fight_tg_transport.NewHandler(bot, fightService)

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

	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		fightTransport.HandleCreateChallenge(ctx, &update)
		return nil
	}, th.CommandEqual("fight"))

	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		fightTransport.HandleListChallenges(ctx, &update)
		return nil
	}, th.CommandEqual("fights"))

	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		if update.CallbackQuery != nil {
			fightTransport.HandleAcceptCallback(ctx, *update.CallbackQuery)
		}
		return nil
	}, th.CallbackDataPrefix("fight:accept:"))

	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		if update.CallbackQuery != nil {
			fightTransport.HandleCancelCallback(ctx, *update.CallbackQuery)
		}
		return nil
	}, th.CallbackDataPrefix("fight:cancel:"))

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
