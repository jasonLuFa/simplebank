package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/jasonLuFa/simplebank/util"
	_ "github.com/lib/pq"
)

// const (
// 	dbDriver = "postgres"
// 	dbSource = "postgresql://admin:admin@localhost:5433/simple_bank?sslmode=disable"
// )

var testQueries *Queries
var db *sql.DB

func TestMain(m *testing.M){
	config,err:= util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:",err)
	}
	db, err = sql.Open(config.DBDriver,config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:",err)
	}

	testQueries = New(db)

	os.Exit(m.Run())
}
