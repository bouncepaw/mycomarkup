run: build
	./mycomarkup

build:
	go build .

test1: build
	./mycomarkup -hypha-name "test doc" -filename "testdata/test1.myco"

test2: build
	./mycomarkup -hypha-name "test doc" -filename "testdata/test2.myco"

test_list_examples: build
	./mycomarkup -hypha-name "test doc" -filename "testdata/list_examples.myco"

generate:
	# stringer -type=TokenKind -trimprefix Token ./parser

indent:
	go run . | indent

help:
	@echo "Mycomarkup: https://mycorrhiza.lesarbr.es/hypha/mycomarkup."
	@echo "Read Makefile too see what you can do."
