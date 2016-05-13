// package main is the breadth of the terraform-remote-env code.
//
// terraform-remote-env is a small little helper for Terraform that
// connects to a remote state, and and exports all root module outputs
// as TF_VAR_ environment variables, to allow for use within Terraform
// count or provider interpolations, or with other tools.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform/state/remote"
)

// programConfig defines the configuration for the program. This will be
// mainly used in arg parsing.
type programConfig struct {
	// The backend to use.
	backend string

	// Backend key/value options.
	backendConfig map[string]string

	// The prefix to add to output.
	prefix string
}

// flagKV is a flag.Value implementation for simple parsing of key=value
// pairs in options.
type flagKV map[string]string

// String presents the default value for the flagKV implementation.
func (v *flagKV) String() string {
	return ""
}

// Set is the option setter for the flagKV implementation.
func (v *flagKV) Set(raw string) error {
	idx := strings.Index(raw, "=")
	if idx == -1 {
		return fmt.Errorf("Error parsing string %s: no \"=\" character found", raw)
	}
	if *v == nil {
		*v = make(map[string]string)
	}

	key, value := raw[0:idx], raw[idx+1:]
	(*v)[key] = value
	return nil
}

// parseArgs parses the command-line arguments given on the command line.
func parseArgs() programConfig {
	var cfg programConfig
	cmdFlags := flag.NewFlagSet("args", flag.ContinueOnError)
	cmdFlags.StringVar(&cfg.backend, "backend", "atlas", "The remote config backend to use")
	cmdFlags.StringVar(&cfg.prefix, "prefix", "", "The prefix to add to output variables")
	cmdFlags.Var((*flagKV)(&cfg.backendConfig), "backend-config", "Backend config parameters, in k=v format")
	cmdFlags.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s OPTIONS\n", os.Args[0])
		cmdFlags.PrintDefaults()
	}
	if err := cmdFlags.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing CLI flags: %s\n\n", err)
		os.Exit(1)
	}
	return cfg
}

// getState loads the Terraform state, provided a certain config.
func getState(cfg programConfig) (map[string]interface{}, error) {
	client, err := remote.NewClient(cfg.backend, cfg.backendConfig)
	if err != nil {
		return nil, err
	}

	state := &remote.State{Client: client}
	if err := state.RefreshState(); err != nil {
		return nil, err
	}

	var outputs map[string]interface{}
	if !state.State().Empty() {
		outputs = state.State().RootModule().Outputs
	}
	return outputs, nil
}

// outputState outputs the state.
//
// The output is a single line of variables, like so:
//
// TF_VAR_foo=bar TF_VAR_baz=qux
//
// Any prefix defined by "-prefix=PREFIX" is also added on.
func outputState(cfg programConfig, outputs map[string]interface{}) string {
	s := []string{}

	for k, v := range outputs {
		if cfg.prefix != "" {
			k = fmt.Sprintf("TF_VAR_%s_%s", cfg.prefix, k)
		} else {
			k = fmt.Sprintf("TF_VAR_%s", k)
		}
		s = append(s, fmt.Sprintf("%s=%s", k, v))
	}

	return strings.Join(s, " ")
}

func main() {
	cfg := parseArgs()
	outputs, err := getState(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting remote state: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", outputState(cfg, outputs))
}
