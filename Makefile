run/Tasky:
	@go run ./cli/tasky

build/Tasky:
	@go build ./cli/tasky

test/Tasky:
	@cd cli/tasky && go test -v ./...