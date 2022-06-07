# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[unreleased]: https://github.com/ndaba1/gommander/compare/v0.1.3...HEAD
[0.1.3]: https://github.com/ndaba1/gommander/compare/v0.1.2...v0.1.3
[0.1.2]: https://github.com/ndaba1/gommander/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/ndaba1/gommander/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/ndaba1/gommander/releases/tag/v0.1.0
