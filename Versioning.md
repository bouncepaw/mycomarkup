# Versioning

Mycomarkup uses [https://semver.org](SemVer), or at least pretends to.

It means that versions follow this pattern: MAJOR.MINOR.PATCH.

Increment MAJOR when you make an API-incompatible change. In our case, it means that users would have to change their code that looks like that:

```go
ctx, _ := mycocontext.ContextFromStringInput(name, content)
ast := mycomarkup.BlockTree(ctx)
result := mycomarkup.BlocksToHTML(ctx, ast)
```

Mycomarkup has a lot more exported symbols than that, and they change often. They are not part of API. One day we should use Go's internal packages for a clearer distinction.

Increment MINOR when you _add_ something new, such as a new block.

Increment PATCH when you _fix_ something, such as a regular bug. Increment when _refactoring_ something too, because bad code is something like a bug.

Each git tag must start with `v`: `v1.0.1`, `v7.8.0`. It is quite ok if there is a tag for each consecutive commit, if it makes sense (it often does).