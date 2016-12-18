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

Eliminate labels with `-K` option.

Example 3:

```bash
$ lltsv -k foo,bar -K file*.log
```

Specify input files as arguments.

Example4:

```bash
$ lltsv -k resptime,status,uri -f 'resptime > 6' access_log
$ lltsv -k resptime,status,uri -f 'resptime > 6' -f 'uri =~ ^/foo' access_log
```

Filter output with "-f" option. Available comparing operators are:

```
  >= > == < <=  (arithmetic (float64))
  == ==*        (string comparison (string))
  =~ !~ =~* !~* (regular expression (string))
```

The comparing operators terminated by __*__ behave in case-insensitive.

You can specify multiple -f options (AND condition).

Example5:

```bash
$ lltsv -k resptime,upstream_resptime,diff -e 'diff = resptime - upstream_resptime' access_log
$ lltsv -k resptime,upstream_resptime,diff_ms -e 'diff_ms = (resptime - upstream_resptime) * 1000' access_log
```

Evaluate value with "-e" option. Available operators are:

```
  + - * / (arithmetic (float64))
```

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
$ wget https://github.com/sonots/lltsv/releases/download/v0.3.0/lltsv_linux_amd64 -O lltsv
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
   0.5.1

AUTHOR(S):
   sonots <sonots@gmail.com>

COMMANDS:
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --key, -k                                            keys to output (multiple keys separated by ,)
   --no-key, -K                                         output without keys (and without color)
   --filter, -f [--filter option --filter option]       filter expression to output
   --expr, -e [--expr option --expr option]             evaluate value by expression to output
   --help, -h                                           show help
   --version, -v                                        print the version
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

mkdir -p pkg && cd pkg && gox ../...
ghr vX.X.X .
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
