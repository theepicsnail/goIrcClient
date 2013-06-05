package main
import (
    gc "code.google.com/p/goncurses"
    "fmt"
)
const (
    ChatAreaColor = iota
    InputAreaColor 
)

type Display struct {
    chatArea *ChatArea
    inputArea *InputField
    mainScreen gc.Window
}

func NewDisplay() *Display {
    disp := new(Display)
    var err error
    disp.mainScreen, err = gc.Init()
    if err != nil {
        panic(err)
    }

    gc.Echo(false)
    gc.CBreak(true)
    gc.Raw(true)

    rows, cols := disp.mainScreen.Maxyx()
    disp.chatArea = NewChatArea(disp.mainScreen.Derived(rows-1, cols, 0, 0))
 
    disp.inputArea = NewInputField(disp.mainScreen.Derived(1, cols, rows-1, 0))
    
    return disp
}
func (d *Display) exit() {
    gc.End()
}

func (d *Display) MainLoop() {
    defer d.exit()

    chatChan := d.chatArea.GetChatChan()

    chatChan <- "Welcome to snails shitty chat thing."
    chatChan <- "Press esc to quit, it may or may not break stuff. "
    chatChan <- "If it does, do a 'reset' to fix it."
    chatChan <- "Use /quit to exit."
    chatChan <- "" 

    userInputChan := d.inputArea.GetLineChan()
    for msg := range(userInputChan) {
        if msg == "/quit" {
            return
        }
        if len(msg) >0 && msg[0] == '/' {
            chatChan <- fmt.Sprintf("Command: %s", msg[1:])
        } else {
            chatChan <- msg
        }
    }
}


func main() {
    d := NewDisplay()
    d.MainLoop()
}
