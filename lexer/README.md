# Lexer
A mycodocument is a fractal thing.

## Nesters
Consecutive lines consisting of these characters are taken together, stripped from these characters and then lexed again. For example, we have this document:

```myco
paragraph

> line
> > line line line
> # hey hey
```

First, the paragraph is seen. Then, we see three lines starting with `>`. The lexer zooms in:

```myco
line
> line line line
# hey hey
```

The paragraph and the heading are easy, but the quote needs further zooming. You get it.

Lists have some tricky rules.

## What can be met
One-liners:
* Horizontal lines
* Headings

Togglers:
* Preformatted

Consecutives:
* Paragraphs
* Rocket links

Nesters:
* Unnumbered lists
* Numbered lists
* Blockquotes
* Most curly tags ig
