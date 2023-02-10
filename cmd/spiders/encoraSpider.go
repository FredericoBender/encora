package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"encora/lib/spider"
	"encora/postgres"

	"github.com/jmoiron/sqlx"
)

var dbUserFlag = flag.String("dbUsername", "postgres", "is the database username connected")
var dbPasswordFlag = flag.String("dbPassword", "123", "is the database password connected")
var dbNameFlag = flag.String("dbName", "encora", "is the database name")
var dbHostFlag = flag.String("dbHost", "localhost", "is the database host")
var dbPortFlag = flag.String("dbPort", "5432", "is the database port")

func main() {
	flag.Parse()

	chromeDriverPath := "C:\\Users\\Frederico Bender\\Documents\\Algoritmos\\dependencias\\chromedriver.exe"
	port := int32(9595)
	debugLevel := int8(1)
	timelimit := 5

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", *dbHostFlag, *dbPortFlag, *dbUserFlag, *dbPasswordFlag, *dbNameFlag)
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	encoraExtractor := &spider.EncoraExtractor{
		DebugLevel: debugLevel,
		Logger:     log.New(os.Stdout, "EncoraExtractor: ", log.LstdFlags),
		TimeLimit:  timelimit,
		Config: spider.EncoraConfig{
			ChromeDriverPath: chromeDriverPath,
			Port:             port,
		},
	}
	encoraData, err := encoraExtractor.Run()
	if err != nil {
		log.Fatal(err)
	}

	ordersExporter := postgres.JobsExporter{
		DB: db,
	}
	err = ordersExporter.Run(&encoraData)
	if err != nil {
		log.Fatal(err)
	}

}
