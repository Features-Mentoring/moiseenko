package server

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lvestera/slot-machine/internal/server/config"
	"github.com/lvestera/slot-machine/internal/server/handlers"
	"github.com/lvestera/slot-machine/internal/storage"
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
	quit := make(chan os.Signal)
	go func() {
		<-quit
		log.Println("Receive interrupt signal. Server Close")
	}()

	repository, err := storage.NewDBRepository()
	if err != nil {
		return err
	}

	return http.ListenAndServe("localhost:8081", Router(s.Cfg, repository))
}

func Router(cfg *config.Config, db *storage.DBRepository) chi.Router {
	r := chi.NewRouter()

	r.Method(http.MethodPost, "/get-config", handlers.GetConfigHandler{Cfg: cfg})
	r.Method(http.MethodPost, "/send-result", handlers.GetResultHandler{Db: *db})

	return r
}
