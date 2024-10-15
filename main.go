package main

import (
	"os"

	"github.com/TlexCypher/gsworm/gsworm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"log"
)

type Env struct {
	Driver string
	Dsn    string
}

func loadEnv() *Env {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("Falied to load dotenv.")
	}
	env := &Env{
		Driver: os.Getenv("Driver"),
		Dsn:    os.Getenv("Dsn"),
	}
	return env
}

func main() {
	env := loadEnv()
	g, err := gsworm.Open(&gsworm.Config{
		Driver: env.Dsn,
		Dsn:    env.Driver,
	})

	if err != nil {
		log.Fatalf("Failed to connect with database.\nError:%v\n", err)
	}
	cols := []string{"primary_id", "col1"}
	types := []string{"VARCHAR(255)", "VARCHAR(255)"}

	err = g.DB.Ping()

	if err := g.Create("test1", cols, types); err != nil {
		log.Fatalf("Failed to create table.")
	}
}
