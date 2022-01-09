package main

import (
	"flag"
	"fmt"
	"go/ast"
	"os"
	"strings"

	"github.com/tony2001/go-list-func"
	"github.com/tony2001/go-untested"
	"golang.org/x/tools/go/packages"
)

func main() {
	var (
		showTested    bool
		showPrivate   bool
		showGenerated bool
		verbose       bool
	)

	flag.BoolVar(&showTested, "tested", false, "show tested funcs")
	flag.BoolVar(&showPrivate, "private", false, "show untested unexported funcs")
	flag.BoolVar(&showGenerated, "generated", false, "show untested generated funcs")
	flag.BoolVar(&verbose, "verbose", false, "show full func names")
	flag.Parse()

	pkgs, err := list.LoadPackages(flag.Args(), true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "LoadPackages(): %v\n", err)
		os.Exit(1)
	}

	pkgMap := make(map[string]untested.Pkg, len(pkgs))

	applyFunc := func(parsedPkg *packages.Package, file *ast.File, decl *ast.FuncDecl) error {
		if parsedPkg.Name == "main" { // XXX ?
			return nil
		}

		pkg, ok := pkgMap[parsedPkg.Name]
		if !ok {
			pkg.Name = parsedPkg.Name
		}

		// contain notest comment - it's neither tested, nor untested
		if untested.HasNoTestComment(decl) {
			return nil
		}

		fn := untested.Func{
			Name:        list.FormatFuncDecl("", decl),
			VerboseName: list.FormatFuncDeclVerbose("", decl),
		}

		if !showPrivate && !untested.IsExported(fn) {
			return nil
		}

		if !showGenerated && untested.IsGeneratedFile(file) {
			return nil
		}

		if strings.HasPrefix(fn.Name, "Test") {
			pkg.Tests = append(pkg.Tests, fn)
		} else if !strings.HasSuffix(parsedPkg.Name, untested.TestPkgSuffix) {
			pkg.Funcs = append(pkg.Funcs, fn)
		}

		pkgMap[parsedPkg.Name] = pkg
		return nil
	}

	if err = list.WalkFuncs(pkgs, applyFunc); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	testedFuncs, untestedFuncs := untested.Analyze(pkgMap)

	var outList []untested.Func
	if showTested {
		outList = testedFuncs
	} else {
		outList = untestedFuncs
	}

	for _, fn := range outList {
		if verbose {
			fmt.Println(fn.VerboseName)
		} else {
			fmt.Println(fn.Name)
		}
	}
}
