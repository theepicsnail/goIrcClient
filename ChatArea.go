package main
import (
    gc "code.google.com/p/goncurses"
    "fmt"
    "time"
)
type ChatArea struct {
    display gc.Window
    chatChan chan string    
}
func NewChatArea(win gc.Window) *ChatArea {
    chat := new(ChatArea)
    chat.display = win
    chat.chatChan = make(chan string)
    go func() {
        for line:=range(chat.chatChan) {
            chat.appendLine(line)
        }
    }()
    win.ScrollOk(true)  
    return chat
}

func (c *ChatArea) appendLine(line string) {
    line = fmt.Sprintf("%v| %v", time.Now().Format("15:04:05"), line)
    rows,_ := c.display.Maxyx()
    c.display.Scroll(1)
    c.display.MovePrint(rows - 1, 0, line)
    c.display.Refresh()
}

func (c *ChatArea) GetChatChan() chan<- string {
    return c.chatChan
}
