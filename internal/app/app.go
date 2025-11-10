package app

import (
	"ai-service/internal/app/router"
	"ai-service/internal/repository"
	"ai-service/internal/util/config"
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	cfg, err := config.LoadConfig("config/conf.json")
	if err != nil {
		log.Fatalf("error parse config file: %v", err)
		return
	}
	repo, milvus, err := repository.New(ctx, cfg)
	defer milvus.Close()
	if err != nil {
		panic(err)
	}
	api := router.NewRouter(cfg, repo).Build(ctx)

	g, wgCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		err = api.Start(cfg.Port)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
		return nil
	})
	g.Go(func() error {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-wgCtx.Done():
			return wgCtx.Err()
		case sig := <-quit:
			fmt.Printf("received signal: [%+v]", sig)
		}
		return api.Shutdown(wgCtx)
	})
	g.Wait()
}
