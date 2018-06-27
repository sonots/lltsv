package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/urfave/cli"
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

	Eliminate labels with "-K" option.

	Example3 $ lltsv -k foo,bar -K file*.log

	Specify input files as arguments.

	Example4 $ lltsv -k resptime,status,uri -f 'resptime > 6' access_log
	         $ lltsv -k resptime,status,uri -f 'resptime > 6' -f 'uri =~ ^/foo' access_log

	Filter output with "-f" option. Available comparing operators are:

    >= > == < <=  (arithmetic (float64))
    == ==* != !=* (string comparison (string))
    =~ !~ =~* !~* (regular expression (string))

        The comparing operators terminated by * behave in case-insensitive.

	You can specify multiple -f options (AND condition).

	Example5 $ lltsv -k resptime,upstream_resptime,diff -f 'diff = resptime - upstream_resptime' access_log
	         $ lltsv -k resptime,upstream_resptime,diff_ms -e 'diff_ms = (resptime - upstream_resptime) * 1000' access_log

	Evaluate value with "-e" option. Available operators are:

	  + - * / (arithmetic (float64))

	Grep period in common log format(common log or ISO 8610) "-t" option. Available operators are:

	Example6 $ lltsv -t 'localtime=01/Jun/2018:00:00:00 +0900~30/Mar/2018:23:59:59 +0900,common' access_log
	         $ lltsv -t 'localtime=2018-01-01T00:00:00+0900~2018-03-30T23:59:59+0900,iso8610' access_log

	Homepage: https://github.com/sonots/lltsv`
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
		cli.StringFlag{
			Name:  "ignore-key, i",
			Usage: "ignored keys to output (multiple keys separated by ,)",
		},
		cli.StringSliceFlag{
			Name:  "filter, f",
			Usage: "filter expression to output",
		},
		cli.StringSliceFlag{
			Name:  "expr, e",
			Usage: "evaluate value by expression to output",
		},
		cli.StringSliceFlag{
			Name:  "timegrep, t",
			Usage: "grep period of time in common log format or iso8601",
		},
	}
	app.Action = doMain

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Fprintf(app.Writer, `%v version %v
Compiler: %s %s
`,
			app.Name,
			app.Version,
			runtime.Compiler,
			runtime.Version())
	}

	app.Run(os.Args)
	os.Exit(exitCode)
}

func doMain(c *cli.Context) error {
	keys := make([]string, 0, 0) // slice with length 0
	if c.String("key") != "" {
		keys = strings.Split(c.String("key"), ",")
	}
	noKey := c.Bool("no-key")
	filters := c.StringSlice("filter")
	exprs := c.StringSlice("expr")
	timegrep := c.StringSlice("timegrep")

	ignoreKeys := make([]string, 0, 0)
	// If -k,--key is specified, -i,--ignore-key is ignored.
	if len(keys) == 0 && c.String("ignore-key") != "" {
		ignoreKeys = strings.Split(c.String("ignore-key"), ",")
	}

	lltsv := newLltsv(keys, ignoreKeys, noKey, filters, exprs, timegrep)

	if len(c.Args()) > 0 {
		for _, filename := range c.Args() {
			file, err := os.Open(filename)
			if err != nil {
				os.Stderr.WriteString("failed to open and read `" + filename + "`.\n")
				exitCode = 1
				return err
			}
			err = lltsv.scanAndWrite(file)
			file.Close()
			if err != nil {
				os.Stderr.WriteString("reading input errored\n")
				exitCode = 1
				return err
			}
		}
	} else {
		file := os.Stdin
		err := lltsv.scanAndWrite(file)
		file.Close()
		if err != nil {
			os.Stderr.WriteString("reading input errored\n")
			exitCode = 1
			return err
		}
	}

	return nil
}
