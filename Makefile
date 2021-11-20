.PHONY: run tidy mock mock-prepare

run:
	go run cmd/main.go -env=local

tidy:
	go mod tidy

mock-prepare:
	go install github.com/golang/mock/mockgen@v1.6.0
	go get -u github.com/golang/mock/gomock
	go get -u github.com/bxcodec/faker/v3

mock:
	mockgen -source=service/interface.go -destination=service/mock/interface_mock.go -package=mock
	mockgen -source=entity/user/interface.go -destination=entity/user/mock/interface_mock.go -package=mock
