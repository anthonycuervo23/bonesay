.PHONY: build
build: build/bonesay build/cowthink

.PHONY: build/bonesay
build/bonesay:
	CGO_ENABLED=0 cd cmd && go build -o ../bin/bonesay -ldflags "-w -s" ./bonesay

.PHONY: build/cowthink
build/cowthink:
	CGO_ENABLED=0 cd cmd && go build -o ../bin/cowthink -ldflags "-w -s" ./cowthink

.PHONY: lint
lint:
	golint ./...
	cd cmd && golint ./...

.PHONY: vet
vet:
	go vet ./...
	cd cmd && go vet ./...

.PHONY: test
test: test/pkg test/cli

.PHONY: test/pkg
test/pkg:
	go test ./...

.PHONY: test/cli
test/cli:
	cd cmd && go test ./...

.PHONY: man
man:
	asciidoctor --doctype manpage --backend manpage doc/bonesay.1.txt.tpl -o doc/bonesay.1

.PHONY: man/preview
man/preview:
	cat doc/bonesay.1 | groff -man -Tascii | less
