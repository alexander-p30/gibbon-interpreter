# What is this?

Gibbon is a interpreter for the Monkey pseudo-language, invented for the [Writing An Interpreter In Go](https://interpreterbook.com) book.

Why Gibbon? Well, the interpreter is implemented using the "tree-walker" pattern, because it relies on an AST (Abstract-Syntax Tree) to execute its input, and it does so by traversing this data structure from "branch" to "branch", just like a [bracchiator](https://en.wikipedia.org/wiki/Brachiator), such as gibbons, would.
