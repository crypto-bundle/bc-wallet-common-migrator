build:
	$(eval build_tag=$(env)-$(shell git rev-parse --short HEAD)-$(shell date +%s))
	$(eval container_registry=$(repository)/bc-platform/bc-wallet-common-migrator)

	docker buildx build --no-cache --ssh default=$(SSH_AUTH_SOCK) --platform linux/amd64,linux/arm64 --push -t $(container_registry):$(build_tag) .

.PHONY: build