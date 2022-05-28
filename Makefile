run: build
	./mycomarkup

build:
	go build ./cmd/mycomarkup

test1: build
	./mycomarkup -hypha-name "test doc" -file-name "testdata/test1.myco"

test2: build
	./mycomarkup -hypha-name "test doc" -file-name "testdata/test2.myco"

test3: build
	./mycomarkup -hypha-name "test doc" -file-name "testdata/test3.myco"

test4: build
	./mycomarkup -hypha-name "test doc" -file-name "testdata/test4.myco"

test_list_examples: build
	./mycomarkup -hypha-name "test doc" -file-name "testdata/list_examples.myco"

test_quotes: build
	./mycomarkup -hypha-name "test doc" -file-name "testdata/test_quotes.myco"

test_p_and_blank: build
	./mycomarkup -hypha-name "test doc" -file-name "testdata/test_p_and_blank.myco"

test_tables: build
	./mycomarkup -hypha-name "test doc" -file-name "testdata/tables.myco"

test_death: build
	./mycomarkup -hypha-name "test doc" -file-name "testdata/death.myco"

test_new_headings: build
	./mycomarkup -hypha-name "test doc" -file-name "testdata/test_new_headings.myco"

test_interwiki: build
	./mycomarkup -hypha-name "test doc" -file-name "testdata/test_interwiki.myco"

test_launchpad: build
	./mycomarkup -hypha-name "test doc" -file-name "testdata/test_launchpad.myco"

help:
	@echo "Mycomarkup: https://mycorrhiza.wiki/hypha/mycomarkup."
	@echo "Read Makefile too see what you can do."
