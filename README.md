# lltsv

List specified keys of LTSV (Labeled Tab Separated Values)

# Description

`lltsv` is a command line tool written in golang to list specified keys of LTSV (Labeled Tab Separated Values) text.

Example 1: 

```bash
$ echo "foo:aaa\tbar:bbb\tbaz:ccc" | lltsv -k foo,bar
foo:aaa   bar:bbb
```

The output is colorized as default when you outputs to a terminal. 
The coloring is disabled if you pipe or redirect outputs.

Example 2:

```bash
$ echo "foo:aaa\tbar:bbb\tbaz:ccc" | lltsv -k foo,bar -K
aaa       bbb
```

You may eliminate labels with `-K` option. 

**How Useful?**

LTSV format is not `awk` friendly (I think), but `lltsv` can help it: 

```bash
$ echo -e "time:2014-08-13T14:10:10Z\tstatus:200\ntime:2014-08-13T14:10:12Z\tstatus:500" \
  | lltsv -k time,status -K | awk '$2 == 500'
2014-08-13T14:10:12Z    500
```

Useful!

## Installation

Executable binaries are available at [releases](https://github.com/sonots/lltsv/releases).

For example, for linux x86_64, 

```bash
$ wget https://github.com/sonots/lltsv/releases/download/v0.1.0/lltsv_linux_amd64 -O lltsv
$ chmod a+x lltsv
```

If you have the go runtime installed, you may use go get. 

```bash
$ go get github.com/sonots/lltsv
```

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

## ToDo

1. write tests

## Build

To build, use go get and make

```
$ go get -d github.com/sonots/lltsv
$ cd $GOPATH/src/github.com/sonots/lltsv
$ make
```

To release binaries, I use [gox](https://github.com/mitchellh/gox) and [ghr](https://github.com/tcnksm/ghr)

```
go get github.com/mitchellh/gox
gox -build-toolchain # only first time
go get github.com/tcnksm/ghr

cd $GOPATH/src/github.com/sonots/lltsv/pkg
gox ../...
ghr <tag> .
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

MEMO: golang version was 5x faster than perl version
