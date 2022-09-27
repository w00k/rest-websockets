package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"w00k/go/rest-ws/database"
	"w00k/go/rest-ws/repository"
	"w00k/go/rest-ws/websocket"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Config struct {
	Port      string
	JWTSecret string
	DataUrl   string
}

type Server interface {
	Config() *Config
	Hub() *websocket.Hub
}

type Broker struct {
	config *Config
	router *mux.Router
	hub    *websocket.Hub
}

func (b *Broker) Config() *Config {
	return b.config
}

func (b *Broker) Hub() *websocket.Hub {
	return b.hub
}

func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("port is required")
	}
	if config.JWTSecret == "" {
		return nil, errors.New("secret is required")
	}
	if config.DataUrl == "" {
		return nil, errors.New("database is required")
	}
	broker := &Broker{
		config: config,
		router: mux.NewRouter(),
		hub:    websocket.NewHub(),
	}
	return broker, nil
}

func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	b.router = mux.NewRouter()
	handler := cors.Default().Handler(b.router)
	binder(b, b.router)
	repo, err := database.NewPostgresRepository(b.config.DataUrl)
	if err != nil {
		log.Fatal(err)
	}
	go b.hub.Run()
	repository.SetRespository(repo)
	log.Println("Starting server on port, ", b.Config().Port)
	if err := http.ListenAndServe(b.config.Port, handler); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
