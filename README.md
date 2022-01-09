# go-untested

go-untested lists untested functions in the package(s).

## Example
```
$ go-untested -verbose ./...
func IsExported(fn Func) bool
func Analyze(pkgs map[string]Pkg) (tested []Func, untested []Func)
func IsGeneratedFile(file *ast.File) bool
func HasNoTestComment(decl *ast.FuncDecl) bool

```

## CLI
```
Usage of go-untested:
  -generated
        show untested generated funcs
  -private
        show untested unexported funcs
  -tested
        show tested funcs
  -verbose
        show full func names
```

