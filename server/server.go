package main

import (
	"context"
	"encoding/json"
	"errors"
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
	cotacao, err := consultaCotacao()
	if err != nil {
		panic(err)
	}
	err = SalvaCotacao(&cotacao.USDBRL)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cotacao.USDBRL)
}

func consultaCotacao() (*CotacaoResponse, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", URL, nil)
	if err != nil {
		return nil, errors.New("Ocorreu um erro ao tentar criar a request da consulta da cotação do dolar")
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New("Ocorreu um erro ao realizar consulta da cotação do dolar")
	}

	defer response.Body.Close()

	resp, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("Ocorreu um erro ao realizar a leitura do response")
	}

	var cotacaoResponse CotacaoResponse
	err = json.Unmarshal(resp, &cotacaoResponse)
	if err != nil {
		return nil, errors.New("Ocorreu um erro ao converter o response")
	}

	return &cotacaoResponse, nil
}

func SalvaCotacao(cotacao *DadosCotacaoResponse) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	dsn := "root:root@tcp(localhost:3306)/goexpert"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return errors.New("Ocorreu um erro ao realizar a conexão com o banco de dados")
	}

	db.WithContext(ctx).Create(cotacao)
	return nil
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
