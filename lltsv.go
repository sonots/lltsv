package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/andrew-d/go-termutil"
	"github.com/mgutz/ansi"
)

type tFuncAppend func([]string, string, string) []string
type tFuncFilter func(string) bool

type Lltsv struct {
	keys        []string
	no_key      bool
	filters     []string
	funcAppend  tFuncAppend
	funcFilters map[string]tFuncFilter
}

func newLltsv(keys []string, no_key bool, filters []string) *Lltsv {
	return &Lltsv{
		keys:        keys,
		no_key:      no_key,
		filters:     filters,
		funcAppend:  getFuncAppend(no_key),
		funcFilters: getFuncFilters(filters),
	}
}

func (lltsv *Lltsv) scanAndWrite(file *os.File) error {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lvs := lltsv.parseLtsv(line)

		if lltsv.filter(lvs) {
			ltsv := lltsv.restructLtsv(lvs)
			os.Stdout.WriteString(ltsv + "\n")
		}
	}
	return scanner.Err()
}

func (lltsv *Lltsv) filter(lvs map[string]string) bool {
	should_output := true

	for key, funcFilter := range lltsv.funcFilters {
		if !funcFilter(lvs[key]) {
			should_output = false
			break
		}
	}

	return should_output
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

func getFuncFilters(filters []string) map[string]tFuncFilter {
	funcFilters := map[string]tFuncFilter{}
	for _, f := range filters {
		token := strings.SplitN(f, " ", 3)
		key := token[0]
		switch token[1] {
		case ">", ">=", "==", "<=", "<":
			r, err := strconv.ParseFloat(token[2], 64)
			if err != nil {
				log.Fatal(err)
			}

			funcFilters[key] = func(val string) bool {
				num, err := strconv.ParseFloat(val, 64)
				if err != nil {
					log.Println(err)
					return false
				}
				switch token[1] {
				case ">":
					return num > r
				case ">=":
					return num >= r
				case "==":
					return num == r
				case "<=":
					return num <= r
				case "<":
					return num < r
				default:
					return false
				}
			}
		case "=~", "!~":
			re := regexp.MustCompile(token[2])
			funcFilters[key] = func(val string) bool {
				switch token[1] {
				case "=~":
					return re.MatchString(val)
				case "!~":
					return !re.MatchString(val)
				default:
					return false
				}
			}
		}
	}
	return funcFilters
}

func keysInMap(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}
