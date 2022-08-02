package main

import (
	"fmt"
	"os"
)

func main() {
	s, err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}

	os.Exit(s)
}

func run() (int, error) {
	if len(os.Args) != 2 {
		return 1, fmt.Errorf("use a single arg to define pipeline")
	}

	tree, err := parse(os.Args[1])
	if err != nil {
		return 1, fmt.Errorf("error parsing pipeline: %w", err)
	}

	//fmt.Println(tree)
	err = execute(tree)
	if err != nil {
		return 1, fmt.Errorf("error executing pipeline: %w", err)
	}

	return 0, nil
}
