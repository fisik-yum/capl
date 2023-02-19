package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/qeesung/image2ascii/convert"
)

func main() {

	files, err := ioutil.ReadDir("frames/")
	if err != nil {
		log.Fatal(err)
	}
	fnames := make([]string, 1)

	sort.Sort(ByNumericalFilename(files))
	//sort.Strings(fnames)
	for _, x := range files {
		fnames = append(fnames, x.Name())
	}

	srcf, err := os.Create("dest.chnm")
	defer srcf.Close()
	fnames = fnames[1:]

	width, height, fps := 72, 54, 30

	convertOptions := convert.DefaultOptions
	convertOptions.FixedWidth = width
	convertOptions.FixedHeight = height
	convertOptions.Colored = false
	converter := convert.NewImageConverter()

	fmt.Fprintf(srcf, "%d %d %d", width, height, fps)
	for a := range fnames {
		fmt.Println("\033[H\033[J")
		fmt.Println(fnames[a])
		afn := filepath.Join("frames", fnames[a])
		fmt.Println(afn)
		asciiArt := converter.ImageFile2ASCIIString(afn, &convertOptions)
		input := strings.ReplaceAll(asciiArt, "\n", "")
		fmt.Fprintf(srcf, "%s\n", input)
	}
	if err != nil {
		fmt.Println(err)
	}
}

type ByNumericalFilename []os.FileInfo

func (nf ByNumericalFilename) Len() int      { return len(nf) }
func (nf ByNumericalFilename) Swap(i, j int) { nf[i], nf[j] = nf[j], nf[i] }
func (nf ByNumericalFilename) Less(i, j int) bool {

	// Use path names
	pathA := nf[i].Name()
	pathB := nf[j].Name()

	// Grab integer value of each filename by parsing the string and slicing off
	// the extension
	a, err1 := strconv.ParseInt(pathA[0:strings.LastIndex(pathA, ".")], 10, 64)
	b, err2 := strconv.ParseInt(pathB[0:strings.LastIndex(pathB, ".")], 10, 64)

	// If any were not numbers sort lexographically
	if err1 != nil || err2 != nil {
		return pathA < pathB
	}

	// Which integer is smaller?
	return a < b
}
