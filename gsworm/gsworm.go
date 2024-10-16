package gsworm

import (
	"database/sql"
	"fmt"
	"log"
)

type Gsworm struct {
	DB *sql.DB
}

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

func (g *Gsworm) Create(table string, cols []string, types []GswType, s *Session) error {
	cldc, err := g.genColsDeclaration(table, cols, types)
	if err != nil {
		log.Fatalf("Number of column and those of types should be equal. Columns:%v, Types:%v\n", len(cols), len(types))
	}
	create := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v (%v);\n", table, cldc)
	_, err = g.DB.Exec(create)
	if err != nil {
		log.Printf("Failed to execute create table query. Published raw query is %v\nError:%v\n", create, err)
		/* TODO: Wrap table declaration error */
		return err
	}
	log.Printf("Success to create table %v\n", table)
	s.ExistTable[table] = true

	return nil
}

func (g *Gsworm) genColsDeclaration(table string, cols []string, types []GswType) (string, error) {
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
		switch t := types[i].(type) {
		case VChar:
			cldc += string(t.SqlType())
		case Int:
			cldc += string(t.SqlType())
		case BInt:
			cldc += string(t.SqlType())
		default:
			log.Fatalf("Unknown type.")
		}
	}
	return cldc, nil
}

func (g *Gsworm) Drop(table string, s *Session) error {
	d := fmt.Sprintf("DROP TABLE %v\n", table)
	if _, err := g.DB.Exec(d); err != nil {
		log.Printf("Failed to drop table. Table:%v\n", table)
		return err
	}
	log.Printf("Success to drop table. Table:%v\n", table)
	s.ExistTable[table] = false

	return nil
}

func (g *Gsworm) Insert(table string, columns []string, values []any, s *Session) error {
	if len(columns) != len(values) {
		log.Fatalf("The number of columns and those of values are not equal when you call insert method.")
	}
	if !s.ExistTable[table] {
		log.Fatalf("%v table does not exist.\n", table)
	}
	insert := g.genInsertStatement(table, columns, values)
	log.Printf("Success to insert records into %v\n", table)
	return nil
}

func (g *Gsworm) genInsertStatement(table string, columns []string, values []any) string {
	//INSERT INTO ${table} (...columns) (...values)
	for i := range len(columns) {
		col, val := columns[i], values[i]
		//colの型とvalueの型が違う場合, errorが必要 --> tableのメタ情報を持っとく必要がある
	}
}
