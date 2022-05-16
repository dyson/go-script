package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	s, err := run()
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	os.Exit(s)
}

func run() (int, error) {
	dir, err := ioutil.TempDir("", "go-script-*")
	if err != nil {
		return 1, err
	}
	defer os.Remove(dir)

	// go mod/sum are now required
	templateMod := `module github.com/dyson/go-script-temp

go 1.18

require github.com/bitfield/script v0.20.0

require (
	bitbucket.org/creachadair/shell v0.0.7 // indirect
	github.com/itchyny/gojq v0.12.7 // indirect
	github.com/itchyny/timefmt-go v0.1.3 // indirect
)`

	templateSum := `bitbucket.org/creachadair/shell v0.0.7 h1:Z96pB6DkSb7F3Y3BBnJeOZH2gazyMTWlvecSD4vDqfk=
bitbucket.org/creachadair/shell v0.0.7/go.mod h1:oqtXSSvSYr4624lnnabXHaBsYW6RD80caLi2b3hJk0U=
github.com/bitfield/script v0.20.0 h1:dqeNh8LKf3MfKN21fpuuXNNZPzLlrG6WnmScHqOt13I=
github.com/bitfield/script v0.20.0/go.mod h1:l3AZPVAtKQrL03bwh7nlNTUtgrgSWurpJSbtqspYrOA=
github.com/google/go-cmp v0.5.4/go.mod h1:v8dTdLbMG2kIc/vJvl+f65V22dbkXbowE6jgT/gNBxE=
github.com/google/go-cmp v0.5.6/go.mod h1:v8dTdLbMG2kIc/vJvl+f65V22dbkXbowE6jgT/gNBxE=
github.com/google/go-cmp v0.5.7 h1:81/ik6ipDQS2aGcBfIN5dHDB36BwrStyeAQquSYCV4o=
github.com/google/go-cmp v0.5.7/go.mod h1:n+brtR0CgQNWTVd5ZUFpTBC8YFBDLK/h/bpaJ8/DtOE=
github.com/itchyny/gojq v0.12.7 h1:hYPTpeWfrJ1OT+2j6cvBScbhl0TkdwGM4bc66onUSOQ=
github.com/itchyny/gojq v0.12.7/go.mod h1:ZdvNHVlzPgUf8pgjnuDTmGfHA/21KoutQUJ3An/xNuw=
github.com/itchyny/timefmt-go v0.1.3 h1:7M3LGVDsqcd0VZH2U+x393obrzZisp7C0uEe921iRkU=
github.com/itchyny/timefmt-go v0.1.3/go.mod h1:0osSSCQSASBJMsIZnhAaF1C2fCBTJZXrnj37mG8/c+A=
github.com/mattn/go-isatty v0.0.14/go.mod h1:7GGIvUiUoEMVVmxf/4nioHXj79iQHKdU27kJ6hsGG94=
github.com/mattn/go-runewidth v0.0.9/go.mod h1:H031xJmbD/WCDINGzjvQ9THkh0rPKHF+m2gUSrubnMI=
golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
golang.org/x/sys v0.0.0-20220227234510-4e6760a101f9/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543/go.mod h1:I/5z698sn9Ka8TeJc9MKroUUfqBBauWjQqLJ2OPfmY0=
golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1/go.mod h1:I/5z698sn9Ka8TeJc9MKroUUfqBBauWjQqLJ2OPfmY0=
gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=`

	templateGo := `package main

import "github.com/bitfield/script"

func main() {
	script.Stdin()%s.Stdout()
}`

	// handle piping stdin -> stdout if no args provided (passthrough)
	input := ""
	switch {
	case len(os.Args) == 2:
		if i := os.Args[1]; i != "" {
			input = fmt.Sprintf(".%s", os.Args[1])
		}
	case len(os.Args) > 2:
		return 1, fmt.Errorf("use a single arg to define pipeline")
	}

	// create mod file
	if err := ioutil.WriteFile(filepath.Join(dir, "go.mod"), []byte(templateMod), 0644); err != nil {
		return 1, err
	}
	defer os.Remove(filepath.Join(dir, "go.mod"))

	// create sum file
	if err := ioutil.WriteFile(filepath.Join(dir, "go.sum"), []byte(templateSum), 0644); err != nil {
		return 1, err
	}
	defer os.Remove(filepath.Join(dir, "go.sum"))

	// create program file
	if err := ioutil.WriteFile(filepath.Join(dir, "main.go"), []byte(fmt.Sprintf(templateGo, input)), 0644); err != nil {
		return 1, err
	}
	defer os.Remove(filepath.Join(dir, "main.go"))

	// run
	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode(), err
		}
	}

	return 0, nil
}
