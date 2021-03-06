## 0.7.0 (2020/11/12)

Enhancements:

* Use mattn-isatty instead of andrew-d/go-termutil #28
* Introduce Go Modules #27
* Use https instead of http in comment. #24
* Add tests. #23

## 0.6.1 (2017/07/14)

Fixes:

* Show compiler version by -v. refs: #19
* feature: introduced --ignore-key. refs #20

## 0.6.0 (2017/04/11)

Fixes:

* sort keys by default. refs: #16
* build: fixed `go get` error for ghr. refs #15

## 0.5.6 (2017/03/21)

Fixes:

* feature: introduced not-equal operators for string.
* expr: re-implemented by go package.

## 0.5.5 (2016/11/24)

Fixes:

* filter: exited when invalid filter expression was given (thanks to @cubicdaiya)

## 0.5.4 (2016/10/31)

Fixes:

* Typo fixed

## 0.5.3 (2016/10/31)

Enhancements:

* Introduce case-insentive comparing operators such as `==*`, `=~*`, `!~*` (thanks to @cubicdaiya)

## 0.5.2 (2016/10/28)

Enhancements:

* Enable perfect string match by '==' (thanks to @cubicdaiya)

Changes:

* Replace codegangsta/cli with urfave/cli (thanks to @b4b4r07)

## 0.5.1 (2016/08/25)

Enhancements:

* Fix typo

## 0.5.0 (2016/08/24)

Enhancements:

* Add new option -expr, -e (thanks to @cubicdaiya)

## 0.4.1 (2015/06/13)

Enhancements:

* filter: fixed panic error when the invalid filter expression is given. (thanks to @cubicdaiya)

## 0.4.0 (2015/10/20)

Enhancements:

* Add `filter` option (thanks to @hirose31)

## 0.3.1 (2014/11/21)

Enhancements:

* Add line feed to error messages

## 0.3.0 (2014/11/21)

Enhancements:

* Read from multiple files

## 0.2.0 (2014/08/13)

Enhancements:

* Read from a file (not only from STDIN)

Changes:

* Print all fields if -k option is not specified
* Print fields in specified keys order

## 0.1.0 (2014/08/13)

* Initial Release

