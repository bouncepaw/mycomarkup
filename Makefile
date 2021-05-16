run: build
	./mycomarkup

build:
	go build .

test1: build
	./mycomarkup -hypha-name "test doc" -filename "testdata/test1.myco"

generate:
	# stringer -type=TokenKind -trimprefix Token ./lexer

indent:
	go run . | indent

help:
	@echo "Mycomarkup: https://mycorrhiza.lesarbr.es/hypha/mycomarkup."
	@echo "Read Makefile too see what you can do."
