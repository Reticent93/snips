package main

import (
	"database/sql"
	"flag"
	"github.com/Reticent93/snips/pkg/models/postgres"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snips    *postgres.SnipModel
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	dsn := os.Getenv("DATABASE_URL")

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//db, err := openDB(*dsn)
	//if err != nil {
	//	errorLog.Fatal(err)
	//}
	defer func(db *bun.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	var app = &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snips:    &postgres.SnipModel{DB: db},
	}
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s\n", *addr)
	err = srv.ListenAndServe()
	if err != nil {

		errorLog.Fatal(nil, err)
	}

}

//func openDB(dsn string) (*sql.DB, error)  {
//	db, err := pg.Open("postgres", dsn)
//	if err != nil {
//		return nil, err
//	}
//	if err = db.Ping(); err != nil {
//		return nil, err
//	}
//	return db, nil
//}
