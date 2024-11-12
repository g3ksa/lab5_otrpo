package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/g3ksa/lab5_otrpo/internal/graph_api/config"
	"github.com/g3ksa/lab5_otrpo/internal/graph_api/handlers"
	"github.com/g3ksa/lab5_otrpo/internal/graph_api/service/storage"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.NewGraphAPIConfig()
	if err != nil {
		slog.Error("error", err)
	}

	uri := fmt.Sprintf("bolt://%s:%d", cfg.Database.Host, cfg.Database.Port)

	db, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(cfg.Database.User, cfg.Database.Password, ""))
	defer db.Close(ctx)

	dbStorage := storage.NewDBStorage(db)

	httpService := handlers.NewHttpServer(dbStorage, cfg.SecretToken)

	httpServer := &http.Server{Addr: cfg.HTTPAddr, Handler: httpService.Router()}
	go func() {
		slog.Info("HTTP server listening on", slog.String("port", cfg.HTTPAddr))
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("listen", err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-signalChan
	slog.Info("Received signal", sig)

	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		slog.Error("error", err)
	}
	slog.Info("http server shutdown")
	cancel()

	slog.Info("Done")
}
