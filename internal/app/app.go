package app

import (
	"TestHitalent/internal/config"
	"TestHitalent/internal/repository"
	"TestHitalent/internal/service"
	"TestHitalent/internal/transport"
	"TestHitalent/pkg/logger"
	"TestHitalent/pkg/postgres"
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"
)

type App struct {
	HiTalentServer *transport.HiTalentServer
	cfg                *config.Config
	ctx                context.Context
	wg                 sync.WaitGroup
	cancel             context.CancelFunc
}

func NewApp(cfg *config.Config,context context.Context) *App {
	db,err := postgres.New(cfg.Postgres)
	if err != nil {
		panic(err)
	}
	repo := repository.NewHiTalentRepository(db, context)
	srv := service.NewHiTalentService(context,repo)
	server := transport.NewHiTalentServer(cfg, srv,context)
	return &App{
		HiTalentServer: server,
		cfg:            cfg,
		ctx:            context,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	errCh := make(chan error, 1)
	a.wg.Add(1)
	go func() {
		logger.GetLoggerFromCtx(a.ctx).Info("Server started on address", zap.Any("address", a.cfg.Host+":"+a.cfg.Port))
		defer a.wg.Done()
		if err := a.HiTalentServer.Run(); err != nil {
			errCh <- err
			a.cancel()
		}
	}()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	select {
	case err := <-errCh:
		logger.GetLoggerFromCtx(a.ctx).Error("error running app", zap.Error(err))
		return err
	case <-a.ctx.Done():
		logger.GetLoggerFromCtx(a.ctx).Info("context done")
	}

	return nil
}
