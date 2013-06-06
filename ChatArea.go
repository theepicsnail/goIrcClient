package main
import (
    gc "code.google.com/p/goncurses"
    "fmt"
    "time"
)
type ChatArea struct {
    display gc.Window
    chatChan chan string   
    chatHistory []string  
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
    rows,_ := chat.display.Maxyx()
    chat.chatHistory = make([]string,rows)
    return chat
}

func (c *ChatArea) appendLine(line string) {
    line = fmt.Sprintf("%v| %v", time.Now().Format("15:04:05"), line)
    rows,_ := c.display.Maxyx()
    
    c.display.Erase()
    for rnum, nextLine := range(c.chatHistory) {
        c.display.MovePrint(rows - 1 - rnum, 0, line)
        c.chatHistory[rnum]= line
        line = nextLine
    }
    c.display.Refresh()
}

func (c *ChatArea) GetChatChan() chan<- string {
    return c.chatChan
}
