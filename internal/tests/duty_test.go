package service

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vitalygudza/duty-app/internal/job"
	"github.com/vitalygudza/duty-app/internal/notifier"
	"github.com/vitalygudza/duty-app/internal/repository"
)

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}

func TestDuty(t *testing.T) {
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

	repo := repository.NewRepository(db)

	notifyClient, err := notifier.NewSpaceClient(
		os.Getenv("SPACE_ENDPOINT"),
		os.Getenv("SPACE_TOKEN"),
	)
	if err != nil {
		logrus.Fatalf("error occurred while initializing the notification client: %s", err.Error())
	}

	// ds := NewDutyService(repo.Team, repo.Teammate, repo.Duty)
	services := NewService(repo, notifyClient)
	dutyJob := job.NewDutyUpdater(services, true)

	dutyJob.Run()
}
