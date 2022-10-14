package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://admin:admin@localhost:5433/simple_bank?sslmode=disable"
)

var testQueries *Queries
var db *sql.DB

func TestMain(m *testing.M){
	var err error

	db, err = sql.Open(dbDriver,dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:s",err)
	}

	testQueries = New(db)

	os.Exit(m.Run())
}
