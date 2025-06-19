package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type ViaBrasil struct {
	Cep     string `json:"cep"`
	Estado  string `json:"state"`
	Cidade  string `json:"city"`
	Bairro  string `json:"neighborhood"`
	Rua     string `json:"street"`
	Servico string `json:"service"`
}

func main() {
	http.HandleFunc("/", BuscaCepHandler)
	http.ListenAndServe(":8080", nil)
}

func BuscaCepHandler(w http.ResponseWriter, r *http.Request) {

	c1 := make(chan ViaCEP)
	c2 := make(chan ViaBrasil)

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	cepParam := r.URL.Query().Get("cep")

	if cepParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	go func() {
		time.Sleep(time.Second * 1)
		cep1, _ := BuscaCepViaCep(cepParam)

		c1 <- *cep1
	}()

	go func() {
		time.Sleep(time.Second * 4)
		cepbrasil, _ := BuscaCepViaBrasil(cepParam)

		c2 <- *cepbrasil
	}()

	select {
	case ViaCep := <-c1:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("===========Resposta API ViaCep============\n"))
		json.NewEncoder(w).Encode(ViaCep)
	case ViaBrasil := <-c2:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("===========Resposta API VIABRASIL============\n"))
		json.NewEncoder(w).Encode(ViaBrasil)
	case <-time.After(time.Second * 10):
		fmt.Println("timeout")
		return
	}

}

func BuscaCepViaCep(cep string) (*ViaCEP, error) {
	resp, _ := http.Get("https://viacep.com.br/ws/" + cep + "/json")

	body, _ := io.ReadAll(resp.Body)

	var c ViaCEP

	json.Unmarshal(body, &c)

	return &c, nil
}

func BuscaCepViaBrasil(cep string) (*ViaBrasil, error) {
	resp, _ := http.Get("https://brasilapi.com.br/api/cep/v1/" + cep)

	body, _ := io.ReadAll(resp.Body)

	var c ViaBrasil

	json.Unmarshal(body, &c)

	return &c, nil
}
