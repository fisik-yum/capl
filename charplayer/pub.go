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
package charplayer

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/stoicperlman/fls"
)

func NewPlayer(filename string) Player {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	//defer f.Close()
	//wi, _, fps := dataOf(film)
	p := Player{
		init:      false,
		framechan: make(chan Frame),
		State: struct {
			IsPlaying bool
			Loop      bool
			pausechan chan struct{}
		}{
			pausechan: make(chan struct{}),
		},
		Source: fls.LineFile(f),
	}
	p.Dimension.Width, p.Dimension.Height, p.fps = dataOf(p.Source)
	return p
}

func dataOf(f *fls.File) (int, int, int) { //width height fps
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
	return keys[0], keys[1], keys[2]
}
