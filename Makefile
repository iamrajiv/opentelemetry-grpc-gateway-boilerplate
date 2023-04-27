clean:
	rm proto/helloworld/v1/*.go

generate:
	buf generate proto/helloworld/v1/*.proto
