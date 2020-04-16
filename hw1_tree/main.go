package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var common = `├───`
var last = `└───`
var separator = "│\t"
var lastSeparator = "\t"


func filterFiles(allFiles []os.FileInfo, includeFiles bool) (ret []os.FileInfo) {
    for _, file := range allFiles {
        if file.IsDir() || includeFiles {
            ret = append(ret, file)
        }
    }
    return
}

func createTree(root string, printFiles bool) string {
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return ""
	}
	files = filterFiles(files, printFiles)

	numFiles := len(files)
	var builder strings.Builder

	for ind, f := range files {
		// Add current directory
		sep := separator
		if ind == numFiles - 1 {
			sep = lastSeparator
			builder.WriteString(last)
		} else {
			builder.WriteString(common)
		}
		builder.WriteString(f.Name())
		
		if f.IsDir() {
			builder.WriteString("\n")
			
			// Retrieve directory's inner structure
			ret := createTree(root + "/" + f.Name(), printFiles)
			rows := strings.Split(ret, "\n")
			rows = rows[:len(rows) - 1]
			for _, row := range rows {
				builder.WriteString(sep + row + "\n")
			}

		} else {
			// Just append the files info
			size := f.Size()
			sizeStr := " (empty)\n"
			if size != 0 {
				sizeStr = fmt.Sprintf(" (%db)\n", size)
			}
			builder.WriteString(sizeStr)
		}
	}
	return builder.String()
}

func dirTree(output io.Writer, path string, printFiles bool) error {
	result := createTree(path, printFiles)
	if result != "" {
		_, _ = fmt.Fprint(output, result)
		return nil
	}
	return fmt.Errorf("something went wrong")
}

func main() {
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}

	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"

	err := dirTree(os.Stdout, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
