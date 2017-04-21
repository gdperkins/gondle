package main

import (
	"flag"
	"net/http"

	"os"

	"github.com/gondle/config"
)

func main() {
	environment := flag.String("e", "development", "")
	flag.Parse()

	config.Init(*environment)
	config := config.GetConfig()
	os.Setenv("env", config.GetString("server.environment"))

	http.HandleFunc(config.GetString("github.hookpath"), HandleHook)
	http.ListenAndServe(config.GetString("server.port"), nil)
}
