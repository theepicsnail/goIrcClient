package main
import (
    gc "code.google.com/p/goncurses"
)
type ChatArea struct {
    display gc.Window
}
func NewChatArea(win gc.Window) *ChatArea {
    chat := new(ChatArea)
    chat.display = win
    return chat
}

func (c *ChatArea) renderWindow(win *Window) {
    rows,_ := c.display.Maxyx()
    startIdx := len(win.history) - rows
    if startIdx < 0 {
        startIdx = 0
    }

    c.display.Erase()
    for rnum, line := range(win.history[startIdx:]) {
        if rnum == rows {
            break
        }
        c.display.MovePrint(rnum, 0, line)
    }
    c.display.Refresh()
}

