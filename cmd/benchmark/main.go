package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"

	benchmarkrunner "github.com/thomasweidner/flashgate-mcp/internal/benchmark"
)

var commitPattern = regexp.MustCompile(`^(unknown|[0-9a-f]{7,40})$`)

func main() {
	os.Exit(run())
}

func run() int {
	var binaryPath string
	var outputPath string
	var budgetPath string
	var protectedBaselineDirectory string
	var commit string
	var repetitions int
	var quick bool
	var workingTreeDirty bool

	flag.StringVar(&binaryPath, "binary", "", "path to the built FlashGate server binary")
	flag.StringVar(&outputPath, "output", "benchmark-result.json", "machine-readable result path")
	flag.StringVar(&budgetPath, "budgets", "", "optional machine-readable regression budget path")
	flag.StringVar(&protectedBaselineDirectory, "protected-baseline-dir", "", "protected repository baseline directory")
	flag.StringVar(&commit, "commit", "unknown", "source commit recorded in the result")
	flag.IntVar(&repetitions, "repetitions", 30, "number of subsequent starts and workflow repetitions")
	flag.BoolVar(&quick, "quick", false, "use ten subsequent starts and workflow repetitions")
	flag.BoolVar(&workingTreeDirty, "working-tree-dirty", false, "record that the binary was built from a dirty working tree")
	flag.Parse()

	if quick {
		repetitions = 10
	}
	if !commitPattern.MatchString(commit) {
		fmt.Fprintln(os.Stderr, "benchmark: commit must be a hexadecimal Git object ID or unknown")
		return 2
	}

	result, err := benchmarkrunner.Run(context.Background(), benchmarkrunner.Options{
		BinaryPath:       binaryPath,
		Commit:           commit,
		WorkingTreeDirty: workingTreeDirty,
		Repetitions:      repetitions,
		BudgetPath:       budgetPath,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "benchmark: %v\n", err)
		return 1
	}

	encoded, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "benchmark: encode result: %v\n", err)
		return 1
	}
	encoded = append(encoded, '\n')
	if err := benchmarkrunner.WriteResultFile(outputPath, protectedBaselineDirectory, encoded); err != nil {
		fmt.Fprintf(os.Stderr, "benchmark: write result: %v\n", err)
		return 1
	}

	if result.BudgetEvaluation.HardFailures > 0 {
		fmt.Fprintf(os.Stderr, "benchmark: %d hard budget failure(s)\n", result.BudgetEvaluation.HardFailures)
		return 3
	}
	fmt.Printf("benchmark complete: repetitions=%d warnings=%d\n", result.Repetitions, len(result.Warnings))
	return 0
}
