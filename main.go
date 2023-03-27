package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"flag"

	"github.com/gkontos/goapi/controller"
	"github.com/gkontos/goapi/db"
	"github.com/gkontos/goapi/logger"
	"github.com/go-chi/chi/v5"
)

var (
	router *chi.Mux
	port   *string
)

func initModule() {

	debug := flag.Bool("debug", getEnvOrBool("APP_DEBUG", false), "sets log level to debug")
	port = flag.String("port", getEnvOrString("APP_PORT", "8080"), "application port")
	console_log := flag.Bool("log_to_console", getEnvOrBool("APP_CONSOLE_LOG", false), "sets log output to console")
	allowed_origins := getEnvOrString("ALLOWED_ORIGIN", "http://localhost")
	flag.Parse()

	logger.InitLogger(*debug, *console_log)
	dbHandler := db.NewDbHandler()

	controlHandler := controller.NewController(dbHandler)
	router = controller.NewRouter(allowed_origins, controlHandler, dbHandler).SetupRouter()
}

// appengine will not run the main method.
// so to run in appengine, the router needs to be initialized in init()
func init() {
	initModule()
	http.Handle("/", router)
}

// since appengine is not running the main method, we can use this method for local dev or running in a container
func main() {
	logger.Logger.Info().Msg(fmt.Sprintf("Listening on port %s", *port))
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		logger.Logger.Error().AnErr("failed to start", err)
	}
}

func getEnvOrString(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defaultValue
	}
	return value
}

func getEnvOrBool(key string, defaultValue bool) bool {
	envValue, exists := os.LookupEnv(key)
	if exists {
		val, err := strconv.ParseBool(envValue)
		if err != nil {
			logger.Logger.Error().Msg(fmt.Sprintf("failed to parse %s, using default value", key))
		}
		return val
	}
	return defaultValue
}
