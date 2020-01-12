prog = app
test:
	go test ./...
build: test
	go build -o $(prog) -ldflags "\
  		-X github.com/longkai/xiaolongtongxue.com/context.v=`git rev-parse --short HEAD` \
  		-X github.com/longkai/xiaolongtongxue.com/context.b=`git rev-parse --abbrev-ref HEAD` \
	"
clean:
	go clean
	[ -f $(prog) ] && rm $(prog) || :