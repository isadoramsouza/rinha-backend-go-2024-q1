package transacao

import (
	"context"

	"github.com/isadoramsouza/rinha-backend-go-2024-q1/internal/domain"
)

type Service interface {
	CreateTransaction(ctx context.Context, t domain.Transacao) (domain.TransacaoResponse, error)
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

var semaphore = make(chan struct{}, 20)

func (s *transacaoService) CreateTransaction(ctx context.Context, t domain.Transacao) (domain.TransacaoResponse, error) {
	responseChan := make(chan domain.TransacaoResponse, 1)
	errChan := make(chan error, 1)

	semaphore <- struct{}{}

	go func() {
		defer func() {
			<-semaphore
		}()

		response, err := s.repository.SaveTransaction(ctx, t)
		if err != nil {
			errChan <- err
			return
		}
		responseChan <- response
	}()

	select {
	case response := <-responseChan:
		return response, nil
	case err := <-errChan:
		return domain.TransacaoResponse{}, err
	case <-ctx.Done():
		return domain.TransacaoResponse{}, ctx.Err()
	}
}

func (s *transacaoService) GetExtrato(ctx context.Context, id int) (domain.Extrato, error) {
	extratoChan := make(chan domain.Extrato, 1)
	errChan := make(chan error, 1)

	semaphore <- struct{}{}

	go func() {
		defer func() {
			<-semaphore
		}()

		extrato, err := s.repository.GetExtrato(ctx, id)
		if err != nil {
			errChan <- err
			return
		}
		extratoChan <- extrato
	}()

	select {
	case extrato := <-extratoChan:
		return extrato, nil
	case err := <-errChan:
		return domain.Extrato{}, err
	case <-ctx.Done():
		return domain.Extrato{}, ctx.Err()
	}
}
