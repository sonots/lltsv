package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/andrew-d/go-termutil"
	"github.com/mgutz/ansi"
)

type tFuncAppend func([]string, string, string) []string
type tFuncFilter func(string) bool
type tFuncTimeGreper func(string) bool

// Lltsv is a context for processing LTSV.
type Lltsv struct {
	keys            []string
	ignoreKeyMap    map[string]struct{}
	noKey           bool
	filters         []string
	exprs           []string
	funcAppend      tFuncAppend
	funcFilters     map[string]tFuncFilter
	exprRunners     map[string]*ExprRunner
	funcTimeGrepers map[string]tFuncTimeGreper
}

func newLltsv(keys []string, ignoreKeys []string, noKey bool, filters []string, exprs []string, timegreps []string) *Lltsv {
	ignoreKeyMap := make(map[string]struct{})
	for _, key := range ignoreKeys {
		ignoreKeyMap[key] = struct{}{}
	}
	return &Lltsv{
		keys:            keys,
		ignoreKeyMap:    ignoreKeyMap,
		noKey:           noKey,
		filters:         filters,
		exprs:           exprs,
		funcAppend:      getFuncAppend(noKey),
		funcFilters:     getFuncFilters(filters),
		exprRunners:     getExprRunners(exprs),
		funcTimeGrepers: getTimeGrepers(timegreps),
	}
}

func (lltsv *Lltsv) scanAndWrite(file *os.File) error {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lvs, keys := lltsv.parseLtsv(line)
		lltsv.expr(lvs)

		if lltsv.filter(lvs) && lltsv.timegrep(lvs) {
			ltsv := lltsv.restructLtsv(lvs, keys)
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

func (lltsv *Lltsv) timegrep(lvs map[string]string) bool {
	shouldOutput := true

	for key, funcTimeGreper := range lltsv.funcTimeGrepers {
		if !funcTimeGreper(lvs[key]) {
			shouldOutput = false
			break
		}
	}
	return shouldOutput
}

// lvs: label and value pairs
func (lltsv *Lltsv) restructLtsv(lvs map[string]string, keys []string) string {
	// specified keys or all keys
	orders := lltsv.keys
	if len(lltsv.keys) == 0 {
		orders = keys
	}
	// make slice with enough capacity so that append does not newly create object
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
	keys := make([]string, 0, len(columns))
	for _, column := range columns {
		lv := strings.SplitN(column, ":", 2)
		if len(lv) < 2 {
			continue
		}
		label, value := lv[0], lv[1]
		lvs[label] = value
		keys = append(keys, label)
	}
	return lvs, keys
}

// Return function pointer to avoid `if` evaluation occurs in each iteration
func getFuncAppend(noKey bool) tFuncAppend {
	if noKey {
		return func(selected []string, label string, value string) []string {
			return append(selected, value)
		}
	}

	if termutil.Isatty(os.Stdout.Fd()) {
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

func getTimeGrepers(greps []string) map[string]tFuncTimeGreper {
	funcTimeGreps := map[string]tFuncTimeGreper{}
	for _, f := range greps {
		token := strings.SplitN(f, "=", 2)
		if len(token) != 2 {
			log.Printf("expression is invalid: %s\n", f)
			continue
		}
		key := token[0]
		token = strings.SplitN(token[1], "~", 2)
		if len(token) != 2 {
			log.Printf("expression is invalid: %s\n", f)
			continue
		}
		from := token[0]
		token = strings.SplitN(token[1], ",", 2)
		if len(token) != 2 {
			log.Printf("expression is invalid: %s\n", f)
			continue
		}
		to := token[0]
		formatType := token[1]
		format := ""
		if formatType == "iso8601" {
			format = "2006-01-02T15:04:05-0700"
		} else if formatType == "common" {
			format = "02/Jan/2006:15:04:05 -0700"
		} else {
			log.Printf("expression is invalid: %s\n", f)
			continue
		}
		start, err := time.Parse(format, from)
		end, err := time.Parse(format, to)
		if err != nil {
			log.Printf("expression is invalid: %s\n", f)
			continue
		}
		funcTimeGreps[key] = func(val string) bool {
			v, err := time.Parse(format, val)
			if err != nil {
				log.Println(err)
				return false
			}
			return start.Unix() <= v.Unix() && v.Unix() <= end.Unix()
		}
	}
	return funcTimeGreps
}
