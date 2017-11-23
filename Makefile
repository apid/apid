APIDGOPACKAGES = $(shell go list ./... | grep -v vendor)
APIDGOFILES = $(shell git ls-files | grep '.go$$' | grep -v vendor)
GOFMTCOUNT = $(shell gofmt -l $(APIDGOFILES) | wc -l)
fresh: update build
build: fmtcheck
	@(go build  -o helper/buildhelper helper/buildhelper.go;)
	@(go test github.com/apid/apid/helper;)
	@(./helper/buildhelper ./glide.lock 2>builderr 1>buildapid; rm helper/buildhelper;)
	@(test -s buildapid || { echo "build script generation failed"; rm builderr buildapid; exit 1;})
	@(chmod +x ./buildapid; echo "building apid..."; ./buildapid; rm builderr buildapid)
update:
	@(rm glide.lock; glide update -v)
fmtcheck:
	go vet $(APIDGOPACKAGES)
	@(if [ $(GOFMTCOUNT) -gt 0 ]; then echo 'Run "go fmt $$(glide novendor)" on your go source code'; exit 1; else echo "Go source format look good"; fi)
