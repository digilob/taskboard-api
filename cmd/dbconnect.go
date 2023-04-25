package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func ConnectToDB() *sql.DB {
	postgresChan := make(chan *sql.DB, 1)

	go func() {
		postgresHost := os.Getenv("DB_HOST")
		postgresDSN := fmt.Sprintf("postgres://%s@%s:5432/tasks?sslmode=disable", os.Getenv("DB_USER"), postgresHost)
		cruddb, err := sql.Open("postgres", postgresDSN)
		if err != nil {
			log.Fatal(fmt.Errorf("error connecting to PostgreSQL db %+v", err))
		}

		log.Info("Pinging the PostgreSQL")
		for {
			if cruddbErr := cruddb.Ping(); cruddbErr != nil {
				log.Errorf("an error occurred connecting to the PostgreSQL db trying again in 20 seconds: %v\n", cruddbErr)
				time.Sleep(time.Second * 20)
			} else {
				log.Info("connected to PostgreSQL db")
				break
			}
		}

		postgresChan <- cruddb
	}()

	// wait for DB to be set up
	// for loop / channel is for applying multiple db connections if needed
	var postgresDB *sql.DB

	done := false

	for !done {
		log.Info("done is ", done)
		select {
		case db := <-postgresChan:
			if postgresDB == nil {
				postgresDB = db
				if db != nil {
					done = true
				}
			}
			log.Info("done is ", done)
		}
	}

	// instantiate a new PostgreSQL connection
	return postgresDB
}
