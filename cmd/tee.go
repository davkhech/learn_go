package main

import (
	"bufio"
	"github.com/akamensky/argparse"
	"os"
)


func check(err error) {
	if err != nil {
		panic(err)
	}
}

func OpenFiles(filenames []string, appendFile bool) []*os.File{
	flag := os.O_CREATE | os.O_WRONLY
	if appendFile {
		flag |= os.O_APPEND
	} else {
		flag |= os.O_TRUNC
	}

	files := make([]*os.File, 0)
	for _, filename := range filenames {
		file, _ := os.OpenFile(filename, flag, 0666)
		files = append(files, file)
	}
	return files
}

func CloseFiles(files []*os.File) {
	for _, file := range files {
		file.Close()
	}
}

func main() {
	parser := argparse.NewParser("tee", "Duplicate of tee command")
	filenames := parser.List("", "file", &argparse.Options{Required: true, Help: "File names to write to"})
	appendFile := parser.Flag("a", "append", &argparse.Options{Required: false, Help: "File write mode"})

	err := parser.Parse(os.Args)
	check(err)

	files := OpenFiles(*filenames, *appendFile)
	defer CloseFiles(files)
	files = append(files, os.Stdout)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "exit" {
			break
		}

		for _, file := range files {
			file.WriteString(text + "\n")
		}
	}
}