package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/nsf/termbox-go"
)

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

	posX := 0
	posY := 0
	width, height := termbox.Size()
	termbox.SetCursor(posX, posY)
	termbox.Flush()
loop:
	for {
		event := termbox.PollEvent()
		if event.Ch != 0 {
			if posX < width-1 {
				termbox.SetChar(posX, posY, event.Ch)
				posX++
			}
		} else {
			switch event.Key {
			case termbox.KeyCtrlC, termbox.KeyCtrlD, termbox.KeyCtrlZ:
				break loop
			case termbox.KeySpace:
				if posX < width-1 {
					posX++
				}
			case termbox.KeyEnter:
				if posY < height-1 {
					posY++
					posX = 0
				}
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				if posX > 0 {
					termbox.SetChar(posX-1, posY, ' ')
					posX--
				}
			}
		}
		termbox.SetCursor(posX, posY)
		termbox.Flush()
	}
}
