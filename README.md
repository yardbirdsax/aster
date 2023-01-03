# aster

`aster` is Go library exposing a high level interface for parsing Go code with
[`go/ast`](https://pkg.go.dev/go/ast). It is named after the ship of the character Ronan the Accuser
in the Marvel Cinematic Universe, the [_Dark
Aster_](https://marvelcinematicuniverse.fandom.com/wiki/Dark_Aster).

## Background

Right now `aster` only has a very limited set of functionality that met my immediate needs. PRs for
more are welcome!

I needed the ability to retrieve objects from code that had comments matching certain patterns
attached to them, similar to how [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) works
with [code markers](https://book.kubebuilder.io/reference/markers.html). I knew `ast` could do this,
but found the interface very complex. What was a [`Decl`](https://pkg.go.dev/go/ast#Decl) and when
is something a `GenDecl` vs `FuncDecl`? How do I figure out what the type of an object is (i.e.
`struct` vs `func` vs something else)? And once I actually retrieve the right `Decl`, how do I get
it's name or what its members / fields are?

The goal of `aster` was to be able to say "here's a pattern, find me all objects that have
comments matching it attached, and tell me what they are and, if applicable, what their fields are".
For a simple example, see the [`astersample`](cmd/astersample/main.go) program.
