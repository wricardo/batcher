test: deps
	go test -timeout 5s -covermode=count -coverprofile=coverage.out ./...
	
deps:
	go get github.com/smartystreets/goconvey/convey
	go get github.com/garyburd/redigo/redis
