run: build
	./mycomarkup

build:
	go build ./cmd/mycomarkup

test1: build
	./mycomarkup -hypha-name "test doc" -filename "testdata/test1.myco"

test2: build
	./mycomarkup -hypha-name "test doc" -filename "testdata/test2.myco"

test3: build
	./mycomarkup -hypha-name "test doc" -filename "testdata/test3.myco"

test_list_examples: build
	./mycomarkup -hypha-name "test doc" -filename "testdata/list_examples.myco"

test_quotes: build
	./mycomarkup -hypha-name "test doc" -filename "testdata/test_quotes.myco"

test_p_and_blank: build
	./mycomarkup -hypha-name "test doc" -filename "testdata/test_p_and_blank.myco"

test_tables: build
	./mycomarkup -hypha-name "test doc" -filename "testdata/tables.myco"


generate:
	# stringer -type=TokenKind -trimprefix Token ./parser

indent:
	go run . | indent

help:
	@echo "Mycomarkup: https://mycorrhiza.lesarbr.es/hypha/mycomarkup."
	@echo "Read Makefile too see what you can do."
