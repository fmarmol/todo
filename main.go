package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/fmarmol/todo/api"
	"github.com/lib/pq"
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
		log.Fatalf("could not connect to db %v:%v\n", url, err)
	}

LOOP:
	for i := 0; i < 3; i++ {
		if err := db.Ping(); err != nil {
			if i < 2 {
				log.Printf("trying to ping db:%v\n", url)
				time.Sleep(time.Second)
			} else {
				log.Fatalf("could not ping db %v:%v\n", url, err)
			}
		} else {
			break LOOP
		}
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS tasks(
			id SERIAL,
			uid INT UNIQUE,
			start_date DATE,
			end_date DATE,
			priority VARCHAR(10) NOT NULL,
			status VARCHAR(10) NOT NULL,
			subject TEXT NOT NULL
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
	rows, err := s.db.Query("SELECT uid, start_date, end_date, priority, status, subject FROM tasks;")
	if err != nil {
		return nil, err
	}

	ret := &api.TaskArray{Tasks: []*api.Task{}}
	for rows.Next() {
		var task api.Task
		var start, end pq.NullTime
		var priority, status string
		if err := rows.Scan(&task.Uid, &start, &end, &priority, &status, &task.Subject); err != nil {
			return nil, err
		}
		if start.Valid {
			task.Start = start.Time.Unix()
		}
		if end.Valid {
			task.End = end.Time.Unix()
		}
		task.Status = api.Task_Status(api.Task_Status_value[status])
		task.Priority = api.Task_Priority(api.Task_Priority_value[priority])
		ret.Tasks = append(ret.Tasks, &task)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	return ret, nil
}

// Create implementation
func (s *Server) Create(ctx context.Context, t *api.Task) (*api.Error, error) {
	_, err := s.db.Exec("INSERT INTO tasks (uid,start_date,end_date,priority, status, subject) VALUES ($1,$2,$3,$4,$5,$6)",
		t.Uid,
		time.Unix(t.Start, 0),
		time.Unix(t.End, 0),
		api.Task_Priority_name[int32(t.Priority)],
		api.Task_Status_name[int32(t.Status)],
		t.Subject)
	if err != nil {
		return &api.Error{Type: api.Error_FAIL, Description: err.Error()}, err
	}
	return &api.Error{Type: api.Error_SUCCESS}, nil
}

// Update implemenation
func (s *Server) Update(ctx context.Context, t *api.Task) (*api.Error, error) {
	_, err := s.db.Exec("UPDATE tasks SET start_date=$1, end_date=$2, priority=$3, status=$4, subject=$5 WHERE ui=$6",
		t.Start, t.End, t.Priority, t.Status, t.Subject)
	if err != nil {
		return &api.Error{Type: api.Error_FAIL, Description: err.Error()}, err
	}
	return &api.Error{Type: api.Error_SUCCESS}, nil
}

// Delete implementation
func (s *Server) Delete(ctx context.Context, t *api.Task) (*api.Error, error) {
	_, err := s.db.Exec("DELETE FROM tasks where uid=$1", t.Uid)
	if err != nil {
		return &api.Error{Type: api.Error_FAIL, Description: err.Error()}, err
	}
	return &api.Error{Type: api.Error_SUCCESS}, nil
}

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
