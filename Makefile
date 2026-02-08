.PHONY: build

build:
	@scripts/build.sh

deploy:
	cd terraform && terraform apply -auto-approve
