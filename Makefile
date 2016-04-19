package = github.com/abeMedia/push-deploy

.PHONY: release

release:
    go get
	mkdir -p release
	GOOS=linux GOARCH=amd64 go build -o release/push-deploy-linux-amd64 $(package)
	GOOS=linux GOARCH=386 go build -o release/push-deploy-linux-386 $(package)
	GOOS=linux GOARCH=arm go build -o release/push-deploy-linux-arm $(package)