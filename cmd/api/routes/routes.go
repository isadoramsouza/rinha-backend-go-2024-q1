package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/isadoramsouza/rinha-backend-go-2024-q1/cmd/api/handler"
	"github.com/isadoramsouza/rinha-backend-go-2024-q1/internal/transacao"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/rueidis"
)

type Router interface {
	MapRoutes()
}

type router struct {
	eng   *gin.Engine
	rg    *gin.RouterGroup
	db    *pgxpool.Pool
	cache rueidis.Client
}

func NewRouter(eng *gin.Engine, db *pgxpool.Pool, redis rueidis.Client) Router {
	return &router{eng: eng, db: db, cache: redis}
}

func (r *router) MapRoutes() {
	r.setGroup()
	r.buildRoutes()
}

func (r *router) setGroup() {
	r.rg = r.eng.Group("")
}

func (r *router) buildRoutes() {
	repo := transacao.NewRepository(r.db, r.cache)
	service := transacao.NewService(repo)
	handler := handler.NewTransacao(service)
	r.rg.POST("/clientes/:id/transacoes", handler.CreateTransaction())
	r.rg.GET("/clientes/:id/extrato", handler.GetExtrato())
}
