package main

import (
	"fmt"
	"strconv"

	"github.com/bitfield/script"
)

func execute(tree *ast) error {
	return newExecutor(tree).execute()
}

type executor struct {
	tree *ast
}

func newExecutor(tree *ast) *executor {
	return &executor{tree: tree}
}

func (e executor) execute() error {
	s := script.Stdin()

	for _, f := range e.tree.functions {
		switch f.name {
		case "First":
			if err := e.argumentsMustLen(f, 1); err != nil {
				return err
			}
			i, err := e.argumentAtoi(f, 0)
			if err != nil {
				return err
			}
			s = s.First(i)
		case "Freq":
			if err := e.argumentsMustLen(f, 0); err != nil {
				return err
			}
			s = s.Freq()
		case "Replace":
			if err := e.argumentsMustLen(f, 2); err != nil {
				return err
			}
			s = s.Replace(f.arguments[0], f.arguments[1])
		default:
			return fmt.Errorf("unknown function %s()", f.name)
		}
	}

	s.Stdout()

	return nil
}

func (e executor) argumentsMustLen(f function, count int) error {
	if len(f.arguments) != count {
		return fmt.Errorf("wrong number of arguments in call to %s(): got %d, expected %d", f.name, len(f.arguments), count)
	}

	return nil
}

func (e executor) argumentAtoi(f function, argument int) (int, error) {
	i, err := strconv.Atoi(f.arguments[argument])
	if err != nil {
		return 0, fmt.Errorf("expected argument in call to %s() to be int: %v", f.name, err)
	}

	return i, nil
}
