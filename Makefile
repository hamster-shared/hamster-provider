VERSION = v1.3.0

web:
	cd frontend && npm install
	cd frontend && npm run build

linux:
	rm -rf ./hamster-provider-$(VERSION)-linux-amd64.tar.gz
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build
	tar -czvf ./hamster-provider-$(VERSION)-linux-amd64.tar.gz ./hamster-provider ./templates ./frontend/dist

macos:
	cd frontend && yarn && yarn build
	rm -rf core/corehttp/dist
	cp -r frontend/dist core/corehttp/
	rm -rf ./hamster-provider-$(VERSION)-darwin-amd64.tar.gz
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64  go build -o build/bin
 	# gon -log-level=info ./build/darwin/gon-sign.json
	#tar -czvf ./hamster-provider-$(VERSION)-darwin-amd64.tar.gz ./hamster-provider ./templates ./frontend/dist

windows:
	go get github.com/akavel/rsrc
	rm -rf ./hamster-provider-$(VERSION)-windows-amd64.zip
	rsrc -arch amd64 -manifest build/windows/hamster-provider.exe.manifest -ico build/windows/icon.ico -o main.syso
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64  go build  -ldflags "-s -w" -o hamster-provider.exe
	zip -r ./hamster-provider-$(VERSION)-windows-amd64.zip ./hamster-provider.exe

docker:
	docker build -t hamstershare/hamster-provider:$(VERSION) .
	docker push hamstershare/hamster-provider:$(VERSION)

all: web linux macos windows

clean:
	rm -rf ./hamster-provider.exe
	rm -f ./hamster-provider
	rm -rf ./hamster-provider-*
