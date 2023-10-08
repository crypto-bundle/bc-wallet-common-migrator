deploy:
	$(if $(and $(env),$(repository)),,$(error 'env' and/or 'repository' is not defined))

	$(eval build_tag=$(env)-$(shell git rev-parse --short HEAD)-$(shell date +%s))
	$(eval container_registry=$(repository)/crypto-bundle/bc-wallet-common-migrator)
	$(eval platform=$(or $(platform),linux/amd64))

	docker build \
		--ssh default=$(SSH_AUTH_SOCK) \
		--no-cache \
		--platform $(platform) \
		--tag $(container_registry):$(build_tag) . \
		--tag $(container_registry)

	docker push $(container_registry):$(build_tag)
	docker push $(container_registry)

.PHONY: deploy