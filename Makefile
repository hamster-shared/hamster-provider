linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build
	tar -czvf ./hamster-provider-linux-amd64.tar.gz ./hamster-provider ./templates ./frontend/dist

clean:
	rm -f ./hamster-provider
