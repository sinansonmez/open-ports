# open-ports

`open-ports` provides a simple command-line workflow for inspecting open network ports
from your terminal.

<img src="https://github.com/user-attachments/assets/97160fca-832e-46a0-9f81-f79eeef8c57c" width="400" height="600" />

## Installation

```sh
brew tap sinansonmez/tools
brew install open-ports
open-ports -h
```

## Usage

Run the CLI to see available commands and flags:

```sh
open-ports -h
```

Sort by a specific field (default: process):

```sh
open-ports -sort process
```

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for release history.

## Development

Build and run the CLI locally:

```sh
make build
./bin/open-ports -h
```
