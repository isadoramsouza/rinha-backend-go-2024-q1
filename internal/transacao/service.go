package transacao

import (
	"context"

	"github.com/isadoramsouza/rinha-backend-go-2024-q1/internal/domain"
)

type Service interface {
	CreateTransaction(ctx context.Context, t domain.Transacao) (domain.TransacaoResponse, error)
	GetBalance(ctx context.Context, id int) (domain.Cliente, error)
	GetExtrato(ctx context.Context, id int) (domain.Extrato, error)
}

type transacaoService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &transacaoService{
		repository: r,
	}
}

func (s *transacaoService) CreateTransaction(ctx context.Context, t domain.Transacao) (domain.TransacaoResponse, error) {
	response, err := s.repository.SaveTransaction(ctx, t)

	return response, err
}

func (s *transacaoService) GetBalance(ctx context.Context, id int) (domain.Cliente, error) {
	cliente, err := s.repository.GetBalance(ctx, id)

	return cliente, err
}

func (s *transacaoService) GetExtrato(ctx context.Context, id int) (domain.Extrato, error) {
	extrato, err := s.repository.GetExtrato(ctx, id)

	return extrato, err
}
