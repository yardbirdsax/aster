// Package aster_test contains all the tests for the `aster` module.
package aster_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yardbirdsax/aster"
)

func TestGetPackageComment(t *testing.T) {
  tests := []struct{
    name string
    directory string
    packages []string
    expectedPackageComment string
  }{
    {
      name: "one package",
      directory: ".",
      packages: []string{"aster"},
      expectedPackageComment: "Package aster provides a high level interface for parsing Go code using the stdlib's\n" +
        "[`go/ast`](https://pkg.go.dev/go/ast) module.\n",
    },
    {
      name: "multi package",
      directory: ".",
      packages: []string{"aster", "aster_test"},
      expectedPackageComment: "aster:\nPackage aster provides a high level interface for parsing Go code using the stdlib's\n" +
        "[`go/ast`](https://pkg.go.dev/go/ast) module.\n" +
        "aster_test:\nPackage aster_test contains all the tests for the `aster` module.\n",
    },
  }

  for _, tc := range tests {
    t.Run(tc.name, func(t *testing.T) {
      tc := tc
      t.Parallel()

      actualPackageComment := aster.FromDirectory(tc.directory).Packages(tc.packages).PackageComment()

      assert.Equal(t, tc.expectedPackageComment, actualPackageComment)
    })
  }
}

func TestMatchComment(t *testing.T) {
  tests := []struct{
    name string
    directory string
    comment string
    expectedResults []aster.Result
  }{
    {
      name: "with comment",
      directory: "sample",
      comment: "aster:",
      expectedResults: []aster.Result{
        {
          Name: "Sample",
          Type: "struct",
          Comments: "aster: hello\n",
          Fields: []aster.Field{
            {
              Name: "Text",
              Type: "string",
              Comments: "Text is some sample text\n",
            },
          },
        },
        {
          Name: "Sampler",
          Type: "func",
          Comments: "aster: hello again\n",
        },
      },
    },
  }

  for _, tc := range tests {
    tc := tc
    t.Run(tc.name, func(t *testing.T) {
      t.Parallel()

      a := aster.FromDirectory(tc.directory)
      actualResults, err := a.MatchComment("aster:")

      assert.Equal(t, tc.expectedResults, actualResults)
      assert.NoError(t, a.Error)
      assert.NoError(t, err)
    })
  }
}
