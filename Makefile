test:
	go test ./... -gcflags=all=-l -cover
	go vet ./...
	golangci-lint run ./... || true
	govulncheck ./...

test2file:
	go test ./... -v -gcflags=all=-l -json > sn_report_test.json
	go test ./... -gcflags=all=-l -coverprofile=sn_report_covprofile
	go vet -json ./... 2> sn_report_vet_report.out
	golangci-lint run --out-format checkstyle ./... > sn_report_report.xml || true

init:
	go mod tidy
	go generate ./...

clean:
	rm -rf ./bin