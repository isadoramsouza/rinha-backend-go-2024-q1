package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/isadoramsouza/rinha-backend-go-2024-q1/cmd/api/routes"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

func main() {

	// var psqlconn string = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"),
	// 	os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	psqlconn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		"admin",
		"rinha",
		"db",
		"5432",
		"rinhabackenddb",
	)

	poolConfig, err := pgxpool.ParseConfig(psqlconn)
	CheckError(err)

	db, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	CheckError(err)

	err = db.Ping(context.Background())
	CheckError(err)

	defer db.Close()

	fmt.Println("Connected!")

	eng := gin.Default()

	router := routes.NewRouter(eng, db)
	router.MapRoutes()

	if err := eng.Run(); err != nil {
		panic(err)
	}
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
