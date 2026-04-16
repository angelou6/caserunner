package runner

import "caserunner/pkg/parser"

type Outcome struct {
	success, fail, error int
}

func RunCase(exec string, cases []parser.TestCase) Outcome {

}

func RunSuite(suite parser.TestSuite) Outcome {

}
