package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lvestera/slot-machine/internal/server/config"
	"github.com/lvestera/slot-machine/internal/server/handlers"
	"github.com/lvestera/slot-machine/internal/storage"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	Cfg *config.Config
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		Cfg: cfg,
	}
}

func (s *Server) Run() error {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()

	repository, err := storage.NewDBRepository(s.Cfg.DBConnection)
	if err != nil {
		return err
	}

	httpServer := &http.Server{
		Addr:    s.Cfg.Host,
		Handler: Router(s.Cfg, repository),
	}

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return httpServer.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		return httpServer.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func Router(cfg *config.Config, db *storage.DBRepository) chi.Router {
	log.Println("Server start at ", cfg.Host)
	r := chi.NewRouter()

	r.Method(http.MethodPost, "/get-config", handlers.GetConfigHandler{Cfg: cfg})
	r.Method(http.MethodPost, "/send-result", handlers.GetResultHandler{Db: *db})
	r.Method(http.MethodGet, "/get-chart", handlers.GetChartHandler{Db: *db})

	return r
}
