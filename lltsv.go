package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/andrew-d/go-termutil"
	"github.com/mgutz/ansi"
)

type tFuncAppend func([]string, string, string) []string

type Lltsv struct {
	keys       []string
	no_key     bool
	funcAppend tFuncAppend
}

func newLltsv(keys []string, no_key bool) *Lltsv {
	return &Lltsv{
		keys:       keys,
		no_key:     no_key,
		funcAppend: getFuncAppend(no_key),
	}
}

func (lltsv *Lltsv) scanAndWrite(file *os.File, filters map[string]Filter) error {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lvs := lltsv.parseLtsv(line)

		should_output := true

		for key, filter := range filters {
			if !filter(lvs[key]) {
				should_output = false
				break
			}
		}

		if should_output {
			ltsv := lltsv.restructLtsv(lvs)
			os.Stdout.WriteString(ltsv + "\n")
		}
	}
	return scanner.Err()
}

// lvs: label and value pairs
func (lltsv *Lltsv) restructLtsv(lvs map[string]string) string {
	// specified keys or all keys
	orders := lltsv.keys
	if len(lltsv.keys) == 0 {
		orders = keysInMap(lvs)
	}
	// make slice with enough capacity so that append does not newly create object
	// cf. http://golang.org/pkg/builtin/#append
	selected := make([]string, 0, len(orders))
	for _, label := range orders {
		value := lvs[label]
		selected = lltsv.funcAppend(selected, label, value)
	}
	return strings.Join(selected, "\t")
}

func (lltsv *Lltsv) parseLtsv(line string) map[string]string {
	columns := strings.Split(line, "\t")
	lvs := make(map[string]string)
	for _, column := range columns {
		l_v := strings.SplitN(column, ":", 2)
		if len(l_v) < 2 {
			continue
		}
		label, value := l_v[0], l_v[1]
		lvs[label] = value
	}
	return lvs
}

// Return function pointer to avoid `if` evaluation occurs in each iteration
func getFuncAppend(no_key bool) tFuncAppend {
	if no_key {
		return func(selected []string, label string, value string) []string {
			return append(selected, value)
		}
	} else {
		if termutil.Isatty(os.Stdout.Fd()) {
			return func(selected []string, label string, value string) []string {
				return append(selected, ansi.Color(label, "green")+":"+ansi.Color(value, "magenta"))
			}
		} else {
			// if pipe or redirect
			return func(selected []string, label string, value string) []string {
				return append(selected, label+":"+value)
			}
		}
	}
}

func keysInMap(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}
