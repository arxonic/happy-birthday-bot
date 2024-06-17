package app

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/arxonic/gmh/internal/config"
	emailController "github.com/arxonic/gmh/internal/controllers/email"
	empController "github.com/arxonic/gmh/internal/controllers/employers"
	v1 "github.com/arxonic/gmh/internal/controllers/http/v1"
	"github.com/arxonic/gmh/internal/controllers/telegram"
	"github.com/arxonic/gmh/internal/controllers/telegram/states"
	"github.com/arxonic/gmh/internal/lib/logger/sl"
	"github.com/arxonic/gmh/internal/services/auth"
	"github.com/arxonic/gmh/internal/services/email"
	"github.com/arxonic/gmh/internal/services/employers"
	"github.com/arxonic/gmh/internal/services/subscribe"
	"github.com/arxonic/gmh/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
)

func Run(cfg *config.Config, log *slog.Logger) {
	// data layer
	// -- init storage
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	// -- init external Employer API
	empAPI := empController.New("http://localhost:4242/users/user/")
	// -- init email server
	emailSrv := emailController.New(cfg.MailServer.Host, cfg.MailServer.Port, cfg.MailServer.Sender, cfg.MailServer.Password)

	// use case
	// -- init employers service
	emloyerService := employers.New(log, empAPI)
	// -- init notify service
	notifyService := email.New(log, emailSrv)
	// -- init auth service
	authService := auth.New(log, storage, storage, storage, notifyService)
	// -- init subscribe service
	subService := subscribe.New(log, storage, storage)

	// transport
	httpRouter := chi.NewRouter()
	v1.NewRouts(httpRouter, log, authService)
	srv := v1.NewServer(cfg.Address, httpRouter)
	go func() {
		if err := v1.Run(srv); err != nil {
			log.Error("failed to run http server", sl.Err(err))
			os.Exit(1)
		}
	}()

	states := states.NewStates()
	bot, err := telegram.NewBot(cfg.TgBotKey, log)
	if err != nil {
		log.Error("failed to init telegram bot", sl.Err(err))
		os.Exit(1)
	}

	bot.Run(states, subService, authService, emloyerService)

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	notify := <-stop

	// TODO stop app

	log.Info("application stopped", slog.String("signal", notify.String()))
}
