package main

import (
	"os"

	"github.com/TlexCypher/gsworm/internal/gsworm"
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
		Driver: os.Getenv("DRIVER"),
		Dsn:    os.Getenv("DSN"),
	}
	return env
}

func main() {
	env := loadEnv()
	g, err := gsworm.Open(&gsworm.Config{
		Driver: env.Driver,
		Dsn:    env.Dsn,
	})
	defer g.Close()

	if err != nil {
		log.Fatalf("Failed to connect with database.\nError:%v\n", err)
	}
	cols := []string{"primary_id", "col1"}
	types := []gsworm.GswType{gsworm.VARCHAR(255), gsworm.BIGINT()}

	err = g.DB.Ping()
	s := &gsworm.Session{
		ExistTable: make(map[string]bool),
	}

	if err := g.Create("test1", cols, types, s); err != nil {
		log.Fatalf("Failed to create table.")
	}
	if err := g.Drop("test1", s); err != nil {
		log.Fatalf("Failed to drop table.")
	}
}
