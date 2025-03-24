package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	appHandlers "github.com/shuklarituparn/tarantool_crud_api/internal/handlers"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	Router *mux.Router
}

func NewServer() *Server {
	s := &Server{Router: mux.NewRouter()}
	s.routes()
	return s
}

func swaggerHandler(w http.ResponseWriter, r *http.Request) {
	basePath, err := os.Getwd()
	if err != nil {
		http.Error(w, "Failed to get working directory", http.StatusInternalServerError)
		return
	}

	filePath := filepath.Join(basePath, "docs", "swagger.json")

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to read swagger.json", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "application/json")
	if _, err := io.Copy(w, file); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func (s *Server) routes() {
	s.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("./static"))))

	s.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		filePath := filepath.Join("/app/static", "welcome.html")
		http.ServeFile(w, r, filePath)
	}).Methods("GET")

	s.Router.HandleFunc("/api/v1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Tarantool KV is up")
	}).Methods("GET")

	apiV1 := s.Router.PathPrefix("/api/v1").Subrouter()
	apiV1.HandleFunc("/kv", appHandlers.CreateKV).Methods("POST")
	apiV1.HandleFunc("/kv/{key}", appHandlers.UpdateKV).Methods("PUT")
	apiV1.HandleFunc("/kv/{key}", appHandlers.GetKV).Methods("GET")
	apiV1.HandleFunc("/kv/{key}", appHandlers.DeleteKV).Methods("DELETE")
	s.Router.PathPrefix("/docs/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"),
	))
	s.Router.HandleFunc("/swagger.json", swaggerHandler)

	s.Router.Handle("/metrics", promhttp.Handler())
	s.Router.Use(gorillaHandlers.CORS(
		gorillaHandlers.AllowedOrigins([]string{"*"}),
		gorillaHandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		gorillaHandlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	))
}
