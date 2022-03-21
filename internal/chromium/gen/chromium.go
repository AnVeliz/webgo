//go:build ignore

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	outFile := os.Args[2]
	writer := getWriter(outFile)

	fmt.Fprintln(writer, "//Auto generated chromium content")
	fmt.Fprintln(writer, "//Do not modify the file manually")
	fmt.Fprintln(writer, "")
	fmt.Fprintln(writer, "package chromium")
	fmt.Fprintln(writer, "")

	fmt.Fprintln(writer, "var Files = []string{")
	err := filepath.Walk(
		"../../assets/chromium/",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				clearName := filepath.ToSlash(fmt.Sprintf("\"%s\",", strings.TrimLeft(path, "..\\..\\")))
				fmt.Fprintln(writer, fmt.Sprintf("\t%s", clearName))
			}

			return nil
		},
	)
	fmt.Fprintln(writer, "}")

	if err != nil {
		panic(err)
	}

	fmt.Fprintln(writer, "")
	fmt.Fprintln(writer, "//Auto generated chromium content. End of file.")
}

func getWriter(outFile string) *os.File {
	writer := os.Stdout
	var err error
	if len(os.Args) == 3 {
		writer, err = os.Create(outFile)
		if err != nil {
			panic(err)
		}
	}

	return writer
}
