go test -coverprofile coverageResult.out ./...
go tool cover -html=coverageResult.out -o coverageResult.html
clear
open coverageResult.html