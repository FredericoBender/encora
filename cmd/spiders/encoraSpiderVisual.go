package main

import (
	"fmt"
	"log"
	"os"

	"encora/bubbletea"
	"encora/lib/spider"
	"encora/postgres"

	"github.com/jmoiron/sqlx"
)

func main() {
	chromeDriverPath := "C:\\Users\\Frederico Bender\\Documents\\Algoritmos\\dependencias\\chromedriver.exe"
	chromeDriverPort := int32(9595)
	debugLevel := int8(0)

	inputs := bubbletea.ReadInputs()
	dbUserName := inputs["dbUsername"]
	dbPassword := inputs["dbPassword"]
	dbName := inputs["dbName"]
	dbHost := inputs["dbHost"]
	dbPort := inputs["dbPort"]

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUserName, dbPassword, dbName)
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	encoraExtractor := &spider.EncoraExtractor{
		DebugLevel: debugLevel,
		Logger:     log.New(os.Stdout, "EncoraExtractor: ", log.LstdFlags),
		Config: spider.EncoraConfig{
			ChromeDriverPath: chromeDriverPath,
			Port:             chromeDriverPort,
		},
		IsVisualExecution: true,
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
