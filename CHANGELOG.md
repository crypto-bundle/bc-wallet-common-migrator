# Change Log

## [v0.0.6] - 11.03.2025
### Added
* Added Helm Job chart
### Changed
* Added support of Go 1.23
* Updated License
  * Copyright - new year 2025
  * Added License banner to *.go files
* Added support of last version of lib-errors/lib-tinyerrors
* Added linter and fixed up all linter issues
* Bump goose-migrator version
* Bump common-lib-* versions:
  * bc-wallet-common-lib-vault v0.0.20
  * bc-wallet-common-lib-config v0.0.8
  * bc-wallet-common-lib-errors v0.0.10
  * bc-wallet-common-lib-logger v0.0.12
  * bc-wallet-common-lib-postgres v0.0.11

## [v0.0.5] - 02.09.2024
### Added
* Added migrator library Helm-chart
* Added support of new version of bc-wallet-common-lib-errors
* Added support of new version of bc-wallet-common-lib-logger
* Bump common-lib-* versions:
  * bc-wallet-common-lib-vault v0.0.19
  * bc-wallet-common-lib-config v0.0.7
  * bc-wallet-common-lib-postgres v0.0.10

## [v0.0.4] - 21.06.2024
### Changed
* Changed config init flow
* Changed MIT License to NON-AI MIT.
* Bump common-lib-vault version - bc-wallet-common-lib-vault v0.0.17
  * Added new environment variables VAULT_AUTH_TOKEN_FILE_PATH and VAULT_AUTH_TOKEN_RENEW_TTL to *-example.env files
  * Implemented new vault client init flow

## [v0.0.3] - 17.06.2024
### Changed
* Changed CLI arguments init flow
* Updated lib-vault version to 0.0.15

## [v0.0.2] - 16.04.2024
### Added
* Migrator moved to another repository - https://github.com/crypto-bundle/bc-wallet-common-migrator
* Bump golang version 1.19 -> 1.22
* Added support of ld-flags in build stage
* Replaced common libs
  * bc-wallet-common-lib-config
  * bc-wallet-common-lib-logger
  * bc-wallet-common-lib-postgres
  * bc-wallet-common-lib-vault

## [v0.0.1] - 11.05.2023
### Added
* Added goose migration tool wrapper
* Added build flow for docker container and local development