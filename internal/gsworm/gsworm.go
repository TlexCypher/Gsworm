package gsworm

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strconv"
)

type Schema struct {
	columns []string
	types   []GswType
}

type Gsworm struct {
	DB      *sql.DB
	Schemas map[string]Schema
}

func Open(c *Config) (*Gsworm, error) {
	db, err := sql.Open(c.Driver, c.Dsn)
	if err != nil {
		return nil, err
	}
	return &Gsworm{DB: db, Schemas: make(map[string]Schema)}, nil
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
	g.Schemas[table] = Schema{
		columns: cols,
		types:   types,
	}
	return nil
}

func (g *Gsworm) genColsDeclaration(table string, cols []string, types []GswType) (string, error) {
	if len(cols) != len(types) {
		return "", GswTblDeclarationErr{Table: table}
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
	delete(g.Schemas, table)
	return nil
}

func (g *Gsworm) Insert(table string, columns []string, values []GswType, s *Session) error {
	if len(columns) != len(values) {
		log.Fatalf("The number of columns and those of values are not equal when you call insert method.")
	}
	if !s.ExistTable[table] {
		log.Fatalf("%v table does not exist.\n", table)
	}
	if !g.typeAssertion(columns, values) {
		log.Fatalf("Failed to type assertion. Some value's type might not be suitable.")
	}
	insert, err := g.genInsertStatement(table, columns, values)
	if err != nil {
		return err
	}
	_, err = g.DB.Exec(insert)
	if err != nil {
		log.Fatalf("Failed to insert records into %v\n", table)
	}
	log.Printf("Success to insert records into %v\n", table)
	return nil
}

func (g *Gsworm) typeAssertion(columns []string, values []GswType) bool {
	s := g.Schemas
	for i, col := range columns {
		if reflect.TypeOf(col) != reflect.TypeOf(values[i]) {
			log.Printf("Failed to do type assertion. column:%v, value:%v, expected type:%v actual type:%v\n",
				col, values[i], reflect.TypeOf(s[col]), reflect.TypeOf(values[i]))
			return false
		}
	}
	return true
}

func (g *Gsworm) genInsertStatement(table string, columns []string, values []GswType) (string, error) {
	strValues := make([]string, len(values))
	radix := 10
	for i, v := range values {
		switch vt := v.(type) {
		case Int:
			log.Printf("Column %v's type is %v\n", columns[i], vt)
			strValues[i] = strconv.Itoa(int(*v.(Int).Value))
			break
		case BInt:
			log.Printf("Column %v's type is %v\n", columns[i], vt)
			strValues[i] = strconv.FormatInt(*(v.(BInt).Value), radix)
			break
		case VChar:
			log.Printf("Column %v's type is %v\n", columns[i], vt)
			strValues[i] = *v.(VChar).Value
			break
		default:
			log.Fatalf("Failed to parse column")
		}
	}
	var insert bytes.Buffer
	allCols, allVals := getColsVals(columns, strValues)
	if len(allCols) != len(allVals) {
		log.Fatalln("The number of values and those of columns should be equal.")
	}
	insert.Write([]byte(fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)", table, allCols, allVals)))
	return insert.String(), nil
}

func getColsVals(columns []string, values []string) (string, string) {
	var cols, vals bytes.Buffer
	for i, col := range columns {
		cols.WriteString(col)
		if i != len(columns) {
			cols.WriteString(", ")
		}
	}

	for i, val := range values {
		vals.WriteString(val)
		if i != len(values) {
			vals.WriteString(", ")
		}
	}
	return cols.String(), vals.String()
}
