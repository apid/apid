fresh: update build
build:
	@(mv build_helper.go.1 build_helper.go; go build build_helper.go; mv build_helper.go build_helper.go.1)
	@(./build_helper 2>builderr 1>buildapid; rm build_helper)
	@test -s buildapid || { echo "build script generation failed"; rm builderr buildapid; exit 1;}
	@(chmod +x ./buildapid; ./buildapid; rm builderr buildapid; echo "build complete")
update:
	@(rm glide.lock; glide update -v)
