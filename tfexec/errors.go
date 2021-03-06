package tfexec

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

var (
	// The "Required variable not set:" case is for 0.11
	missingVarErrRegexp  = regexp.MustCompile(`Error: No value for required variable|Error: Required variable not set:`)
	missingVarNameRegexp = regexp.MustCompile(`The root module input variable "(.+)" is not set, and has no default|Error: Required variable not set: (.+)`)

	usageRegexp = regexp.MustCompile(`Too many command line arguments|^Usage: .*Options:.*|Error: Invalid -\d+ option`)

	// "Could not load plugin" is present in 0.13
	noInitErrRegexp = regexp.MustCompile(`Error: Could not satisfy plugin requirements|Error: Could not load plugin`)

	noConfigErrRegexp = regexp.MustCompile(`Error: No configuration files`)
)

func parseError(err error, stderr string) error {
	if _, ok := err.(*exec.ExitError); !ok {
		return err
	}

	switch {
	case missingVarErrRegexp.MatchString(stderr):
		name := ""
		names := missingVarNameRegexp.FindStringSubmatch(stderr)
		for i := 1; i < len(names); i++ {
			name = strings.TrimSpace(names[i])
			if name != "" {
				break
			}
		}
		
		return &ErrMissingVar{name}
	case usageRegexp.MatchString(stderr):
		return &ErrCLIUsage{stderr: stderr}
	case noInitErrRegexp.MatchString(stderr):
		return &ErrNoInit{stderr: stderr}
	case noConfigErrRegexp.MatchString(stderr):
		return &ErrNoConfig{stderr: stderr}
	default:
		return errors.New(stderr)
	}
}

type ErrNoSuitableBinary struct {
	err error
}

func (e *ErrNoSuitableBinary) Error() string {
	return fmt.Sprintf("no suitable terraform binary could be found: %s", e.err.Error())
}

// ErrVersionMismatch is returned when the detected Terraform version is not compatible with the
// command or flags being used in this invocation.
type ErrVersionMismatch struct {
	MinInclusive string
	MaxExclusive string
	Actual       string
}

func (e *ErrVersionMismatch) Error() string {
	return fmt.Sprintf("unexpected version %s (min: %s, max: %s)", e.Actual, e.MinInclusive, e.MaxExclusive)
}

type ErrNoInit struct {
	stderr string
}

func (e *ErrNoInit) Error() string {
	return e.stderr
}

type ErrNoConfig struct {
	stderr string
}

func (e *ErrNoConfig) Error() string {
	return e.stderr
}

// ErrCLIUsage is returned when the combination of flags or arguments is incorrect.
//
//  CLI indicates usage errors in three different ways: either
// 1. Exit 1, with a custom error message on stderr.
// 2. Exit 1, with command usage logged to stderr.
// 3. Exit 127, with command usage logged to stdout.
// Currently cases 1 and 2 are handled.
// TODO KEM: Handle exit 127 case. How does this work on non-Unix platforms?
type ErrCLIUsage struct {
	stderr string
}

func (e *ErrCLIUsage) Error() string {
	return e.stderr
}

// ErrManualEnvVar is returned when an env var that should be set programatically via an option or method
// is set via the manual environment passing functions.
type ErrManualEnvVar struct {
	name string
}

func (err *ErrManualEnvVar) Error() string {
	return fmt.Sprintf("manual setting of env var %q detected", err.name)
}

type ErrMissingVar struct {
	VariableName string
}

func (err *ErrMissingVar) Error() string {
	return fmt.Sprintf("variable %q was required but not supplied", err.VariableName)
}
