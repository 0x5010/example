package main

import (
	"fmt"
	"strings"

	termbox "github.com/nsf/termbox-go"
)

var (
	current string
	curev   termbox.Event
)

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func mouseButtonStr(k termbox.Key) string {
	switch k {
	case termbox.MouseLeft:
		return "MouseLeft"
	case termbox.MouseMiddle:
		return "MouseMiddle"
	case termbox.MouseRight:
		return "MouseRight"
	case termbox.MouseRelease:
		return "MouseRelease"
	case termbox.MouseWheelUp:
		return "MouseWheelUp"
	case termbox.MouseWheelDown:
		return "MouseWheelDown"
	}
	return "Key"
}

func modStr(m termbox.Modifier) string {
	var out []string
	if m&termbox.ModAlt != 0 {
		out = append(out, "ModAlt")
	}
	if m&termbox.ModMotion != 0 {
		out = append(out, "ModMotion")
	}
	return strings.Join(out, " | ")
}

func redrawAll() {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)
	tbprint(0, 0, termbox.ColorMagenta, coldef, "Press 'q' to quit")
	tbprint(0, 1, coldef, coldef, current)
	switch curev.Type {
	case termbox.EventKey:
		tbprint(0, 2, coldef, coldef,
			fmt.Sprintf("EventKey: k: %d, c: %c, mod: %s", curev.Key, curev.Ch, modStr(curev.Mod)))
	case termbox.EventMouse:
		tbprint(0, 2, coldef, coldef,
			fmt.Sprintf("EventMouse: x: %d, y: %d, b: %s, mod: %s",
				curev.MouseX, curev.MouseY, mouseButtonStr(curev.Key), modStr(curev.Mod)))
	case termbox.EventNone:
		tbprint(0, 2, coldef, coldef, "EventNone")
	}
	tbprint(0, 3, coldef, coldef, fmt.Sprintf("%d", curev.N))
	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputAlt | termbox.InputMouse)
	redrawAll()

	data := make([]byte, 0, 64)
L:
	for {
		if cap(data)-len(data) < 32 {
			newdata := make([]byte, len(data), len(data)+32)
			copy(newdata, data)
			data = newdata
		}
		beg := len(data)
		d := data[beg : beg+32]
		switch ev := termbox.PollRawEvent(d); ev.Type {
		case termbox.EventRaw:
			data = data[:beg+ev.N]
			current = fmt.Sprintf("%q", data)
			if current == `"q"` {
				break L
			}

			for {
				ev := termbox.ParseEvent(data)
				if ev.N == 0 {
					break
				}
				curev = ev
				copy(data, data[curev.N:])
				data = data[:len(data)-curev.N]
			}
		case termbox.EventError:
			panic(ev.Err)
		}
		redrawAll()
	}
}
