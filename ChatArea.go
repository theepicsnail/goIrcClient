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
    
    c.display.Erase()
    for rnum, line := range(win.history) {
        if rnum == rows {
            break
        }
        c.display.MovePrint(rnum, 0, line)
    }
    c.display.Refresh()
}

