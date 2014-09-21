test:
	go get github.com/smartystreets/goconvey
	./test.sh
	
deps:
	go get -u ./...
