mkdir local
go test ./... -coverpkg=./... -coverprofile=local/coverage.out -covermode=atomic -count=1
go tool cover -html=local/coverage.out -o local/coverage.html