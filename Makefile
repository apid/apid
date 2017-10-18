fresh: update build
build:
	@(go build  -o helper/buildhelper helper/buildhelper.go;)
	@(go test github.com/apid/apid/helper;)
	@(./helper/buildhelper ./glide.lock 2>builderr 1>buildapid; rm helper/buildhelper;)
	@(test -s buildapid || { echo "build script generation failed"; rm builderr buildapid; exit 1;})
	@(chmod +x ./buildapid; echo "building apid..."; ./buildapid; rm builderr buildapid)
update:
	@(rm glide.lock; glide update -v)
