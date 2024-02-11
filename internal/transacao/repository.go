package transacao

import (
	"context"
	"errors"
	"time"

	"github.com/isadoramsouza/rinha-backend-go-2024-q1/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrDuplicateApelido = errors.New("duplicate apelido")
	ErrNotFound         = errors.New("cliente not found")
)

type Repository interface {
	SaveTransaction(ctx context.Context, t domain.Transacao, saldo int64) error
	GetBalance(ctx context.Context, id int) (domain.Cliente, error)
	GetExtrato(ctx context.Context, id int) (domain.Extrato, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) SaveTransaction(ctx context.Context, t domain.Transacao, saldo int64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // Rollback da transação em caso de erro

	// Inserção da transação
	_, err = tx.Exec(ctx,
		"INSERT INTO public.transacoes(tipo, descricao, valor, cliente_id) VALUES ($1, $2, $3, $4)",
		t.Tipo, t.Descricao, t.Valor, t.ClienteID)
	if err != nil {
		return err
	}

	// Atualização do saldo do cliente
	_, err = tx.Exec(ctx,
		"UPDATE public.clientes SET saldo = $1 WHERE id = $2",
		saldo, t.ClienteID)
	if err != nil {
		return err
	}

	// Commit da transação
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetBalance(ctx context.Context, id int) (domain.Cliente, error) {
	query := "SELECT limite, saldo FROM clientes WHERE id=$1;"
	row := r.db.QueryRow(ctx, query, id)
	c := domain.Cliente{}
	var limite, saldo int64
	err := row.Scan(&limite, &saldo)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return domain.Cliente{}, ErrNotFound
		}
		return domain.Cliente{}, err
	}
	c.Limite = limite
	c.Saldo = saldo
	return c, nil
}

func (r *repository) GetExtrato(ctx context.Context, id int) (domain.Extrato, error) {
	cliente, err := r.GetBalance(ctx, id)
	if err != nil {
		return domain.Extrato{}, err
	}
	query := `SELECT valor, tipo, descricao, realizada_em FROM transacoes t where cliente_id = $1 LIMIT 10;`
	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return domain.Extrato{}, ErrNotFound
		}
		return domain.Extrato{}, err
	}

	var transacoes []domain.UltimaTransacao

	for rows.Next() {
		t := domain.UltimaTransacao{}
		_ = rows.Scan(&t.Valor, &t.Tipo, &t.Descricao, &t.RealizadaEm)
		transacoes = append(transacoes, t)
	}

	extrato := domain.Extrato{
		Saldo: domain.Saldo{
			Total:       cliente.Saldo,
			DataExtrato: time.Now().UTC(),
			Limite:      cliente.Limite,
		},
		UltimasTransacoes: transacoes,
	}

	return extrato, nil
}
