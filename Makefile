VERSION = v1.2.0

web:
	cd frontend && npm install
	cd frontend && npm run build

linux:
	rm -rf ./hamster-provider-$(VERSION)-linux-amd64.tar.gz
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build
	tar -czvf ./hamster-provider-$(VERSION)-linux-amd64.tar.gz ./hamster-provider ./templates ./frontend/dist

macos:
	rm -rf ./hamster-provider-$(VERSION)-darwin-amd64.tar.gz
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64  go build
	tar -czvf ./hamster-provider-$(VERSION)-darwin-amd64.tar.gz ./hamster-provider ./templates ./frontend/dist

windows:
	rm -rf ./hamster-provider-$(VERSION)-windows-amd64.tar.gz
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64  go build
	tar -czvf ./hamster-provider-$(VERSION)-windows-amd64.tar.gz ./hamster-provider.exe ./templates ./frontend/dist

all: web linux macos windows

clean:
	rm -rf ./hamster-provider.exe
	rm -f ./hamster-provider
	rm -rf ./hamster-provider-*
