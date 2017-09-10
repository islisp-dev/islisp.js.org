VERSION := $(shell cd $(GOPATH)/src/github.com/ta2gch/iris/ && git rev-parse --short HEAD)
all:
	gsed -i -e "s/version\s*=\s*\".*\"/version = \"$(VERSION)\"/g" main.go
	gopherjs build main.go
	git commit -a -m 'Updated the implementation'
	git push
