package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"reinaldosantosrosa/Pos-Go-Expert/Cotacao/Banco"
	"reinaldosantosrosa/Pos-Go-Expert/Cotacao/Util"
	"time"
)

type CotacaoDolar struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func main() {

	http.HandleFunc("/cotacao", cotacao)
	http.HandleFunc("/", IncluirCotacao)
	http.ListenAndServe(":8080", nil)

}

func cotacao(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Println("Request cotacao iniciada")

	defer log.Println("Request cotacao finalizada")

	select {
	case <-time.After(200 * time.Millisecond):
		log.Println("Realizando sua Cotação")

		cotacao, error := BuscaCotacao()

		if error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Contenty-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(cotacao.USDBRL.Bid)

		err := Util.AppendCreateArq("Valor da cotacao do Dolar:  "+string(cotacao.USDBRL.Bid)+"\n", "Arquivo.txt")

		if err != nil {
			panic(err)
		}

	case <-ctx.Done():
		log.Println("Request cancelada pelo cliente")
	}
}

func IncluirCotacao(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	amountParam := r.URL.Query().Get("amount")
	if amountParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	log.Println("Request do banco de dados iniciada")

	defer log.Println("Request do banco finalizada")

	select {
	case <-time.After(300 * time.Millisecond):
		log.Println("Registrando Cotação no banco de dados")
		Banco.InsertCotation(time.Now(), amountParam)

	case <-ctx.Done():
		log.Println("Erro ao acessar o banco de dados")
	}

}

func BuscaCotacao() (*CotacaoDolar, error) {

	resp, error := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")

	if error != nil {
		return nil, error
	}

	defer resp.Body.Close()

	body, error := io.ReadAll(resp.Body)

	if error != nil {
		return nil, error
	}

	var c CotacaoDolar

	error = json.Unmarshal(body, &c)

	if error != nil {
		return nil, error
	}

	return &c, nil
}
