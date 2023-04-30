package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", Cotacao)
	http.ListenAndServe(":8080", mux)
}

func Cotacao(w http.ResponseWriter, r *http.Request) {
	cotacao := consultaCotacao()
	SalvaCotacao(&cotacao.USDBRL)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cotacao.USDBRL)
}

func consultaCotacao() CotacaoResponse {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", URL, nil)
	if err != nil {
		panic(err)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	resp, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var cotacaoResponse CotacaoResponse
	err = json.Unmarshal(resp, &cotacaoResponse)
	if err != nil {
		panic(err)
	}

	return cotacaoResponse
}

func SalvaCotacao(cotacao *DadosCotacaoResponse) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	dsn := "root:root@tcp(localhost:3306)/goexpert"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&DadosCotacaoResponse{})

	db.WithContext(ctx).Create(cotacao)
}

const URL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

type CotacaoResponse struct {
	USDBRL DadosCotacaoResponse
}

type DadosCotacaoResponse struct {
	ID         int    `gorm:"primaryKey"`
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
}
