package Banco

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	database, err := sql.Open("sqlite3", "./cotacoes.db")

	if err != nil {
		log.Fatal(err)
	}

	db = database
}

func Close() {
	db.Close()
}

func InsertCotation(dateCotation time.Time, amount string) error {
	statemant, _ := db.Prepare("Create table if not exists cotacao (id integer primary key, DT_Cotacao DATETIME, VR_Cotacao string )")

	statemant.Exec()

	statemant, _ = db.Prepare("insert into cotacao (dt_cotacao, vr_cotacao) values (?,?)")

	statemant.Exec(dateCotation, amount)

	rows, _ := db.Query("SELECT id, DT_Cotacao, VR_Cotacao FROM cotacao")

	var id int
	var dtCotacao time.Time
	var vrCotacao string

	for rows.Next() {
		rows.Scan(&id, &dtCotacao, &vrCotacao)
		fmt.Printf("ID: %d, Data: %s, Valor: %s\n", id, dtCotacao.Format("2006-01-02 15:04:05"), vrCotacao)
	}

	return nil
}
