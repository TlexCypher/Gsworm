package gsworm

import "fmt"

type GswErr struct{}

type GswTblDeclarationErr struct {
	Table string
}

func (ge *GswTblDeclarationErr) Error() string {
	return fmt.Sprintf("Failed to declare %v\n", ge.Table)
}
