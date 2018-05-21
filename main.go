package main

import (
	"flag"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jung-kurt/gofpdf"

	"github.com/k0kubun/pp"
)

type filesPath []string

func (i *filesPath) String() string {
	return "list string"
}

func (i *filesPath) Set(value string) error {
	pp.Println(value)
	*i = strings.Split(value, ",")
	return nil
}

func main() {
	var err error
	var files []string
	var out string

	if os.Args[1] != "-files" {
		pp.Println("Use flag -files path1 path2...")
		os.Exit(1)
	}

	for i := 2; i < len(os.Args) && os.Args[i] != "-out"; i++ {
		files = append(files, os.Args[i])
	}

	// start a new pdf file
	pdfFile := gofpdf.New("P", "mm", "A4", "")

	if len(files) <= 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	for i, arg := range os.Args {
		if arg == "-out" && i+1 < len(os.Args) {
			out = os.Args[i+1]
			break
		}
		out, err = os.Getwd()
		if err != nil {
			pp.Fatal(err)
		}
	}

	for _, f := range files {
		matched, err := regexp.MatchString(`\.png$|\.jpg$|\.svg$|\.jpeg$`, f)
		if err != nil {
			pp.Fatal(err)
		}
		if matched {
			ext := path.Ext(f)
			if ext == ".jpeg" {
				ext = ".jpg"
			}
			ext = strings.Replace(ext, ".", "", 1)
			pp.Printf("Adding: %s\n", f)
			pdfFile.AddPage()
			pdfFile.ImageOptions(f, 0, 0, 0, 0, false, gofpdf.ImageOptions{ImageType: ext}, 0, "")
		}
	}

	err = pdfFile.OutputFileAndClose(filepath.Join(out, "output.pdf"))
	if err != nil {
		pp.Fatal(err)
	}
	pp.Printf("Generated output to: %s\n", out)
}
