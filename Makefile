all: version main.js api/eval.js
	git commit -a -m 'Updated the implementation'
	git push

main.js: main.go
	gopherjs build main.go -o main.js

api/eval.js: api/eval.go
	gopherjs build api/eval.go -o api/eval.js

.PHONY: version

VERSION := $(shell cd $(GOPATH)/src/github.com/ta2gch/iris/ && git rev-parse --short HEAD)

version:
	gsed -i -e "s/version\s*=\s*\".*\"/version = \"$(VERSION)\"/g" main.go api/eval.go
