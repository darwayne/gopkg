/*
 * Copyright (c) 2018. Darwayne
 */

package db

import (
	"github.com/gobuffalo/packr"
	"github.com/jinzhu/gorm"
	"github.com/rubenv/sql-migrate"
	"log"
)

func RunMigrations(db *gorm.DB) {
	box := packr.NewBox("../../pkg/db/migrate")

	migrations := &migrate.PackrMigrationSource{
		Box: box,
	}

	_, err := migrate.Exec(db.DB(), "postgres", migrations, migrate.Up)

	if err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}
}

func GetDB(connectionStr string) *gorm.DB {
	db, err := gorm.Open("postgres", connectionStr)
	db.DB().SetMaxIdleConns(5)
	db.DB().SetMaxOpenConns(5)

	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}

	return db
}
