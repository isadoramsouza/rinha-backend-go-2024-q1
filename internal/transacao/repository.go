package transacao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/isadoramsouza/rinha-backend-go-2024-q1/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/rueidis"
)

var (
	ErrDuplicateApelido = errors.New("duplicate apelido")
	ErrNotFound         = errors.New("cliente not found")
	LimitErr            = errors.New("limit error")
)

type Repository interface {
	SaveTransaction(ctx context.Context, t domain.Transacao) (domain.TransacaoResponse, error)
	GetBalance(ctx context.Context, id int) (domain.Cliente, error)
	GetExtrato(ctx context.Context, id int) (domain.Extrato, error)
}

type repository struct {
	db    *pgxpool.Pool
	cache rueidis.Client
}

func NewRepository(db *pgxpool.Pool, redis rueidis.Client) Repository {
	return &repository{
		db:    db,
		cache: redis,
	}
}

func (r *repository) SaveTransaction(ctx context.Context, t domain.Transacao) (domain.TransacaoResponse, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return domain.TransacaoResponse{}, err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		"INSERT INTO transacoes(tipo, descricao, valor, cliente_id) VALUES ($1, $2, $3, $4)",
		t.Tipo, t.Descricao, t.Valor, t.ClienteID)
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Message == "Débito excede o limite do cliente" {
			return domain.TransacaoResponse{}, LimitErr
		}
		return domain.TransacaoResponse{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.TransacaoResponse{}, err
	}

	cliente, _ := r.GetBalance(ctx, t.ClienteID)
	if err != nil {
		return domain.TransacaoResponse{}, err
	}
	query := `SELECT valor, tipo, descricao, realizada_em FROM transacoes t where cliente_id = $1 ORDER BY realizada_em DESC LIMIT 10;`
	rows, err := r.db.Query(ctx, query, t.ClienteID)
	if err != nil {
		return domain.TransacaoResponse{}, err
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
	go r.SaveExtratoInCache(t.ClienteID, extrato)

	response := domain.TransacaoResponse{
		Saldo:  cliente.Saldo,
		Limite: cliente.Limite,
	}

	return response, nil
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
	extratoInCache, _ := r.GetExtratoInCache(id)
	fmt.Println(extratoInCache)
	if extratoInCache.Saldo.Limite != 0 {
		fmt.Println("CACHEEEEE")
		return extratoInCache, nil
	}
	fmt.Println("NOOO CACHEEEEE")

	cliente, err := r.GetBalance(ctx, id)
	if err != nil {
		return domain.Extrato{}, err
	}
	query := `SELECT valor, tipo, descricao, realizada_em FROM transacoes t where cliente_id = $1 ORDER BY realizada_em DESC LIMIT 10;`
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

func (r *repository) SaveExtratoInCache(clienteID int, extrato domain.Extrato) error {
	ctx := context.Background()
	extratoString, err := sonic.MarshalString(extrato)
	fmt.Println(extratoString)
	if err != nil {
		fmt.Println("ERROOO CACHE")

		return err
	}
	fmt.Println("SALVANDO CACHE")
	extratoInCached := r.cache.B().Set().Key("extrato-id-" + fmt.Sprint(clienteID)).Value(extratoString).Build()
	for _, resp := range r.cache.DoMulti(ctx, extratoInCached) {
		if err := resp.Error(); err != nil {
			fmt.Println("ERROOO CACH2E")

			return err
		}
	}
	return nil
}

func (r *repository) GetExtratoInCache(clienteID int) (domain.Extrato, error) {
	fmt.Println(fmt.Sprint(clienteID))
	ctx := context.Background()
	extratoResult, err := r.cache.Do(ctx, r.cache.B().Get().Key("extrato-id-"+fmt.Sprint(clienteID)).Build()).ToString()
	if err != nil {
		fmt.Println(err)
		return domain.Extrato{}, err
	}
	fmt.Println(extratoResult)
	var extrato domain.Extrato
	err = json.Unmarshal([]byte(extratoResult), &extrato)
	if err != nil {
		fmt.Println("ERROOO qqqGETCACH2E")

		return domain.Extrato{}, err
	}
	return extrato, nil
}
