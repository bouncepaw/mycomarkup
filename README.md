# Mycomarkup
[![godocs.io](http://godocs.io/github.com/bouncepaw/mycomarkup?status.svg)](http://godocs.io/github.com/bouncepaw/mycomarkup)

**Mycomarkup** is a markup language designed to be used in [Mycorrhiza Wiki](https://mycorrhiza.lesarbr.es). This project is both a library to be used in the wiki engine and a command-line tool for processing mycomarkup files in other projects.

See [the Mycorrhiza docs](https://mycorrhiza.lesarbr.es/hypha/mycomarkup) on the markup language itself. The rest of the document provides documentation on the library and the command only.

Also see [our kanban board](https://github.com/bouncepaw/mycomarkup/projects/1).

## Running
```
Usage of mycomarkup:
  -filename string
        File with mycomarkup. (default "/dev/stdin")
  -hypha-name string
        Set hypha name. Relative links depend on it.
```

Set the parameters and run the program. The output will be written to `stdout`. The output is a poorly-formatted HTML code. In the future, more front-ends will be available.

Please note that transclusion is not supported in CLI.

## Embedding
Current API for using Mycomarkup in other Go projects is highly unstable, but it does exist. Consult the docs and Mycorrhiza Wiki source code for inspiration.

## Contributing
All pull requests are welcome. Feel free to open issues. Also, pay a visit to the [MycorrhizaWiki Telegram chat](https://t.me/mycorrhizadev). Also consider donating on [Boosty](https://boosty.to/bouncepaw).
