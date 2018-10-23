go test ./ -v -cover -coverprofile=./output/c.out

go tool cover -html=./output/c.out -o ./output/coverage.html