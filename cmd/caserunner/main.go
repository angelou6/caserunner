package main

import (
	"caserunner/internal/parser"
	"caserunner/internal/runner"
	"flag"
	"fmt"
	"os"
)

func main() {
	verbose := flag.Bool("verbose", false, "Salida detallada")
	halt := flag.Bool("halt", false, "El programa se detiene cuando encuentra un error")
	flag.Parse()

	code := os.Args[len(os.Args)-1]
	testfile := os.Args[len(os.Args)-2]

	data, err := os.ReadFile(testfile)
	if err != nil {
		fmt.Println("Error leyendo archivo:", err)
		return
	}

	par := parser.New()
	err = par.ParseFile(string(data), code)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	runner.RunFile(par, *verbose, *halt)
}
