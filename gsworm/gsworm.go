package gsworm

import (
	"database/sql"
	"fmt"
	"log"
)

type Gsworm struct {
	DB *sql.DB
}

type GswColType interface {
	Column()
}

type Number struct{}

func (n *Number) Column() {}

type String struct{}

func (s *String) Column() {}

func Open(c *Config) (*Gsworm, error) {
	db, err := sql.Open(c.Driver, c.Dsn)
	if err != nil {
		return nil, err
	}
	return &Gsworm{DB: db}, nil
}

func (g *Gsworm) Close() error {
	if err := g.DB.Close(); err != nil {
		log.Printf("Failed to close database connection. Error:%v\n", err)
		return err
	}
	return nil
}

func (g *Gsworm) Create(table string, cols []string, types []string) error {
	cldc, err := g.genColsDeclaration(table, cols, types)
	if err != nil {
		log.Fatalf("Number of column and those of types should be equal. Columns:%v, Types:%v\n", len(cols), len(types))
	}
	create := fmt.Sprintf("CREATE TABLE %v (%v);\n", table, cldc)
	_, err = g.DB.Exec(create)
	if err != nil {
		log.Printf("Failed to execute create table query. Published raw query is %v\nError:%v\n", create, err)
		/* TODO: Wrap table declaration error */
		return err
	}
	return nil
}

func (g *Gsworm) genColsDeclaration(table string, cols []string, types []string) (string, error) {
	if len(cols) != len(types) {
		return "", &GswTblDeclarationErr{Table: table}
	}
	cldc := ""
	for i, col := range cols {
		if len(cldc) != 0 {
			cldc += ", "
		}
		cldc += col
		cldc += " "
		cldc += types[i]
	}
	return cldc, nil
}
