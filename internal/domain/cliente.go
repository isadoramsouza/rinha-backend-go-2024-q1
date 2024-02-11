package domain

type Cliente struct {
	ID     int    `json:"id"`
	Nome   string `json:"nome"`
	Limite int64  `json:"limite"`
	Saldo  int64  `json:"saldo"`
}
