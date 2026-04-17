package main

import (
	"caserunner/pkg/parser"
	"caserunner/pkg/runner"
	"flag"
	"fmt"
	"os"
)

func main() {
	verbose := flag.Bool("verbose", false, "Salida detallada")
	halt := flag.Bool("halt", false, "Pruebas paran cuando se encuetra un error")
	silent := flag.Bool("silent", false, "Salida silenciosa")
	flag.Parse()

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	code := os.Args[len(os.Args)-1]
	file, err := os.ReadFile(os.Args[len(os.Args)-2])
	if err != nil {
		println("No se pudo leer las pruebas")
		os.Exit(1)
	}

	suite, err := parser.ParseFile(string(file), code)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	res := runner.RunSuite(suite, *verbose, *halt, *silent)
	res.PrintResults()
}
