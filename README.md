# Cosmos Go Transaction

Cosmos SDK Transaction using pass for key management.

Shows building and signing a tx, which tbh is long winded.

## Installation

- `go mod download`
- `go run main.go`
- `go build`

## Usage

The pass backend uses the appName passed to `keyring.New` to create a namespace
in pass. In the case that the appName is `test` it will look for keys under
`keyring-test`.

```sh
pass show | grep -C 2 bot-1

├── keyring-test
│   ├── bot-1.info

```

```sh
./cosmos-tx
```

The process writes data to the console.
