package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 800*time.Millisecond)

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)

	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	body, error := io.ReadAll(res.Body)

	if error != nil {
		panic(error)
	}

	log.Println(string(body))

	amount := "2.2"

	println(amount)

	ctx2, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)

	defer cancel()

	println("http://localhost:8080/incluirCotacao/?amount=" + amount)

	req, err = http.NewRequestWithContext(ctx2, "POST", "http://localhost:8080?amount="+amount, nil)

	if err != nil {
		panic(err)
	}

	res2, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}

	defer res2.Body.Close()

	body2, error := io.ReadAll(res2.Body)

	if error != nil {
		panic(error)
	}

	log.Println(string(body2))

}
