package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/fmarmol/todo/api"
	_ "github.com/lib/pq"
)

type ConfigDB struct {
	Host     string
	Port     string
	User     string
	Password string
}

type Server struct {
	db *sql.DB
}

func NewServer(cfg ConfigDB) *Server {
	db, err := sql.Open(
		"postgres",
		fmt.Sprintf("postgres://%v:%v@%v:%v/todo", cfg.User, cfg.Password, cfg.Host, cfg.Port),
	)
	if err != nil {
		log.Fatal("could not connecto to db")
	}
	if err := db.Ping(); err != nil {
		log.Fatal("could not ping db")
	}

}

func (s *Server) List(ctx context.Context, i api.Empty) (*api.TaskArray, error) {}
func (s *Server) Create(ctx context.Context, t *api.Task) (*Error, error)       {}
func (s *Server) Update(ctx context.Context, t *api.Task) (*Error, error)       {}
func (s *Server) Delete(ctx context.Context, t *api.Task) (*Error, error)       {}

func main() {

}
