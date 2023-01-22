package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/vitalygudza/duty-app/internal/handler"
	myHttp "github.com/vitalygudza/duty-app/internal/http"
	"github.com/vitalygudza/duty-app/internal/job"
	"github.com/vitalygudza/duty-app/internal/notifier"
	"github.com/vitalygudza/duty-app/internal/repository"
	"github.com/vitalygudza/duty-app/internal/scheduler"
	"github.com/vitalygudza/duty-app/internal/service"
)

const (
	debugMode = true
)

// @title Duty manager API
// @version 2.0
// @description API Server for Duty manager application

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("initializing config error: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	notifyClient, err := notifier.NewSpaceClient(
		os.Getenv("SPACE_ENDPOINT"),
		os.Getenv("SPACE_TOKEN"),
	)
	if err != nil {
		logrus.Fatalf("error occurred while initializing the notification client: %s", err.Error())
	}

	scheduler, err := scheduler.NewScheduler()
	if err != nil {
		logrus.Fatalf("error occurred while initializing the scheduler: %s", err.Error())
	}

	repo := repository.NewRepository(db)
	services := service.NewService(repo, notifyClient)
	handlers := handler.NewHandler(services)
	dutyJob := job.NewDutyUpdater(services, debugMode)

	err = scheduler.AddJob(viper.GetString("duty_job_schedule"), dutyJob)
	if err != nil {
		logrus.Fatalf("error occurred while add job to scheduler: %s", err.Error())
	}

	scheduler.Start()

	srv := new(myHttp.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != http.ErrServerClosed {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	gracefulShutdown()

	ctx := scheduler.Stop()
	<-ctx.Done()

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db conection close: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}

func gracefulShutdown() {
	logrus.Println("Duty app started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Println("Duty app is shutting down")
}
