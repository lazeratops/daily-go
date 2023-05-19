mkdir local
v=`go version | { read _ _ v _; echo ${v#go}; }`

go test ./... -coverpkg=./... -coverprofile=local/coverage.out -covermode=atomic -count=1
go tool cover -html=local/coverage.out -o local/coverage-$v.html