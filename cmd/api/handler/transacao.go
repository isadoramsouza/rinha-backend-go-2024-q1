package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/isadoramsouza/rinha-backend-go-2024-q1/internal/domain"
	transacao "github.com/isadoramsouza/rinha-backend-go-2024-q1/internal/transacao"
	"github.com/isadoramsouza/rinha-backend-go-2024-q1/pkg/web"
)

var (
	ErrInvalidJson = errors.New("invalid json")
	ErrNotFound    = errors.New("cliente not found")
	InvalidDtoErr  = errors.New("invalid request")
	LimitErr       = errors.New("limit error")
)

type TransacaoRequest struct {
	Valor     int64  `json:"valor" validate:"required,gt=0"`
	Tipo      string `json:"tipo" validate:"required,len=1"`
	Descricao string `json:"descricao" validate:"required,len=10"`
}

type TransacaoResponse struct {
	Limite int64 `json:"limite"`
	Saldo  int64 `json:"saldo"`
}

func (t *TransacaoRequest) Validate() error {

	if (t.Tipo != "c" && t.Tipo != "d") || (len(t.Descricao) < 1 || len(t.Descricao) > 10) { // validate body
		return InvalidDtoErr
	}

	return nil
}

type TransacaoController struct {
	transacaoService transacao.Service
}

func NewTransacao(s transacao.Service) *TransacaoController {
	return &TransacaoController{
		transacaoService: s,
	}
}

func (t *TransacaoController) CreateTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {

		input := &TransacaoRequest{}

		err := c.ShouldBindJSON(input)
		if err != nil {
			web.Error(c, http.StatusUnprocessableEntity, ErrInvalidJson.Error())
			return
		}

		if err := input.Validate(); err != nil {
			web.Error(c, http.StatusUnprocessableEntity, ErrInvalidJson.Error())
			return
		}

		id, _ := strconv.Atoi(c.Param("id"))

		cliente, err := t.transacaoService.GetBalance(c, id)
		if err != nil {
			if err.Error() == ErrNotFound.Error() {
				web.Error(c, http.StatusNotFound, ErrNotFound.Error())
				return
			}
			web.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		var newBalance int64

		if "d" == input.Tipo {
			newBalance = cliente.Saldo - input.Valor
		} else {
			newBalance = cliente.Saldo + input.Valor
		}

		if (cliente.Limite + newBalance) < 0 {
			web.Error(c, http.StatusUnprocessableEntity, LimitErr.Error())
			return
		}

		newTransacao := domain.Transacao{
			ClienteID: id,
			Tipo:      input.Tipo,
			Descricao: input.Descricao,
			Valor:     input.Valor,
		}

		err = t.transacaoService.CreateTransaction(c, newTransacao, newBalance)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		newTransacaoResponse := TransacaoResponse{
			Limite: cliente.Limite,
			Saldo:  newBalance,
		}
		web.Success(c, http.StatusOK, newTransacaoResponse)
	}
}

func (t *TransacaoController) GetExtrato() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))

		extrato, err := t.transacaoService.GetExtrato(c, id)
		if err != nil {
			if err.Error() == ErrNotFound.Error() {
				web.Error(c, http.StatusNotFound, ErrNotFound.Error())
				return
			}
			web.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		web.Success(c, http.StatusOK, extrato)
	}
}
