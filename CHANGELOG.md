# CHANGELOG

## 0.3.1

* Fix bug when invoking `GetAll()` with non-readable sysctls
* Upgrade all vendored dependencies
* Run nilness check in CI

## 0.3.0

* Correctly evaluate which sysctl files are readable when running `GetPattern()`
* Upgrade all vendored dependencies
* Fix builds for Go 1.18
* Code linting
* Test improvements

## 0.2.0

* Add `Client` type
* Upgrade all dependencies
* Improve test coverage

## 0.1.1

* Fix configuration file parsing bug
* Go modules support

## 0.1.0

* Support for getting and setting sysctls
* Support for loading sysctl from configuration files
