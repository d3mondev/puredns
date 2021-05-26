# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
- No changes yet.

## [2.0.1] - 2021-06-25
### Fixed
- Wildcard subdomains with only CNAME records were not being filtered properly.

## [2.0.0] - 2021-05-03
### Added
- Stdin can be used in place of the domain list or wordlist files. See help for examples.
- Quiet flag (-q, --quiet) to silence output. Only valid domains are output to stdout when quiet mode is on. [#4](https://github.com/d3mondev/puredns/issues/4)
- Attempt to detect DNS load balancing during wildcard detection. Use flag -n, --wildcard-tests to specify the number of DNS queries to perform to detect all the possible IPs for a subdomain.
- Add ability to specify a maximum batch size of domains to process at once during wildcard detection with --wildcard-batch. This is to help prevent memory issues that can happen on very large lists (70M+ wildcard subdomains).
- Progress bar during wildcard detection.
- Selected options are displayed at the start of the program.
- Add sponsors command to view active [Github sponsors](https://github.com/sponsors/d3mondev).

### Changed
- Complete rewrite in Go for more stability and to prepare new features.
- Some command line flags have changed to be POSIX compliant, use --help on commands to see the changes.
- Rewrite wildcard detection algorithm to be more robust.
- Remove dependency on 'pv' and do progress bar and rate limiting internally instead.
- Massdns output file is now written in -o Snl format.
- A default list of public resolvers is no longer provided as a reference. Best results will be obtained by curating your own list, for example using [public-dns.info](https://public-dns.info/nameservers-all.txt) and [DNS Validator](https://github.com/vortexau/dnsvalidator).
- Remove --write-answers command line option since the full wildcard answers are no longer kept in memory to optimize for large files. This might come back in a future release if requested.

### Fixed
- Massdns and wildcard detection will retry on SERVFAIL errors.
- Add missing entries in the massdns cache that resulted in a higher number of DNS queries being made during wildcard detection.
- Fix many edge cases happening around wildcard detection.

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

[Unreleased]: https://github.com/d3mondev/puredns/compare/v2.0.0...HEAD
[1.0.0]: https://github.com/d3mondev/puredns/releases/tag/v1.0.0
[1.0.1]: https://github.com/d3mondev/puredns/releases/tag/v1.0.1
[1.0.2]: https://github.com/d3mondev/puredns/releases/tag/v1.0.2
[1.0.3]: https://github.com/d3mondev/puredns/releases/tag/v1.0.3
[2.0.0]: https://github.com/d3mondev/puredns/releases/tag/v2.0.0
