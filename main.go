package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jung-kurt/gofpdf"

	"github.com/k0kubun/pp"
)

func main() {
	var err error
	count := 0
	dir := flag.String("dir", "", "path to the directory contains your images")
	out := flag.String("out", "", "path to output pdf file")

	flag.Parse()

	// start a new pdf file
	pdfFile := gofpdf.New("P", "mm", "A4", "")

	// if dir is not provide we use current directory
	if *dir == "" {
		*dir, err = os.Getwd()
		if err != nil {
			pp.Fatal(err)
		}
	}

	if *out == "" {
		*out, err = os.Getwd()
		if err != nil {
			pp.Fatal(err)
		}
	}

	filesOfDir, err := ioutil.ReadDir(*dir)
	if err != nil {
		pp.Fatal(err)
	}

	pp.Printf("Processing file in directory %s\n", *dir)
	for _, f := range filesOfDir {
		filename := f.Name()
		matched, err := regexp.MatchString(`\.png$|\.jpg$|\.svg$|\.jpeg$`, filename)
		if err != nil {
			pp.Fatal(err)
		}
		if matched {
			count++
			ext := path.Ext(filename)
			if ext == ".jpeg" {
				ext = ".jpg"
			}
			ext = strings.Replace(ext, ".", "", 1)
			pp.Printf("Adding: %s\n", filename)
			pdfFile.AddPage()
			pdfFile.ImageOptions(filepath.Join(*dir, filename), 0, 0, 0, 0, false, gofpdf.ImageOptions{ImageType: ext}, 0, "")
		}
	}

	if count == 0 {
		pp.Printf("There is no images in this directory\n")
		os.Exit(0)
	}

	err = pdfFile.OutputFileAndClose(filepath.Join(*out, "output.pdf"))
	if err != nil {
		pp.Fatal(err)
	}
	pp.Printf("Generated output to: %s\n", *out)
}
