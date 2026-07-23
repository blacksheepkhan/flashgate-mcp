package benchmark

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

var benchmarkCommitPattern = regexp.MustCompile(`^(unknown|[0-9a-f]{7,40})$`)

func loadValidatedBaselineArtifact(artifactPath string, budgetPath string) (Result, []byte, error) {
	artifactName := filepath.Base(artifactPath)
	raw, err := os.ReadFile(artifactPath)
	if err != nil {
		return Result{}, nil, benchmarkReadError("artifact", artifactPath, err)
	}
	var result Result
	if err := decodeStrictJSON(raw, &result); err != nil {
		return Result{}, nil, fmt.Errorf("artifact %s JSON: %w", artifactName, err)
	}
	if err := validateBaselineSchemaInvariants(result); err != nil {
		return Result{}, nil, fmt.Errorf("artifact %s schema: %w", artifactName, err)
	}
	if err := validateArtifactBudgets(artifactName, budgetPath, result); err != nil {
		return Result{}, nil, err
	}
	return result, raw, nil
}

func benchmarkReadError(area string, path string, err error) error {
	cause := "read failed"
	if errors.Is(err, os.ErrNotExist) {
		cause = "file does not exist"
	} else if errors.Is(err, os.ErrPermission) {
		cause = "permission denied"
	}
	return fmt.Errorf("%s %s: %s", area, filepath.Base(path), cause)
}

func validateArtifactBudgets(artifactName string, budgetPath string, result Result) error {
	recomputed, err := EvaluateBudgets(budgetPath, result)
	if err != nil {
		return fmt.Errorf("artifact %s canonical budget evaluation: %w", filepath.Base(artifactName), err)
	}
	matches := reflect.DeepEqual(result.BudgetEvaluation, recomputed)
	if recomputed.HardFailures > 0 {
		return fmt.Errorf(
			"artifact %s budget failure: recomputed %d hard failure(s): %s; embedded budget_evaluation matches=%t",
			filepath.Base(artifactName), recomputed.HardFailures, strings.Join(recomputed.Messages, "; "), matches,
		)
	}
	if !matches {
		return fmt.Errorf(
			"artifact %s budget_evaluation mismatch: embedded hard=%d soft=%d messages=%q; recomputed hard=%d soft=%d messages=%q",
			filepath.Base(artifactName),
			result.BudgetEvaluation.HardFailures, result.BudgetEvaluation.SoftWarnings, result.BudgetEvaluation.Messages,
			recomputed.HardFailures, recomputed.SoftWarnings, recomputed.Messages,
		)
	}
	return nil
}

func validateBaselineSchemaInvariants(result Result) error {
	if result.SchemaVersion != ResultSchemaVersion {
		return fmt.Errorf("schema_version=%q, want %q", result.SchemaVersion, ResultSchemaVersion)
	}
	if result.BenchmarkSuiteVersion != BenchmarkSuiteVersion || result.WorkflowCatalogVersion != WorkflowCatalogVersion || result.CorpusVersion != CorpusVersion {
		return fmt.Errorf("benchmark version provenance is invalid")
	}
	if result.Project != "flashgate-mcp" {
		return fmt.Errorf("project=%q, want flashgate-mcp", result.Project)
	}
	if !benchmarkCommitPattern.MatchString(result.Commit) {
		return fmt.Errorf("commit does not match the benchmark schema")
	}
	if result.GeneratedAtUTC.IsZero() {
		return fmt.Errorf("generated_at_utc must be a non-zero date-time")
	}
	if result.GoVersion == "" || result.OS == "" || result.Architecture == "" {
		return fmt.Errorf("runtime platform provenance is incomplete")
	}
	if result.RuntimeMode != "direct_stdio" || result.Transport != "stdio" || result.ExecutionBackend != "current_process" || result.Profile != "read_only" {
		return fmt.Errorf("execution provenance is invalid")
	}
	if result.Parallelism != 1 || result.Repetitions < 1 {
		return fmt.Errorf("parallelism or repetitions violate the schema")
	}
	if len(result.ToolsList) < 2 {
		return fmt.Errorf("tools_list_measurements has %d item(s), want at least 2", len(result.ToolsList))
	}
	if len(result.Workflows) < 10 {
		return fmt.Errorf("workflow_measurements has %d item(s), want at least 10", len(result.Workflows))
	}
	if err := validateStartMeasurementSchema("start_measurements.first_process_start", result.StartMeasurements.FirstProcessStart); err != nil {
		return err
	}
	if err := validateStartMeasurementSchema("start_measurements.subsequent_process_start", result.StartMeasurements.SubsequentProcessStart); err != nil {
		return err
	}
	if err := validateResourceSchema("resource_measurements", result.Resources); err != nil {
		return err
	}
	for index, measurement := range result.ToolsList {
		if measurement.Profile != "read_only" && measurement.Profile != "default" {
			return fmt.Errorf("tools_list_measurements[%d].profile=%q is invalid", index, measurement.Profile)
		}
		if measurement.ToolCount < 0 || measurement.SchemaCount < 0 {
			return fmt.Errorf("tools_list_measurements[%d] contains a negative count", index)
		}
	}
	for index, measurement := range result.Workflows {
		name := fmt.Sprintf("workflow_measurements[%d]", index)
		if measurement.Name == "" || measurement.Repetitions < 1 || measurement.Calls < 0 || measurement.WarningCount < 0 {
			return fmt.Errorf("%s identity or counts violate the schema", name)
		}
		for _, metric := range []struct {
			name    string
			summary MetricSummary
		}{
			{"request_bytes", measurement.RequestBytes}, {"response_bytes", measurement.ResponseBytes},
			{"result_bytes", measurement.ResultBytes}, {"duration_ns", measurement.DurationNS},
			{"read_bytes", measurement.ReadBytes}, {"written_bytes", measurement.WrittenBytes},
			{"scanned_bytes", measurement.ScannedBytes}, {"entries", measurement.Entries},
			{"approx_tokens_bytes4", measurement.ApproxTokensBytes4},
		} {
			if err := validateMetricSchema(name+"."+metric.name, metric.summary); err != nil {
				return err
			}
		}
		if err := validateResourceSchema(name+".resources", measurement.Resources); err != nil {
			return err
		}
	}
	if result.BudgetEvaluation.SchemaVersion != BudgetSchemaVersion || result.BudgetEvaluation.HardFailures < 0 || result.BudgetEvaluation.SoftWarnings < 0 {
		return fmt.Errorf("budget_evaluation violates the schema")
	}
	return nil
}

func validateStartMeasurementSchema(name string, measurement StartMeasurement) error {
	if measurement.Repetitions < 1 || measurement.WarningCount < 0 {
		return fmt.Errorf("%s counts violate the schema", name)
	}
	for _, metric := range []struct {
		name    string
		summary MetricSummary
	}{
		{"duration_ns", measurement.DurationNS}, {"request_bytes", measurement.RequestBytes},
		{"response_bytes", measurement.ResponseBytes}, {"result_bytes", measurement.ResultBytes},
	} {
		if err := validateMetricSchema(name+"."+metric.name, metric.summary); err != nil {
			return err
		}
	}
	for code, count := range measurement.ExitStatuses {
		if count < 0 {
			return fmt.Errorf("%s.exit_statuses.%s is negative", name, code)
		}
	}
	return validateResourceSchema(name+".resources", measurement.Resources)
}

func validateResourceSchema(name string, resources ResourceSummary) error {
	if resources.Status != "supported" && resources.Status != "not_supported" {
		return fmt.Errorf("%s.status=%q is invalid", name, resources.Status)
	}
	for _, metric := range []struct {
		name    string
		summary *MetricSummary
	}{
		{"idle_working_set_bytes", resources.IdleWorkingSetBytes},
		{"peak_working_set_bytes", resources.PeakWorkingSetBytes},
		{"user_cpu_ns", resources.UserCPUNS},
		{"system_cpu_ns", resources.SystemCPUNS},
	} {
		if metric.summary != nil {
			if err := validateMetricSchema(name+"."+metric.name, *metric.summary); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateMetricSchema(name string, summary MetricSummary) error {
	if summary.Samples < 0 {
		return fmt.Errorf("%s.samples is negative", name)
	}
	return nil
}
