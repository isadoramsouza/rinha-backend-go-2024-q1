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

func (t *TransacaoRequest) Validate() error {

	if (t.Tipo != "c" && t.Tipo != "d") || (len(t.Descricao) < 1 || len(t.Descricao) > 10) {
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

		newTransacao := domain.Transacao{
			ClienteID: id,
			Tipo:      input.Tipo,
			Descricao: input.Descricao,
			Valor:     input.Valor,
		}

		response, err := t.transacaoService.CreateTransaction(c, newTransacao, id)
		if err != nil {
			if err.Error() == ErrNotFound.Error() {
				web.Error(c, http.StatusNotFound, ErrNotFound.Error())
				return
			}
			if err.Error() == LimitErr.Error() {
				web.Error(c, http.StatusUnprocessableEntity, LimitErr.Error())
				return
			}
			web.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		web.Success(c, http.StatusOK, response)
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
