package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/mmkamron/basefit/internal/data"
	"github.com/mmkamron/basefit/internal/mailer"
	"github.com/mmkamron/basefit/internal/pkg/config"
	"github.com/mmkamron/basefit/internal/pkg/db"
)

type application struct {
	config *config.Config
	logger *slog.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {
	cfg := config.Load("./config/local.yaml")
	log := slog.Default()
	db := db.Load(cfg)
	defer db.Close()
	app := &application{
		config: cfg,
		logger: log,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.Smtp.Host, cfg.Smtp.Port, cfg.Smtp.Username, cfg.Smtp.Password, cfg.Smtp.Sender),
	}
	if err := app.serve(); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         app.config.Addr,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	//TODO:find why graceful shutdown is not working
	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.logger.Info(s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		app.logger.Info("completing background tasks...")

		app.wg.Wait()
		shutdownError <- nil
	}()

	app.logger.Info("starting the server")
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Info("stopped server")
	return nil
}
