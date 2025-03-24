package config

import (
	"os"
	"github.com/shuklarituparn/tarantool_crud_api/pkg/logger"

	tarantool "github.com/tarantool/go-tarantool"
)

var Conn *tarantool.Connection

func InitDB() {
	var err error
	password := os.Getenv("TARANTOOL_PASSWORD")
	if password == "" {
		password = "Hellopassword123"
	}

	logger.Log.Info("Connecting to Tarantool on tarantool:3301 with user kv_user")
	Conn, err = tarantool.Connect("tarantool:3301", tarantool.Opts{
		User: "kv_user",
		Pass: password,
	})
	if err != nil {
		logger.Log.Fatal("Failed to connect to Tarantool: ", err)
	}

	logger.Log.Info("Successfully connected to Tarantool")
}
