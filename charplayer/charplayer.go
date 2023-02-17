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
	"time"

	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
	"github.com/stoicperlman/fls"
)

type Player struct {
	//Mutex  *sync.Mutex
	Source    *fls.File
	init      bool
	Screen    tcell.Screen
	fps       int
	framechan chan Frame
	State     struct {
		IsPlaying bool
		Loop      bool
		pausechan chan struct{}
	}
	Dimension struct {
		Height int
		Width  int
	}
}

type Frame struct {
	data []string
	N    int
}

func (p *Player) Init() {
	var err error
	if !p.init {
		p.Screen, err = tcell.NewScreen()
		if err != nil {
			log.Fatal(err)
		}
		p.Screen.Init()
		p.init = !p.init
	} else {
		os.Exit(1)
	}
}

func (p *Player) Close() {
	if p.init {
		p.Screen.Fini()
		os.Exit(1)
	} else {
		return
	}
}

func (p *Player) GetFrame(n int64) Frame {
	pos, _ := p.Source.SeekLine(n, io.SeekStart)
	//defer p.wg.Done()
	p.Source.Seek(pos, io.SeekStart)
	reader := bufio.NewReader(p.Source.File)
	line, _, _ := reader.ReadLine()
	/*if err != nil && !p.Loop {
		os.Exit(1)
		return errors.New("EOF")
		//cleanup
	}*/
	return Frame{
		data: tochunks(string(line), p.Dimension.Width),
		N:    int(n),
	}
	//p.State.framechan <- struct{}{}
	//fmt.Println(tochunks(string(line), p.Dimension.Width))

}

func (p *Player) DrawFrame(f Frame) {
	x, y := 0, 0
	//p.wg.Wait()
	//p.Screen.Clear()
	for _, i := range f.data {
		for _, c := range i {
			var comb []rune
			w := runewidth.RuneWidth(c)
			if w == 0 {
				comb = []rune{c}
				c = ' '
				w = 1
			}
			p.Screen.SetContent(x, y, c, comb, tcell.StyleDefault)
			//p.screen.Show()
			x += w
		}
		x = 0
		y++
	}
	afc := []rune("Frame " + strconv.Itoa(f.N))
	for x, i := range afc {
		p.Screen.SetContent(x, y, i, nil, tcell.Style.Foreground(tcell.Style.Background(tcell.StyleDefault, tcell.ColorGray), tcell.ColorBlack))
	}
	p.Screen.Show()

	//p.screen.GetContent()
}

func (p *Player) Play(exitchan chan<- struct{}) {
	p.State.IsPlaying = true
	go func() {
		for {
			ev := p.Screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape:
					exitchan <- struct{}{}
				case tcell.KeyEnter:
					if p.State.IsPlaying {
						p.State.IsPlaying = !p.State.IsPlaying
					} else {
						p.State.IsPlaying = !p.State.IsPlaying
						p.State.pausechan <- struct{}{}
					}
				case tcell.KeyTab:
					p.State.Loop = !p.State.Loop
				}
			case *tcell.EventResize:
				p.Screen.Sync()
			}
		}
	}()
	go func() {
		i := 0
		for {
			//time.Sleep(time.Duration(uint(time.Second) / uint(p.fps)))
			//go func() {
			if p.State.IsPlaying {
				p.framechan <- p.GetFrame(int64(i))
				i++
			} else {
				p.framechan <- p.GetFrame(int64(i))
				<-p.State.pausechan
			}
			//}()
		}
	}()
	go func() {
		fc := 0
		for {
			f := <-p.framechan
			time.Sleep(time.Duration(uint(time.Second) / uint(p.fps+1)))
			if f.N < fc {
				continue
			} else {
				fc = f.N
				p.DrawFrame(f)
			}
		}
	}()

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
