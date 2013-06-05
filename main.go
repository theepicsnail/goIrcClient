package main
import (
    gc "code.google.com/p/goncurses"
    "fmt"
    "time"
)
const (
    ChatAreaColor = iota
    InputAreaColor 
)

type Display struct {
    chatArea gc.Window
    inputArea *InputField
    mainScreen gc.Window
    ChatChan chan<- string //Input only channel to add strings to the display
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
    disp.chatArea = disp.mainScreen.Derived(rows-1, cols, 0, 0)
    disp.chatArea.ScrollOk(true)
 
    disp.inputArea = NewInputField(disp.mainScreen.Derived(1, cols, rows-1, 0))
    
    ch := make(chan string)
    disp.ChatChan = ch
    go func() {
        for line := range(ch) {
            disp.appendLine(line)
        }
    }()
    return disp
}
func (d *Display) appendLine(line string) {
    line = fmt.Sprintf("%v| %v", time.Now().Format("15:04:05"), line)
    rows,_ := d.chatArea.Maxyx()
    d.chatArea.Scroll(1)
    d.chatArea.MovePrint(rows - 1, 0, line)
    d.chatArea.Refresh()

}
func (d *Display) exit() {
    gc.End()
}

func (d *Display) MainLoop() {
    defer d.exit()
    defer close(d.ChatChan)


    d.ChatChan <- "Welcome to snails shitty chat thing."
    d.ChatChan <- "Press esc to quit, it may or may not break stuff. "
    d.ChatChan <- "If it does, do a 'reset' to fix it."
    d.ChatChan <- "Use /quit to exit."
    d.ChatChan <- "" 
    userInputChan := d.inputArea.GetLineChan()
    for msg := range(userInputChan) {
        if msg == "/quit" {
            return
        }
        if len(msg) >0 && msg[0] == '/' {
            d.ChatChan <- fmt.Sprintf("Command: %s", msg[1:])
        } else {
            d.ChatChan <- msg
        }
    }
}


func main() {
    d := NewDisplay()
    d.MainLoop()
}
