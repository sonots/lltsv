package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/mgutz/ansi"
)

type tFuncAppend func([]string, string, string) []string
type tFuncFilter func(string) bool

// Lltsv is a context for processing LTSV.
type Lltsv struct {
	keys         []string
	ignoreKeyMap map[string]struct{}
	noKey        bool
	filters      []string
	exprs        []string
	funcAppend   tFuncAppend
	funcFilters  map[string]tFuncFilter
	exprRunners  map[string]*ExprRunner
}

func newLltsv(keys []string, ignoreKeys []string, noKey bool, filters []string, exprs []string) *Lltsv {
	ignoreKeyMap := make(map[string]struct{})
	for _, key := range ignoreKeys {
		ignoreKeyMap[key] = struct{}{}
	}
	return &Lltsv{
		keys:         keys,
		ignoreKeyMap: ignoreKeyMap,
		noKey:        noKey,
		filters:      filters,
		exprs:        exprs,
		funcAppend:   getFuncAppend(noKey),
		funcFilters:  getFuncFilters(filters),
		exprRunners:  getExprRunners(exprs),
	}
}

func (lltsv *Lltsv) scanAndWrite(file *os.File) error {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lvs, labels := lltsv.parseLtsv(line)
		lltsv.expr(lvs)

		if lltsv.filter(lvs) {
			ltsv := lltsv.restructLtsv(lvs, labels)
			os.Stdout.WriteString(ltsv + "\n")
		}

	}
	return scanner.Err()
}

func (lltsv *Lltsv) filter(lvs map[string]string) bool {
	shouldOutput := true

	for key, funcFilter := range lltsv.funcFilters {
		if !funcFilter(lvs[key]) {
			shouldOutput = false
			break
		}
	}

	return shouldOutput
}

func (lltsv *Lltsv) expr(lvs map[string]string) {
	ctx := &ExprContext{
		vars: lvs,
	}
	for key, runner := range lltsv.exprRunners {
		v, err := evalExpr(runner.expr, ctx)
		if err == nil {
			lvs[key] = v.String()
		}
	}
}

// lvs: label and value pairs
func (lltsv *Lltsv) restructLtsv(lvs map[string]string, labels []string) string {
	// Sort in the order of -k, or follow the order of the input file.
	orders := lltsv.keys
	if len(lltsv.keys) == 0 {
		orders = labels
	}
	// Make slice with enough capacity so that append does not newly create object
	// cf. https://golang.org/pkg/builtin/#append
	selected := make([]string, 0, len(orders))
	for _, label := range orders {
		if _, ok := lltsv.ignoreKeyMap[label]; ok {
			continue
		}
		value := lvs[label]
		selected = lltsv.funcAppend(selected, label, value)
	}
	return strings.Join(selected, "\t")
}

func (lltsv *Lltsv) parseLtsv(line string) (map[string]string, []string) {
	columns := strings.Split(line, "\t")
	lvs := make(map[string]string)
	labels := make([]string, 0, len(columns))
	for _, column := range columns {
		lv := strings.SplitN(column, ":", 2)
		if len(lv) < 2 {
			continue
		}
		label, value := lv[0], lv[1]
		lvs[label] = value
		labels = append(labels, label)
	}
	return lvs, labels
}

// Return function pointer to avoid `if` evaluation occurs in each iteration
func getFuncAppend(noKey bool) tFuncAppend {
	if noKey {
		return func(selected []string, label string, value string) []string {
			return append(selected, value)
		}
	}

	if isatty.IsTerminal(os.Stdout.Fd()) {
		return func(selected []string, label string, value string) []string {
			return append(selected, ansi.Color(label, "green")+":"+ansi.Color(value, "magenta"))
		}
	}

	// if pipe or redirect
	return func(selected []string, label string, value string) []string {
		return append(selected, label+":"+value)
	}
}

func getFuncFilters(filters []string) map[string]tFuncFilter {
	funcFilters := map[string]tFuncFilter{}
	for _, f := range filters {
		token := strings.SplitN(f, " ", 3)
		if len(token) < 3 {
			log.Fatalf("filter expression is invalid: %s\n", f)
		}
		key := token[0]
		switch token[1] {
		case ">", ">=", "<=", "<":
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
				case "<=":
					return num <= r
				case "<":
					return num < r
				default:
					return false
				}
			}
		case "==":
			funcFilters[key] = func(val string) bool {
				return val == token[2]
			}
		case "==*":
			funcFilters[key] = func(val string) bool {
				return strings.ToLower(val) == strings.ToLower(token[2])
			}
		case "!=":
			funcFilters[key] = func(val string) bool {
				return val != token[2]
			}
		case "!=*":
			funcFilters[key] = func(val string) bool {
				return strings.ToLower(val) != strings.ToLower(token[2])
			}
		case "=~", "!~", "=~*", "!~*":
			if token[1] == "=~*" || token[1] == "!~*" {
				token[2] = strings.ToLower(token[2])
			}
			re := regexp.MustCompile(token[2])
			funcFilters[key] = func(val string) bool {
				switch token[1] {
				case "=~":
					return re.MatchString(val)
				case "!~":
					return !re.MatchString(val)
				case "=~*":
					return re.MatchString(strings.ToLower(val))
				case "!~*":
					return !re.MatchString(strings.ToLower(val))
				default:
					return false
				}
			}
		}
	}
	return funcFilters
}

func getExprRunners(exprs []string) map[string]*ExprRunner {
	funcExprs := make(map[string]*ExprRunner, len(exprs))
	for _, f := range exprs {
		token := strings.SplitN(f, "=", 2)
		if len(token) != 2 {
			log.Printf("expression is invalid: %s\n", f)
			continue
		}

		expr, err := parseExpr(token[1])
		if err != nil {
			log.Printf("expression is invalid: %s\n", f)
			continue
		}

		key := strings.Trim(token[0], " ")

		funcExprs[key] = &ExprRunner{
			expr: expr,
		}
	}
	return funcExprs
}
