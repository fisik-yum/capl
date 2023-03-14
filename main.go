/*
capl - A program to view asciimation movies
Copyright (C) 2021  fisik_yum
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/stoicperlman/fls"
)

var file string
var ffile *fls.File
var width int
var fps float64

type Frame struct {
	data []string
	N    int
	err  error
}

func init() {
	flag.StringVar(&file, "f", "", ".CHNM file to read")
	flag.Parse()
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	ffile = fls.LineFile(f)
	width, _, fps = dataOf(ffile)
}

func main() {
	ti := time.NewTicker(time.Duration(uint(time.Second) / uint(fps)))
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fchan := make(chan Frame, 1)
	pchan := make(chan struct{})
	pstate := true
	start := time.Now()
	go func() {
		i := 0
		for {
			select {
			case fchan <- GetFrame(int64(i)):
				i++
			}

		}
	}()
	go func() {
		for {
			_, key, err := keyboard.GetKey()
			if err != nil {
				panic(err)
			}
			if key == keyboard.KeyEsc {
				keyboard.Close()
				ti.Stop()
				os.Exit(1)
			} else if key == keyboard.KeyEnter {
				if pstate == false {
					pstate = !pstate
					pchan <- struct{}{}
				} else {
					pstate = false
				}
			}
		}
	}()

	for {
		if pstate == false {
			<-pchan
		}
		select {
		case a := <-fchan:
			if a.err != nil {
				fmt.Println(time.Since(start))
				keyboard.Close()
                ti.Stop()
				os.Exit(0)
			}
			DrawFrame(a.data)
            fmt.Println(a.N)
		}
		<-ti.C
	}
}

func tochunks(s string, w int) []string {
	var chunks []string = make([]string, 0, (len(s)-1)/w+1)
	currentLen := 0
	currentStart := 0
	for i := range s {
		if currentLen == w {
			chunks = append(chunks, s[currentStart:i])
			currentLen = 0
			currentStart = i
		}
		currentLen++
	}
	chunks = append(chunks, s[currentStart:])
	return chunks
}

func GetFrame(n int64) Frame {
	pos, _ := ffile.SeekLine(n, io.SeekStart)
	ffile.Seek(pos, io.SeekStart)
	reader := bufio.NewReader(ffile.File)
	line, _, err := reader.ReadLine()
	return Frame{
		data: tochunks(string(line), width),
		N:    int(n),
		err:  err,
	}
}

func DrawFrame(f []string) {
	os.Stdout.WriteString("\033[H\033[2J")
	for _, x := range f {
		os.Stdout.WriteString(x + "\n")
	}
}

func dataOf(f *fls.File) (int, int, float64) {
	pos, _ := f.SeekLine(0, io.SeekStart)
	f.Seek(pos, io.SeekStart)
	reader := bufio.NewReader(f)
	line, _, err := reader.ReadLine()
	if err != nil {
		log.Fatal(err)
	}
	tentative := strings.Split(string(line), " ")
	if len(tentative) < 3 {
		log.Fatal("Invalid metadata descriptor")
	}
	keys := [3]int{}

	for a := 0; a < 3; a++ {
		keys[a], err = strconv.Atoi(tentative[a])
		if err != nil {
			log.Fatal("Invalid file descriptor")
		}
	}
	return keys[0], keys[1], float64(keys[2])
}
