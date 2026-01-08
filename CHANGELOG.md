# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [v2.1.0] - 2026-01-08
### Added
- Added `-sort` option to order ports by process (default), port, pid, protocol, address, or status.
- Added sorting coverage tests.

### Changed
- Updated the Go module path to `github.com/sinansonmez/open-ports` for pkg.go.dev publishing.

## [v2.0.0] - 2026-01-08
### Added
- Added table rendering tests for headers and row values.

### Fixed
- Fixed Windows builds by using platform-specific process termination helpers.

### Changed
- Run `make test` in the release workflow before building archives.

## [v1.0.0] - 2026-01-08
### Added
- Created repository.

[v2.1.0]: https://github.com/sinansonmez/open-ports/releases/tag/v2.1.0
[v2.0.0]: https://github.com/sinansonmez/open-ports/releases/tag/v2.0.0
[v1.0.0]: https://github.com/sinansonmez/open-ports/releases/tag/v1.0.0
