package database

import (
	"context"
	"crowdfunding/config"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Initpostgre(ctx context.Context) *gorm.DB {
	connection := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		config.GlobalEnv.PostgreHost,
		config.GlobalEnv.PostgreUser,
		config.GlobalEnv.PostgrePassword,
		config.GlobalEnv.PostgreDBName,
		config.GlobalEnv.PostgrePort,
		config.GlobalEnv.PostgreSSLMode)

	db, err := gorm.Open(postgres.Open(connection), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database postgre")
	}

	return db
}
