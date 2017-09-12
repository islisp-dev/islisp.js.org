VERSION := $(shell cd $(GOPATH)/src/github.com/ta2gch/iris/ && git rev-parse --short HEAD)
all:
	gsed -i -e "s/version\s*=\s*\".*\"/version = \"$(VERSION)\"/g" script/main.go api/eval.go
	gopherjs build script/main.go -o script/main.js
	gopherjs build api/eval.go -o api/eval.js
	git commit -a -m 'Updated the implementation'
	git push
