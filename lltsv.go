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

	Example2 $ "foo:aaa\tbar:bbb" | lltsv -k foo,bar -K
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
	if len(c.String("key")) == 0 {
		os.Stderr.WriteString("-k <key> option is required\n")
		exitCode = 1
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	keys := strings.Split(c.String("key"), ",")
	no_key := c.Bool("no-key")
	funcAppend := getFuncAppend(no_key)

	for scanner.Scan() {
		text := scanner.Text()
		lvs := strings.Split(text, "\t")
		// make slice with enough capacity so that append does not newly create object
		// cf. http://golang.org/pkg/builtin/#append
		selected := make([]string, 0, len(keys))
		for _, lv := range lvs {
			l_v := strings.SplitN(lv, ":", 2)
			if len(l_v) < 2 {
				continue
			}
			label, value := l_v[0], l_v[1]
			if stringInSlice(label, keys) {
				selected = funcAppend(selected, label, value, lv)
			}
		}
		os.Stdout.WriteString(strings.Join(selected, "\t") + "\n")
	}
	if err := scanner.Err(); err != nil {
		os.Stderr.WriteString("reading standard input errored")
		exitCode = 1
		return
	}
}

// Return function pointer to avoid `if` evaluation occurs in each iteration
func getFuncAppend(no_key bool) func([]string, string, string, string) []string {
	if no_key {
		return func(selected []string, label string, value string, lv string) []string {
			return append(selected, value)
		}
	} else {
		if termutil.Isatty(os.Stdout.Fd()) {
			return func(selected []string, label string, value string, lv string) []string {
				return append(selected, ansi.Color(label, "green")+":"+ansi.Color(value, "magenta"))
			}
		} else {
			return func(selected []string, label string, value string, lv string) []string {
				return append(selected, lv)
			}
		}
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
