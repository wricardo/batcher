test: deps
	./test.sh
	
deps:
	go get github.com/smartystreets/goconvey/convey
	go get github.com/garyburd/redigo/redis
