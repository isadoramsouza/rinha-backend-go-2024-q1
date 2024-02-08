package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Router interface {
	MapRoutes()
}

type router struct {
	eng *gin.Engine
	rg  *gin.RouterGroup
	db  *pgxpool.Pool
}

func NewRouter(eng *gin.Engine, db *pgxpool.Pool) Router {
	return &router{eng: eng, db: db}
}

func (r *router) MapRoutes() {
	r.setGroup()
	r.buildRoutes()
}

func (r *router) setGroup() {
	r.rg = r.eng.Group("")
}

func (r *router) buildRoutes() {
	//repo := pessoa.NewRepository(r.db, r.cache)
	//service := pessoa.NewService(repo)
	//handler := handler.NewPessoa(service)
	r.rg.POST("/clientes/:id/transacoes")
	r.rg.GET("/clientes/:id/extrato")
}
