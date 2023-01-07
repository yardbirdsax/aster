/*
Package aster provides a high level interface for parsing Go code using the stdlib's
[`go/ast`](https://pkg.go.dev/go/ast) module.
*/
package aster

import (
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"sort"
	"strings"
)

// FromDirectory creates a new `aster` struct from the contents of the specified directory.
func FromDirectory(path string) *Aster {
	a := &Aster{
		fileset: token.NewFileSet(),
	}

	a.packages, a.Error = parser.ParseDir(a.fileset, path, nil, parser.ParseComments)

	return a
}

/*
Aster is the main struct and pipeline object for the `aster` module. It provides all
the methods and interfaces used by any wrapper methods at the package level.
*/
type Aster struct {
	// Error is populated when an error occurs in the course of processing results.
	Error error
	// fileset is the `token.Fileset` object used when parsing all the underyling Go code.
	fileset *token.FileSet
	// packages is the map of `ast.Package` structs populated from the parsed Go files.
	packages map[string]*ast.Package
}

// Packages filters the returned packages to only those with the specified names.
func (a *Aster) Packages(names []string) *Aster {
	newPackages := map[string]*ast.Package{}

	for pkgName, pkg := range a.packages {
		for _, name := range names {
			if pkgName == name {
				newPackages[pkgName] = pkg
				continue
			}
		}
	}
	a.packages = newPackages

	return a
}

// PackageComment returns the comments for all parsed packages. If more than one is present,
// it will return a concatenated group of them with each one prefaced with the package name.
// If an error has occurred during processing, then it will return a blank string.
func (a *Aster) PackageComment() string {
	if a.Error != nil {
		return ""
	}

	packageComments := []string{}

	for pkgName, pkg := range a.packages {
		var b strings.Builder
		if len(a.packages) > 1 {
			b.WriteString(pkgName + ":\n")
		}
		for _, f := range pkg.Files {
			b.WriteString(f.Doc.Text())
		}

		packageComments = append(packageComments, b.String())
	}

	sort.Strings(packageComments)

	return strings.Join(packageComments, "")
}

// MatchComment returns all declarations that have a related comment matching the given string.
// If an error has occurred during processing, it returns an empty slice.
func (a *Aster) MatchComment(match string) ([]Result, error) {
	results := []Result{}
	if a.Error != nil {
		return results, a.Error
	}

	r, err := regexp.Compile(match)
	if err != nil {
		a.Error = err
		return results, a.Error
	}

	for _, pkg := range a.packages {
		for _, f := range pkg.Files {
			matchingCommentGroups := []*ast.CommentGroup{}
			for _, cGroup := range f.Comments {
				if r.MatchString(cGroup.Text()) {
					matchingCommentGroups = append(matchingCommentGroups, cGroup)
					continue
				}
			}

			if len(matchingCommentGroups) > 0 {
				for _, cGroup := range matchingCommentGroups {
					cGroupLine := a.fileset.Position(cGroup.End()).Line

					for _, decl := range f.Decls {
						declLine := a.fileset.Position(decl.Pos()).Line
						if cGroupLine == declLine-1 {
							result, err := resultFromDecl(decl)
							if err != nil {
								a.Error = err
								return results, a.Error
							}
							results = append(results, result)
						}
					}
				}
			}
		}
	}

	return results, nil
}
