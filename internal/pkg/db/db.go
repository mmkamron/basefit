package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/mmkamron/basefit/internal/pkg/config"
)

func Load(conf *config.Config) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", conf.DB.Host, conf.DB.Port, conf.DB.User, conf.DB.Password, conf.DB.Name)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	duration, err := time.ParseDuration(conf.DB.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(conf.DB.MaxOpenConns)
	db.SetMaxIdleConns(conf.DB.MaxIdleConns)
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
