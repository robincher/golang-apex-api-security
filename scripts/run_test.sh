go get -u github.com/GovTechSG/test-suites-apex-api-security

go test ./ -v -cover -coverprofile=./output/c.out

go tool cover -html=./output/c.out -o ./output/coverage.html