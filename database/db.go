package database

import (
	"database/sql"

	_ "github.com/lib/pq"
	"log"
	"os"

	"github.com/joho/godotenv"
)
var DB *sql.DB

func InitDB() *sql.DB {
	err := godotenv.Load()
	if err!=nil{
		log.Fatal(err)
		return nil
	}
	connStr :=os.Getenv("connStr")
	//using docker for pulling image of postgres
	db,err:=sql.Open("postgres",connStr)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
		return nil
	}
	DB=db
	return db
}