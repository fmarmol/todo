package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/fmarmol/todo/api"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

// ConfigDB struct
type ConfigDB struct {
	Host     string
	Port     string
	User     string
	Password string
}

// Server struct
type Server struct {
	db *sql.DB
}

// NewServer creates a new server
func NewServer(cfg *ConfigDB) *Server {
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

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXIST tasks(
			id SERIAL,
			uid INT,
			start TIMESTAMP,
			end TIMESTAMP,
			priority VARCHAR(10),
			status VARCHAR(10),
			subject TEXT
		)`)
	if err != nil {
		log.Fatal("something got wrong", err)
	}
	s := new(Server)
	s.db = db
	return s
}

// List implementation
func (s *Server) List(ctx context.Context, i *api.Empty) (*api.TaskArray, error) { return nil, nil }

// Create implementation
func (s *Server) Create(ctx context.Context, t *api.Task) (*api.Error, error) { return nil, nil }

// Update implemenation
func (s *Server) Update(ctx context.Context, t *api.Task) (*api.Error, error) { return nil, nil }

// Delete implementation
func (s *Server) Delete(ctx context.Context, t *api.Task) (*api.Error, error) { return nil, nil }

// NewConfig creates a new config from os env variables
func NewConfig() *ConfigDB {
	return &ConfigDB{
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
	}
}

// Update the default config
func (cfg *ConfigDB) Update(user, password, host, port string) {
	if user != "" {
		cfg.User = user
	}
	if password != "" {
		cfg.Password = password
	}
	if host != "" {
		cfg.Host = host
	}
	if port != "" {
		cfg.Port = port
	}
}

func main() {
	listen := flag.String("listen", "8080", "port listening")
	host := flag.String("host", "", "postgres host")
	port := flag.String("port", "", "postgres port")
	user := flag.String("user", "", "postgres user")
	password := flag.String("password", "", "postgres user's password")
	flag.Usage()

	cfg := NewConfig()
	cfg.Update(*user, *password, *host, *port)

	lis, err := net.Listen("tpc", ":"+*listen)
	if err != nil {
		log.Fatal()

	}
	defer lis.Close()
	server := grpc.NewServer()
	api.RegisterToDoApiServer(server, NewServer(cfg))
	server.Serve(lis)
}
