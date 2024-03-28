package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mmkamron/basefit/internal/pkg/config"
	"github.com/mmkamron/basefit/internal/pkg/db"
	"github.com/mmkamron/basefit/internal/server/handler"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run() {
	router := http.NewServeMux()

	config := config.Load("./config/local.yaml")
	fmt.Println(config)
	db := db.Load(config)

	handler := handler.New(db)

	router.HandleFunc("GET /trainers", handler.Read)
	router.HandleFunc("POST /trainer", handler.Create)
	router.HandleFunc("PUT /trainer/:id", handler.Update)
	router.HandleFunc("DELETE /trainer/:id", handler.Delete)

	log.Fatal(http.ListenAndServe(":1337", router))
}
