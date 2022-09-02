package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"w00k/go/rest-ws/database"
	"w00k/go/rest-ws/repository"

	"github.com/gorilla/mux"
)

type Config struct {
	Port      string
	JWTSecret string
	DataUrl   string
}

type Server interface {
	Config() *Config
}

type Broker struct {
	config *Config
	router *mux.Router
}

func (b *Broker) Config() *Config {
	return b.config
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
	}
	return broker, nil
}

func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	b.router = mux.NewRouter()
	binder(b, b.router)
	repo, err := database.NewPostgresRepository(b.config.DataUrl)
	if err != nil {
		log.Fatal(err)
	}
	repository.SetRespository(repo)
	log.Println("Starting server on port, ", b.Config().Port)
	if err := http.ListenAndServe(b.config.Port, b.router); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
