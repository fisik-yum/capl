package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var wi int
var he int
var file string

// starwars height 14 (actual 13) width 67
func init() {
	flag.IntVar(&wi, "w", 0, "width of source file")
	flag.IntVar(&he, "h", 0, "height of source file, including runlength compression integer")
	flag.StringVar(&file, "f", "", "sourcefile name, including extension")
	flag.Parse()
}
func main() {
	f, err := os.Open(file)
	check(err)
	defer f.Close()
	s := bufio.NewScanner(f)
	srcf, err := os.Create(fmt.Sprint(file, ".chnm"))
	check(err)
	h := true
	ln := 1
	//h = s.Scan()
	for h {
		a := make([]string, 0)
		for x := 0; x < he; x++ {
			a = append(a, s.Text())
			h = s.Scan()

		}
		a = append(a, s.Text())

		fmt.Printf("Line %d \n", ln)
		a = a[1:] //reduce preceding whitespace
		fmt.Println(len(a))
		//fmt.Println(a)
		d, err := strconv.Atoi(a[0]) //get display time
		a = a[1:]
		check(err)
		fmt.Println((len(a)))
		for l := range a {
			a[l] = extend(a[l], wi)
		}
		for i := 0; i < d; i++ {

			fmt.Fprintf(srcf, "%s\n", strings.Join(a, ""))
		}
		ln++
	}
}

func extend(st string, l int) string {
	if len(st) < l {
		for len(st) < l {
			st = st + " "
		}
		return st
	} else if len(st) > l {
		return st[:l]
	} else {
		return st
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
