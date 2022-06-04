# Mycomarkup
[![godocs.io](http://godocs.io/github.com/bouncepaw/mycomarkup/v5?status.svg)](http://godocs.io/github.com/bouncepaw/mycomarkup/v5)

**Mycomarkup** is a markup language designed to be used in [Mycorrhiza Wiki](https://mycorrhiza.wiki). This project is
both a library for the wiki engine and a command-line tool for processing Mycomarkup files in other projects.

See [the Mycorrhiza docs](https://mycorrhiza.wiki/help/en/mycomarkup) on the markup language itself. The rest of the document provides documentation on the library and the command only.

Also see [our kanban board](https://github.com/bouncepaw/mycomarkup/projects/1) and [this hypha](https://mycorrhiza.wiki/hypha/release/roadmap).

## Running
```
Usage of mycomarkup:
  -file-name string
        File with mycomarkup. (default "/dev/stdin")
  -hypha-name string
        Set hypha name. Relative links depend on it.
```

Set the parameters and run the program. The output will be written to `stdout`. The output is a poorly-formatted HTML code. In the future, more front-ends will be available.

Please note that transclusion is not supported in CLI.

## Embedding
Mycomarkup provides an API for Go projects. Consult the docs and Mycorrhiza Wiki source code for inspiration.

## Contributing
All pull requests are welcome. Feel free to open issues. Also, pay a visit to the [Mycorrhiza Wiki Telegram chat](https://t.me/mycorrhizadev). Also consider donating on [Boosty](https://boosty.to/bouncepaw).
