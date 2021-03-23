run: build
	./mycomarkup

build: generate
	go build .

generate:
	stringer -type=TokenKind -trimprefix Token ./lexer

indent:
	go run . | indent

help:
	@echo "Mycomarkup: https://mycorrhiza.lesarbr.es/hypha/mycomarkup."
	@echo "Read Makefile too see what you can do."
