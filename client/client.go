package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	valorCambio := ConsultaBid()
	SalvaBid(valorCambio.Bid)

}

const URL_COTACAO = "http://localhost:8080/cotacao"

type ValorCambio struct {
	Bid string `json:"bid"`
}

func ConsultaBid() ValorCambio {
	c := http.Client{}

	response, err := c.Get(URL_COTACAO)
	if err != nil {
		fmt.Println("Ocorreu um erro ao tentar consultar a cotação do dolar")
		panic(err)
	}

	defer response.Body.Close()

	resp, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Ocorreu um erro ao tentar ler o response da consulta da cotação do dolar")
		panic(err)
	}

	var valorCambio ValorCambio
	err = json.Unmarshal(resp, &valorCambio)
	if err != nil {
		fmt.Println("Ocorreu um erro ao realizar o parse do valor do cambio")
		panic(err)
	}

	return valorCambio
}

func SalvaBid(bid string) {
	file, err := os.OpenFile("cotacao.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	_, err = file.WriteString("Dólar: " + bid + "\n")
	if err != nil {
		panic(err)
	}
}
