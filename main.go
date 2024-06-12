package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

func drawCircle(s tcell.Screen, centerX, centerY, radius int, circleStyle, borderStyle tcell.Style) {
	for y := -radius - 1; y <= radius+1; y++ {
		for x := -radius - 1; x <= radius+1; x++ {
			distance := x*x + y*y
			if distance <= radius*radius {
				s.SetCell(centerX+x, centerY+y, circleStyle, ' ')
			} else if distance <= (radius+1)*(radius+1) {
				s.SetCell(centerX+x, centerY+y, borderStyle, ' ')
			}
		}
	}
	s.Show()
}

func drawRectangle(s tcell.Screen, topLeftX, topLeftY, width, height int, rectStyle, borderStyle tcell.Style) {
	for y := 0; y <= height; y++ {
		for x := 0; x <= width; x++ {
			if x == 0 || x == width || y == 0 || y == height {
				s.SetCell(topLeftX+x, topLeftY+y, borderStyle, ' ')
			} else {
				s.SetCell(topLeftX+x, topLeftY+y, rectStyle, ' ')
			}
		}
	}
	s.Show()
}

func main() {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorBlack).
		Background(tcell.ColorWhite))
	s.Clear()

	quit := make(chan struct{})
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape, tcell.KeyEnter:
					close(quit)
					return
				case tcell.KeyCtrlL:
					s.Sync()
				}
			case *tcell.EventResize:
				s.Sync()
			}
		}
	}()

	w, h := s.Size()
	centerX := w / 2
	centerY := h / 4 // Top center
	radius := 6      // Small circle radius
	rectWidth := 8
	rectWidthMin := 3
	rectHeight := 4
	rectTopLeftX := centerX - rectWidth*3
	firstLineY := centerY + radius + 2
	secondLineY := firstLineY + rectHeight + 1
	thirdLineY := secondLineY + rectHeight + 1
	fourthLineY := thirdLineY + rectHeight + 1
	orange := tcell.NewRGBColor(255, 165, 0)
	red := tcell.NewRGBColor(255, 0, 0)
	grey := tcell.ColorGrey
	blink := true

loop:
	for {
		select {
		case <-quit:
			break loop
		case <-time.After(time.Second * 1): // Blink interval
		}

		if blink {
			drawCircle(s, centerX, centerY, radius, tcell.StyleDefault.Background(orange), tcell.StyleDefault.Background(grey))
		} else {
			drawCircle(s, centerX, centerY, radius, tcell.StyleDefault.Background(tcell.ColorWhite), tcell.StyleDefault.Background(grey))
		}

		// Draw the hours, 4 boxes, each box represents 5 hour slot
		n5HourSwitched := time.Now().Hour() / 5
		//fmt.Printf("hour: %d n5HourSwitched: %d\n", time.Now().Hour(), n5HourSwitched)
		for i := 1; i <= n5HourSwitched; i++ {
			drawRectangle(s, rectTopLeftX+rectWidth*i, firstLineY, rectWidth, rectHeight, tcell.StyleDefault.Background(orange), tcell.StyleDefault.Background(grey))
		}
		for i := n5HourSwitched + 1; i < 5; i++ {
			drawRectangle(s, rectTopLeftX+rectWidth*i, firstLineY, rectWidth, rectHeight, tcell.StyleDefault.Background(tcell.ColorWhite), tcell.StyleDefault.Background(grey))
		}

		// Draw the hours, 4 boxes, each box represents 1 hour slot
		n4HourSwitched := time.Now().Hour() % 5
		//fmt.Printf("hour: %d n4HourSwitched: %d\n", time.Now().Hour(), n4HourSwitched)
		for i := 1; i <= n4HourSwitched; i++ {
			drawRectangle(s, rectTopLeftX+rectWidth*i, secondLineY, rectWidth, rectHeight, tcell.StyleDefault.Background(orange), tcell.StyleDefault.Background(grey))
		}
		for i := n4HourSwitched + 1; i < 5; i++ {
			drawRectangle(s, rectTopLeftX+rectWidth*i, secondLineY, rectWidth, rectHeight, tcell.StyleDefault.Background(tcell.ColorWhite), tcell.StyleDefault.Background(grey))
		}

		// Draw the minutes, 12 boxes, each box represents 5 min slot
		n5minSwitched := time.Now().Minute() / 5
		//fmt.Printf("minute: %d n5minSwitched: %d\n", time.Now().Minute(), n5minSwitched)
		for i := 1; i <= n5minSwitched; i++ {
			if i%3 == 0 {
				drawRectangle(s, rectTopLeftX+4+rectWidthMin*i, thirdLineY, rectWidthMin, rectHeight, tcell.StyleDefault.Background(red), tcell.StyleDefault.Background(grey))
			} else {
				drawRectangle(s, rectTopLeftX+4+rectWidthMin*i, thirdLineY, rectWidthMin, rectHeight, tcell.StyleDefault.Background(orange), tcell.StyleDefault.Background(grey))
			}
		}
		for i := n5minSwitched + 1; i < 12; i++ {
			drawRectangle(s, rectTopLeftX+4+rectWidthMin*i, thirdLineY, rectWidthMin, rectHeight, tcell.StyleDefault.Background(tcell.ColorWhite), tcell.StyleDefault.Background(grey))
		}

		// Draw the hours, 4 boxes, each box represents 1 minute slot
		n4MinuteSwitched := time.Now().Minute() % 5
		//fmt.Printf("minute: %d n4MinuteSwitched: %d\n", time.Now().Minute(), n4MinuteSwitched)
		for i := 1; i <= n4MinuteSwitched; i++ {
			drawRectangle(s, rectTopLeftX+rectWidth*i, fourthLineY, rectWidth, rectHeight, tcell.StyleDefault.Background(orange), tcell.StyleDefault.Background(grey))
		}
		for i := n4MinuteSwitched + 1; i < 5; i++ {
			drawRectangle(s, rectTopLeftX+rectWidth*i, fourthLineY, rectWidth, rectHeight, tcell.StyleDefault.Background(tcell.ColorWhite), tcell.StyleDefault.Background(grey))
		}

		blink = !blink
	}

	s.Fini()
}
