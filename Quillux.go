package main

import (
	"math"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/nsf/termbox-go"
)

type Pos struct {
	x uint16
	y uint16
}

type Editor struct {
	lines  []string
	cursor Pos
	size   Pos
}

func newEditor(width, height uint16) *Editor {
	return &Editor{
		lines: []string{""},
		cursor: Pos{
			x: 0,
			y: 0,
		},
		size: Pos{
			x: width,
			y: height,
		},
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		termbox.Interrupt()
	}()

	w, h := termbox.Size()
	editor := newEditor(uint16(w), uint16(h))

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGWINCH)
		for {
			<-ch
			w, h := termbox.Size()
			editor.size.x = uint16(w)
			editor.size.y = uint16(h)
		}
	}()

	for {
		editor.render()
		editor.handleEvent(termbox.PollEvent())
	}
}

func (e *Editor) handleEvent(event termbox.Event) {
	if event.Ch != 0 {
		e.lines[e.cursor.y] += string(event.Ch)
		e.cursor.x++
	} else {
		switch event.Key {
		case termbox.KeyCtrlC, termbox.KeyCtrlD, termbox.KeyCtrlZ:
			termbox.Close()
			os.Exit(0)
		case termbox.KeyEnter:
			e.lines = append(e.lines, "")
			e.cursor.y++
			e.cursor.x = 0
		case termbox.KeyBackspace, termbox.KeyBackspace2:
			if e.cursor.x > 0 {
				e.lines[e.cursor.y] = e.lines[e.cursor.y][:e.cursor.x-1] + e.lines[e.cursor.y][e.cursor.x:]
				e.cursor.x--
			} else if len(e.lines) > 1 {
				e.lines = e.lines[:e.cursor.y]
				e.cursor.y--
				e.cursor.x = uint16(len(e.lines[e.cursor.y]))
			}
		case termbox.KeySpace:
			e.lines[e.cursor.y] += " "
			e.cursor.x++
		}
	}
}

func (e *Editor) render() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	here := Pos{
		x: 0,
		y: 0,
	}
	cursor := Pos{
		x: 0,
		y: 0,
	}
	offset := Pos{
		x: 0,
		y: 0,
	}

	lineNumberLength := uint16(math.Log10(float64(e.cursor.y+1)) + 1)

	const bottomMargin uint16 = 1
	if e.cursor.y+bottomMargin <= e.size.y {
		offset.y = 0
	} else {
		offset.y = e.cursor.y + bottomMargin - e.size.y
	}
	for index, line := range e.lines {
		if index < int(offset.y) {
			continue
		}
		for i, char := range strconv.Itoa(index + 1) {
			termbox.SetChar(i, int(here.y), char)
		}
		termbox.SetChar(int(lineNumberLength), int(here.y), '|')
		cursor = Pos{
			x: lineNumberLength + 1,
			y: here.y,
		}
		for _, char := range line {
			if here.x+lineNumberLength+2 >= e.size.x {
				here.x = 0
				here.y++
				termbox.SetChar(int(lineNumberLength), int(here.y), '|')
			}
			if here.y >= e.size.y {
				break
			}
			cursor = Pos{
				x: here.x + lineNumberLength + 2,
				y: here.y,
			}
			termbox.SetChar(int(cursor.x), int(cursor.y), char)
			here.x++
		}
		if here.x+lineNumberLength+2 >= e.size.x {
			here.x = 0
			here.y++
		}
		if here.y >= e.size.y {
			break
		}
		here.y++
		here.x = 0
	}
	termbox.SetCursor(int(cursor.x+1), int(cursor.y))
	termbox.Flush()
}
