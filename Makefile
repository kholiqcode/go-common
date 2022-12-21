.PHONY:

# ==============================================================================
# Common support

test:
	go clean -testcache && go test -v -cover -race -coverprofile=coverage.out -covermode=atomic ./...

# ==============================================================================
# Mocks

mocks:
	mockgen -package s3_mock -destination pkg/s3/mock/s3_client_mock.go -source=pkg/s3/s3_client.go