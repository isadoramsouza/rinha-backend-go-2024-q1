package domain

import "time"

type Transacao struct {
	ID          int       `json:"id"`
	Tipo        string    `json:"tipo"`
	Descricao   string    `json:"descricao"`
	Valor       int64     `json:"valor"`
	ClienteID   int       `json:"cliente_id"`
	RealizadaEm time.Time `json:"realizada_em"`
}

type TransacaoResponse struct {
	Limite int64 `json:"limite"`
	Saldo  int64 `json:"saldo"`
}
