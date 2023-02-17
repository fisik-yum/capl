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
	"capl/charplayer"
	"flag"
)

var file string

func init() {
	flag.StringVar(&file, "f", "", ".CHNM file to read")

}

func main() {
	echan := make(chan struct{})
	//debug.SetGCPercent()
	//debug.SetMemoryLimit(100000000000)
	flag.Parse()
	scr := charplayer.NewPlayer(file)
	scr.Init()
	//scr.Close()
	go scr.Play(echan)
	//scr.DrawFrame(900)
	//fmt.Println(scr.Frame.Data)
	//scr.Close()
	<-echan
	scr.Close()

}
