package main

/*
 * _ blank identifier
 * use for neccessary import package but don't need to directly call it in the code
 */

import (
	"authentication-service/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var counts int64

type Config struct {
	Repo data.Repository
}

// type Config struct {
// 	DB     *sql.DB
// 	Models data.Models
// }

func main() {
	log.Println("Starting authentication service")

	// connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}

	// set up Config
	app := Config{}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	/*
	 * - if there is error try to connecting DB
	 * time sleep for 2 seconds and try again
	 * - if connecting attempt is more than 10 times
	 * give up and return nil
	 */

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		/*
		 * return out of loop and function
		 * if connecting attempt is more than 10
		 */
		if counts > 10 {
			log.Println(err)
			return nil
		}

		/*
		 * if there is error try to connecting DB
		 * time sleep for 2 seconds and try again
		 */
		log.Println("Backing off for two seconds")
		time.Sleep(2 * time.Second)
		continue
	}
}

func (app *Config) setupRepo(conn *sql.DB) {
	db := data.NewPostgresRepository(conn)
	app.Repo = db
}
