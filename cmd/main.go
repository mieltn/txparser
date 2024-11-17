package main

import (
	"context"
	"flag"
	"github.com/mieltn/txparser/internal/client"
	"github.com/mieltn/txparser/internal/config"
	"github.com/mieltn/txparser/internal/handlers"
	"github.com/mieltn/txparser/internal/logger"
	"github.com/mieltn/txparser/internal/repositories/inmemory"
	"github.com/mieltn/txparser/internal/router"
	"github.com/mieltn/txparser/internal/server"
	"github.com/mieltn/txparser/internal/services"
	"github.com/mieltn/txparser/pkg/eth"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	mode := flag.String("mode", "", "mode of operation")
	flag.Parse()

	l := logger.New(*mode)

	cfg := config.Config{
		Mode: *mode,
	}
	if err := config.Load(&cfg); err != nil {
		l.Errorf("load config failed: %v", err)
		time.Sleep(time.Second * 1)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	srv, background, err := buildDeps(cfg, l)
	if err != nil {
		l.Errorf("build deps failed: %v", err)
		time.Sleep(time.Second * 1)
		os.Exit(1)
	}

	if background != nil {
		go background(ctx)
	}

	srv.Run(ctx)
	time.Sleep(time.Second * 2)
}

func buildDeps(cfg config.Config, l logger.Logger) (app, func(ctx context.Context), error) {
	// init inmemory repos
	inmemTransactions := inmemory.NewTransactions(l)
	inmemAddresses := inmemory.NewAddresses(l)

	blockStart, err := eth.ParseBigInt256(cfg.App.StartBlock)
	if err != nil {
		return nil, nil, err
	}
	inmemLastProcessed := inmemory.NewProcessedBlockRepo(l, blockStart)

	// client
	rpcClient := client.New(l, cfg.Eth.Url, cfg.Eth.Retry, cfg.Eth.RetryIn, cfg.Eth.Timeout)

	// parser
	serviceTxparser := services.NewTxparser(l, rpcClient, inmemTransactions, inmemAddresses, inmemLastProcessed, cfg.App.PollWorkers)

	// handlers
	handlerTransactions := handlers.NewTransactions(l, serviceTxparser)
	handlerAddresses := handlers.NewAddresses(l, serviceTxparser)

	// router
	router := router.New(
		handlerAddresses,
		handlerTransactions,
	)

	srv := server.New(l, cfg, router)

	backgroundFunc := func(ctx context.Context) {
		time.Sleep(time.Second * 10)
		serviceTxparser.Start(ctx)
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second * time.Duration(cfg.App.PollIntervalSec)):
				if err = serviceTxparser.ParseBlocks(ctx); err != nil {
					l.Errorf("failed to parse new blocks: %v", err)
				}
			}
		}
	}

	return srv, backgroundFunc, nil
}
