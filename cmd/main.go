package main

import (
	"net/http"

	"github.com/shuklarituparn/tarantool_crud_api/config"
	"github.com/shuklarituparn/tarantool_crud_api/internal/prometheus"
	"github.com/shuklarituparn/tarantool_crud_api/internal/server"
	"github.com/shuklarituparn/tarantool_crud_api/pkg/logger"

	"github.com/joho/godotenv"
)

// @title Key-Value Store API
// @version 1.0
// @description Документация по API для простого хранилища KV на базе Tarantool.
// @contact.name API Support
// @contact.url https://shukla.ru
// @contact.email shukla.r@phystech.edu
// @license.name MIT License
// @license.url https://opensource.org/licenses/MIT
// @in header
// @name Authorization
func main() {
	if err := godotenv.Load(); err != nil {
		logger.Log.Warn("No .env file found, using default environment variables")
	}
	prometheus.Init()
	logger.Log.Info("Prometheus metrics initialized")
	logger.Log.Info("Initializing Key-Value Store...")

	config.InitDB()

	s := server.NewServer()
	logger.Log.Info("Server running on port :5005")

	if err := http.ListenAndServe(":5005", s.Router); err != nil {
		logger.Log.Fatal("Server failed to start: ", err)
	}
}
