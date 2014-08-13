package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/andrew-d/go-termutil"
	"github.com/codegangsta/cli"
	"github.com/mgutz/ansi"
)

// os.Exit forcely kills process, so let me share this global variable to terminate at the last
var exitCode = 0

func main() {
	app := cli.NewApp()
	app.Name = "lltsv"
	app.Version = Version
	app.Usage = `List specified keys of LTSV (Labeled Tab Separated Values)

	Example1 $ echo "foo:aaa\tbar:bbb" | lltsv -k foo,bar
	foo:aaa   bar:bbb

	The output is colorized as default when you outputs to a terminal. 
	The coloring is disabled if you pipe or redirect outputs.

	Example2 $ echo "foo:aaa\tbar:bbb" | lltsv -k foo,bar -K
	aaa       bbb

	You may eliminate labels with "-K" option`
	app.Author = "sonots"
	app.Email = "sonots@gmail.com"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "key, k",
			Usage: "keys to output (multiple keys separated by ,)",
		},
		cli.BoolFlag{
			Name:  "no-key, K",
			Usage: "output without keys (and without color)",
		},
	}
	app.Action = doMain
	app.Run(os.Args)
	os.Exit(exitCode)
}

func doMain(c *cli.Context) {
	keys := make([]string, 0, 0) // slice with length 0
	if c.String("key") != "" {
		keys = strings.Split(c.String("key"), ",")
	}
	no_key := c.Bool("no-key")
	funcAppend := getFuncAppend(no_key)

	var file *os.File
	if len(c.Args()) > 0 {
		filename := c.Args()[0]
		var err error
		file, err = os.Open(filename)
		if err != nil {
			os.Stderr.WriteString("failed to open and read " + filename)
			exitCode = 1
			return
		}
	} else {
		file = os.Stdin
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lvs := parseLtsv(line)
		ltsv := restructLtsv(keys, lvs, funcAppend)
		os.Stdout.WriteString(ltsv + "\n")
	}
	if err := scanner.Err(); err != nil {
		os.Stderr.WriteString("reading standard input errored")
		exitCode = 1
		return
	}
}

func restructLtsv(keys []string, lvs map[string]string, funcAppend func([]string, string, string) []string) string {
	// specified keys or all keys
	orders := keys
	if len(keys) == 0 {
		orders = keysInMap(lvs)
	}
	// make slice with enough capacity so that append does not newly create object
	// cf. http://golang.org/pkg/builtin/#append
	selected := make([]string, 0, len(orders))
	for _, label := range orders {
		value := lvs[label]
		selected = funcAppend(selected, label, value)
	}
	return strings.Join(selected, "\t")
}

func parseLtsv(line string) map[string]string {
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
func getFuncAppend(no_key bool) func([]string, string, string) []string {
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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
