package main

import (
	"os"
	"strings"

	"github.com/codegangsta/cli"
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

	lltsv := newLltsv(keys, no_key)

	if len(c.Args()) > 0 {
		for _, filename := range c.Args() {
			file, err := os.Open(filename)
			if err != nil {
				os.Stderr.WriteString("failed to open and read `" + filename + "`.\n")
				exitCode = 1
				return
			}
			lltsv.scanAndWrite(file)
			file.Close()
		}
	} else {
		file := os.Stdin
		lltsv.scanAndWrite(file)
		file.Close()
	}
}
