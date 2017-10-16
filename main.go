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
	url := fmt.Sprintf("postgres://%v:%v@%v:%v/todo?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port)
	db, err := sql.Open(
		"postgres",
		url,
	)
	if err != nil {
		log.Fatalf("could not connecto to db: %v\n", url)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("could not ping db: %v\n", url)
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS tasks(
			id SERIAL,
			uid INT,
			start_date TIMESTAMP,
			end_date TIMESTAMP,
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
func (s *Server) List(ctx context.Context, i *api.Empty) (*api.TaskArray, error) {
	rows, err := s.db.Query("SELECT uid, start, end, priority, status, subject FROM tasks;")
	if err != nil {
		return nil, err
	}

	ret := &api.TaskArray{Tasks: []*api.Task{}}
	for rows.Next() {
		var task *api.Task
		if err := rows.Scan(task.Uid, task.Start, task.End, task.Priority, task.Status, task.Subject); err != nil {
			return nil, err
		}
		ret.Tasks = append(ret.Tasks, task)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	return ret, nil
}

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
	flag.Parse()

	cfg := NewConfig()
	cfg.Update(*user, *password, *host, *port)

	lis, err := net.Listen("tcp", ":"+*listen)
	if err != nil {
		log.Fatal(err)

	}
	defer lis.Close()
	server := grpc.NewServer()
	api.RegisterToDoApiServer(server, NewServer(cfg))
	server.Serve(lis)
}
