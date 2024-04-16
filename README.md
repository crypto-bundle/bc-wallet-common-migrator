# Bc Wallet Common Migrator

## Description

Bc-wallet-common-migrator is a database migration tool based on the [Goose library](https://github.com/pressly/goose)

### Install
```bash
go install github.com/crypto-bundle/bc-wallet-common-migrator/cmd@latest
```
This will install the bc-wallet-common-migrator binary to your $GOPATH/bin directory.

### Usage
```
Usege:
    bc-wallet-common-migrator [OPTIONS] GOOSE_COMMAND
    
Options:
    -dir string
        directory with migration files (default "./migrations")
    -envPath string
        environment file with migration settings
```

### Development
```bash
go build -o bc-wallet-common-migrator ./cmd/... && mv bc-wallet-common-migrator $GOPATH/bin/
```


## Contributors

* Maintainer - [@gudron (Alex V Kotelnikov)](https://github.com/gudron)
* Author [@d.burnyshev (Dmitry R Burnyshev)](https://github.com/qrinef)

## Licence

**bc-wallet-common-migrator** is licensed under the [MIT](./LICENSE) License.