export TF_ACC=1

.PHONY: test
test:
	go test ./... -v -timeout 120m -coverprofile coverage.txt -covermode atomic
