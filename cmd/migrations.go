package main

import (
	"database/sql"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	log "github.com/sirupsen/logrus"
)

func HandleMigrations(db *sql.DB) error {
	if db == nil {
		log.Fatal("Database is nil")
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Error("an error occurred with migrations1: ", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		log.Error("error when migrating: ", err)
	}

	version, _, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.WithField("version", version).Error(err)
		return err
	}

	if err = m.Up(); err != nil {
		log.Error("an error occurred with migrations: ", err)
	}

	nversion, _, err := m.Version()
	if err != nil {
		log.Error(err)
		return err
	}

	log.Infof("migrated PostgreSQL DB from version %d to version %d", version, nversion)

	return nil
}
