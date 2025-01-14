package main

import (
	"Cotacao/Util"
	"encoding/json"
	"io"
	"log"
	"net/http"
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
	http.ListenAndServe(":8080", nil)

}

func cotacao(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Println("Request iniciada")

	defer log.Println("Request finalizada")

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
		json.NewEncoder(w).Encode("O Valor da Cotação de hoje " + string(cotacao.USDBRL.Bid))

		err := Util.AppendCreateArq("Valor da cotacao do Dolar: "+string(cotacao.USDBRL.Bid)+"\n", "Arquivo.txt")

		if err != nil {
			panic(err)
		}

	case <-ctx.Done():
		log.Println("Request cancelada pelo cliente")
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
