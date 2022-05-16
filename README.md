# go-script

Run https://github.com/bitfield/script directly from the CLI.

## Installation

```bash
$ go install github.com/dyson/go-script@latest
```

Optionally alias go-script, eg:

```bash
echo 'alias gs="go-script"' >> ~/.bash_profile
```

## Example

go-script wraps its input between `script.Stdin.` and `.Stdout` and so works with all `bitfield/script` functions that return a pipeline.

A contrived example:

```bash
$ echo "cat cat cat dog bird bird bird bird" | gs 'Replace(" ", "\n").Freq().First(1)'
4 bird
```
