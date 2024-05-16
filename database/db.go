package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/V-Ader/Loyality_GO/config"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	c := config.GetDBConfig()
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.Username, c.Password, c.Database)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Connected to the database")
	return db, nil
}
