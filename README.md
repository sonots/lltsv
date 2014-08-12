# lltsv

List specified keys of LTSV (Labeled Tab Separated Values)

# Description

`lltsv` is a command line golang tool to list specified keys of LTSV (Labeled Tab Separated Values) format.

Example: 

```bash
$ echo "foo:aaa\tbar:bbb" | lltsv -k foo,bar
foo:aaa   bar:bbb
```

The output is colorized as default when you outputs to a terminal. 
The coloring is disabled if you pipe or redirect outputs.

Example2:

```bash
$	"foo:aaa\tbar:bbb" | lltsv -k foo,bar -K
aaa       bbb
```

You may eliminate labels with `-K` option


## Usage

```
$ lltsv -h
NAME:
   lltsv - List specified keys of LTSV (Labeled Tab Separated Values)

USAGE:
   lltsv [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --key, -k            keys to output (multiple keys separated by ,)
   --no-key, -K         output without keys (and without color)
   --version, -v        print the version
   --help, -h           show help
```

## Installation

To install, use go get and make install.

```bash
$ go get -d github.com/sonots/lltsv
$ cd $GOPATH/src/github.com/sonots/lltsv
$ make install 
```

## Contribution

1. Fork (https://github.com/sonots/lltsv/fork)
2. Create a feature branch
3. Commit your changes
4. Rebase your local changes against the master branch
5. Run test suite with the go test ./... command and confirm that it passes
6. Run gofmt -s
7. Create new Pull Request

## Copyright

See [LICENSE](./LICENSE)

## Special Thanks

This is a golang fork of perl version created by [id:key_amb](http://keyamb.hatenablog.com/). Thanks!
