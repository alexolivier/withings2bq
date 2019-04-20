build:
	dep ensure --vendor-only
	docker build .
	rm -rf vendor