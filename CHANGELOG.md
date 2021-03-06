# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.1] - 2022-07-16

### Added

- Updated the Readme documentation to include the newly added functionality for arg types and methods.
- Added new predefined argument types: `file` for filenames and directories and `uint` for unsigned integers
- Enabled validation of argument values using regular expressions via the `Arg.ValidatorRegex()` method
- Integrated the `AllowNegativeNumbers` setting into the program
- Enabled customization of the help headings and help values to be used for sub-cmds, arguments, flags and options.
- Implemented a simple caching policy to improve performance
- Integrated `DisableColor` setting and functionality
- Created a simple contributing guide

### Changed

- Improved general package performance and structs creation latency
- Updated existing and added more benchmarks
- Modified some of the error messages to be more concise
- Renamed struct fields for Options, Arguments and Flags to make them public and enable creation using constructors
- Modified the `FormatGenerator` interface to take in a pointer to the app for settings integration
- Moved the suggest subcommands functionality from utils to a method of the `Command` struct

### Removed

- Removed unused assertion utility functions such as the assertArrEq and assertStructEq.
- Removed the `.compare()` method for the StructComp interface implementations. Now comparing structs using the reflect.DeepEqual method.

## [0.2.0] - 2022-06-30

### Added

- Added support for multiple validator functions for a single argument
- Added loose support for argument types i.e `<int:arg1>` will ensure only integer values are accepted. Others "types" include `float` and `bool` with `str` being the default.
- Added more tests for package components, significantly improved the coverage
- Created functions for indent and dedent functionality in the utils file
- Introduced a new and better word wrap function
- Command disuccion printing now functions correctly
- Created custom assertion utilities for easier and better testing
- Created a shell script for use in comparing benchmarks along with a make target for running benchstat
- Configured CI test workflow to create an automated issue on w/f failure

### Changed

- Moved helpwriter functionality into a separate file
- Removed the doc.go file and moved package docstring to gommander.go file
- Renamed the app.Info() method to app.Discussion()
- Refactored the tests to use the newly created assertion utilities
- Enabled some more linters in the golangci-lint config
- Modified program to check for test mode before exiting when performing runtime checks
- Updated issue template to include assignee and added feature_request template

### Fixed

- Corrected the package inner parse method to use the args passed as input, rather than reading directly from os.Args
- Patched issue #31 causing arg parsing for the root cmd to fail

## [0.1.7] - 2022-06-25

### Added

- Integrated codecov into the CI workflow and added a coverage badge in the README
- Created an issue template for the project and added the CODEOWNERS file

### Changed

- Moved all the parser error handling functionality into a single method rather than spread out all over the parser code.
- Edited the fields of the Error struct to make errors more simplified
- Refactored the suggest cmds functionality to make it more "suggestive"

### Fixed

- Patched the parser error causing the uncaught unknown-command error
- Fixed the parser issue #21 causing parsing of multiple positional args to fail

## [0.1.6] - 2022-06-24

### Fixed

- Patched issue #15 causing runtime errors when optional arguments are included in a subcommand

## [0.1.5] - 2022-06-22

### Fixed

- Changed the variable naming convention from snake-case to camel-case to resolve linter warnings
- Resolved all other golangci-lint warnings and errors
- Removed deprecated linters from config

### Changed

- Refactored the previously named `Error` designation into `ErrorMsg` to avoid naming conflicts
- Renamed the program error from `gommander.GommanderError` into `gommander.Error`

## [0.1.4] - 2022-06-12

### Fixed

- Patched the parser error causing parsing of empty argument values
- Updated argument tests to include the format with the arg-default value included

### Added

- Created formatter utility functions for: adding & printing directly, printing by color rather than by designation, coloring and printing directly
- Set argument default values to be printed out if any are present
- Added new error-handling section in the README along with screenshots corresponding to the errors
- Documented the public formatter interface in the README in its own section

### Changed

- Removed the unused new variadic theme function
- Removed all the resolved todos

## [0.1.3] - 2022-06-07

### Fixed

- Patched critical bug causing runtime error when no args are supplied to the program

## [0.1.2] - 2022-06-07

### Added

- Added a new argument method for adding custom validator functions `.ValidatorFunc()`
- Added more method docstrings and new examples in docs
- Included screenshots for the README defined examples

### Changed

- Made the `fmter.Add()` and `fmter.Print()` methods public
- Commands can now have multiple aliases instead of just one

### Fixed

- The `ShowCommandAliases` setting now functions correctly
- Resolved the uncaught parser errors issue for subcommands
- Patched the parser complications when parsing arguments for nested subcommands
- Fixed the indentation issues in the README code examples

## [0.1.1] - 2022-06-02

### Fixed

- Patched parser bug causing unnecessary help printing
- Fixed simple typos and formatting issue in README

## [0.1.0] - 2022-06-01

### Added

- Basic package functionality
- Simple examples and indepth docs
- Fair amount of tests and coverage

[unreleased]: https://github.com/ndaba1/gommander/compare/v0.2.1...HEAD
[0.2.1]: https://github.com/ndaba1/gommander/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/ndaba1/gommander/compare/v0.1.7...v0.2.0
[0.1.7]: https://github.com/ndaba1/gommander/compare/v0.1.6...v0.1.7
[0.1.6]: https://github.com/ndaba1/gommander/compare/v0.1.5...v0.1.6
[0.1.5]: https://github.com/ndaba1/gommander/compare/v0.1.3...v0.1.5
[0.1.4]: https://github.com/ndaba1/gommander/compare/v0.1.3...v0.1.4
[0.1.3]: https://github.com/ndaba1/gommander/compare/v0.1.2...v0.1.3
[0.1.2]: https://github.com/ndaba1/gommander/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/ndaba1/gommander/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/ndaba1/gommander/releases/tag/v0.1.0
