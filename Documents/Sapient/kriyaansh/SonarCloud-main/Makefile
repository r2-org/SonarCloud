export GO111MODULE=on
export TF_LOG=DEBUG
SRC=$(shell find . -name '*.go')

.PHONY: clean build run install

build:
	go build -o terraform-provider-sonarcloud

run: 
	cd example
	terraform init
	terraform apply --auto-approve

clean:
	rm -rf terraform-sonarcloud .terraform terraform.tfstate crash.log terraform.tfstate.backup
