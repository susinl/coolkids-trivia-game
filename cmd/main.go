package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	_ "time/tzdata"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/cors"
	"github.com/susinl/coolkids-trivia-game/database"
	"github.com/susinl/coolkids-trivia-game/logz"
	"github.com/susinl/coolkids-trivia-game/middleware"
	"github.com/susinl/coolkids-trivia-game/question"

	"github.com/spf13/viper"
)

func init() {
	runtime.GOMAXPROCS(1)
	initTimezone()
	initViper()
}

func initTimezone() {
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Printf("error loading location 'Asia/Bangkok': %v\n", err)
	}
	time.Local = ict
}

func initViper() {
	viper.SetDefault("app.name", "coolkids-trivia-game")
	viper.SetDefault("app.port", "9090")
	viper.SetDefault("app.timeout", "60s")
	viper.SetDefault("app.context", "/coolkids")

	viper.SetDefault("log.level", "debug")
	viper.SetDefault("log.env", "dev")

	viper.SetDefault("mysql.host", "localhost")
	viper.SetDefault("mysql.port", "3306")
	viper.SetDefault("mysql.username", "sa")
	viper.SetDefault("mysql.database", "db")
	viper.SetDefault("mysql.timeout", 100)

	viper.SetDefault("question.timeout", "20s")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func main() {
	logger, err := logz.NewLogConfig()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	db, err := database.NewMySQLConn()
	if err != nil {
		logger.Error(err.Error())
	}
	defer db.Close()

	route := mux.NewRouter()

	cfgCors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},                                               // All origins
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}, // Allowing only get, just an example
		AllowedHeaders:   []string{"Content-Type", "Origin", "Authorization", "Accept"},
		AllowCredentials: true,
	})

	middle := middleware.NewMiddleware(logger)

	mux := route.PathPrefix(viper.GetString("app.context")).Subrouter()
	mux.Use(middle.JsonMiddleware)

	mux.Handle("/start", question.NewStartQuestion(
		logger,
		question.NewQueryParticipantByCodeFn(db),
		question.NewQueryQuestionByStatusFn(db),
		question.NewUpdateQuestionStatusAndParticipantInfoFn(db),
	)).Methods(http.MethodPost)

	mux.Handle("/submit", question.NewSubmitAnswer(
		logger,
		question.NewQueryParticipantAndAnswerFn(db),
		question.NewUpdateParticipantAnswerAndStatusFn(db),
	)).Methods(http.MethodPost)

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", viper.GetString("app.port")),
		Handler:      cfgCors.Handler(route),
		ReadTimeout:  viper.GetDuration("app.timeout"),
		WriteTimeout: viper.GetDuration("app.timeout"),
		IdleTimeout:  viper.GetDuration("app.timeout"),
	}

	logger.Info(fmt.Sprintf("⇨ http server started on [::]:%s", viper.GetString("app.port")))

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Info(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("app.timeout"))
	defer cancel()

	srv.Shutdown(ctx)

	logger.Info("shutting down")
	os.Exit(0)
}
