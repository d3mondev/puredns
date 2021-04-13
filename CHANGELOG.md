# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
- No changes yet.

## [1.0.3] - 2021-04-12
### Fixed
- Remove Cloudflare DNS from the list of trusted resolvers. [Here's why](https://twitter.com/d3mondev/status/1381678504450924552?s=20).
- Increase the default rate limit per trusted resolver to 50.
- Adjust massdns command line parameter -s to limit the size of the initial burst of queries sent to the trusted resolvers.

## [1.0.2] - 2021-03-22
### Fixed
- Fix a badly handled exception during wildcard detection that was halting the process.

## [1.0.1] - 2020-10-15
### Fixed
- Fix a bug where valid subdomains were not saved to a file. [#1](https://github.com/d3mondev/puredns/issues/1)

## [1.0.0] - 2020-08-02
### Added
- Initial implementation.

[Unreleased]: https://github.com/d3mondev/puredns/compare/v1.0.3...HEAD
[1.0.0]: https://github.com/d3mondev/puredns/releases/tag/v1.0.0
[1.0.1]: https://github.com/d3mondev/puredns/releases/tag/v1.0.1
[1.0.2]: https://github.com/d3mondev/puredns/releases/tag/v1.0.2
[1.0.3]: https://github.com/d3mondev/puredns/releases/tag/v1.0.3
