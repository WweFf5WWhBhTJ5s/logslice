// logslice is a fast structured log query tool that parses and filters
// JSON logs from files or stdin using a minimal query DSL.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/parser"
)

const usage = `logslice — fast structured log query tool

Usage:
  logslice [options] [file...]

Options:
  -f, --filter <expr>   Filter expression (e.g. 'level=error', 'status>=400').
                        May be repeated to AND multiple filters.
  -o, --output <fmt>    Output format: json (default), pretty, compact.
  -F, --fields <list>   Comma-separated list of fields to include in output.
  --skip-invalid        Skip lines that are not valid JSON (default: true).
  --stop-on-invalid     Stop processing on the first invalid JSON line.
  -h, --help            Show this help message.

Examples:
  logslice -f 'level=error' app.log
  cat app.log | logslice -f 'status>=500' -o pretty
  logslice -f 'level=error' -F 'time,msg,error' app.log
`

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() { fmt.Fprint(stderr, usage) }

	var filterExprs multiFlag
	fs.Var(&filterExprs, "f", "Filter expression (repeatable)")
	fs.Var(&filterExprs, "filter", "Filter expression (repeatable)")

	outputFmt := fs.String("o", "json", "Output format: json, pretty, compact")
	fs.StringVar(outputFmt, "output", "json", "Output format: json, pretty, compact")

	fields := fs.String("F", "", "Comma-separated fields to include")
	fs.StringVar(fields, "fields", "", "Comma-separated fields to include")

	skipInvalid := fs.Bool("skip-invalid", true, "Skip invalid JSON lines")
	stopOnInvalid := fs.Bool("stop-on-invalid", false, "Stop on first invalid JSON line")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Build filter chain from expressions.
	chain, err := filter.NewChain(filterExprs)
	if err != nil {
		return fmt.Errorf("invalid filter: %w", err)
	}

	// Parse field selection.
	var selectedFields []string
	if *fields != "" {
		for _, f := range strings.Split(*fields, ",") {
			if f = strings.TrimSpace(f); f != "" {
				selectedFields = append(selectedFields, f)
			}
		}
	}

	// Build formatter.
	fmt_, err := output.New(*outputFmt, stdout, selectedFields)
	if err != nil {
		return fmt.Errorf("invalid output format: %w", err)
	}

	// Determine input sources.
	readers := fs.Args()

	processReader := func(r io.Reader, name string) error {
		p := parser.New(r)
		for {
			entry, err := p.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				if *stopOnInvalid {
					return fmt.Errorf("%s: %w", name, err)
				}
				if !*skipInvalid {
					fmt.Fprintf(stderr, "warning: %s: %v\n", name, err)
				}
				continue
			}
			if chain.Match(entry) {
				if werr := fmt_.Write(entry); werr != nil {
					return werr
				}
			}
		}
		return nil
	}

	if len(readers) == 0 {
		// Read from stdin when no files are provided.
		return processReader(stdin, "<stdin>")
	}

	for _, path := range readers {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		if err := processReader(f, path); err != nil {
			f.Close()
			return err
		}
		f.Close()
	}
	return nil
}

// multiFlag allows a flag to be specified multiple times.
type multiFlag []string

func (m *multiFlag) String() string { return strings.Join(*m, ", ") }
func (m *multiFlag) Set(v string) error {
	*m = append(*m, v)
	return nil
}
