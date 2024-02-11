package domain

import "time"

type Extrato struct {
	Saldo             Saldo             `json:"saldo"`
	UltimasTransacoes []UltimaTransacao `json:"ultimas_transacoes"`
}

type Saldo struct {
	Total       int64     `json:"total"`
	DataExtrato time.Time `json:"data_extrato"`
	Limite      int64     `json:"limite"`
}

type UltimaTransacao struct {
	Tipo        string    `json:"tipo"`
	Descricao   string    `json:"descricao"`
	Valor       int64     `json:"valor"`
	RealizadaEm time.Time `json:"realizada_em"`
}
