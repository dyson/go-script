package main

type ast struct {
	functions []function
}

func newAST() *ast {
	return &ast{
		functions: []function{},
	}
}

type function struct {
	name      string
	arguments []string
}
