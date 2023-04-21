# Bc Wallet Common Migrator

### Install
```bash
go install gitlab.heronodes.io/bc-platform/bc-wallet-common-migrator/cmd@latest
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
