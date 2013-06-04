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
    inputArea gc.Window
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
 
    disp.inputArea = disp.mainScreen.Derived(1, cols, rows-1, 0)
    disp.inputArea.Keypad(true)
    
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

    userInputChan := make(chan string)
    go func() {
        for msg := range(userInputChan) {
            if len(msg) >0 && msg[0] == '/' {
                d.ChatChan <- fmt.Sprintf("Command: %s", msg[1:])
            } else {
                d.ChatChan <- msg
            }
        }
    }()
    defer close(userInputChan)

    d.ChatChan <- "Welcome to snails shitty chat thing."
    d.ChatChan <- "Press esc to quit, it may or may not break stuff. "
    d.ChatChan <- "If it does, do a 'reset' to fix it."
 
    buffer := ""
    for {
        d.chatArea.Refresh()
        key := d.inputArea.GetChar()
        switch key {
            case gc.Key(27):
                return
            case gc.KEY_RETURN:
                userInputChan <- buffer
                buffer = ""
            case gc.Key(127)://backspace
                l := len(buffer)
                if l > 0 {
                    buffer = buffer[:l-1]
                }
            default:
                buffer = fmt.Sprintf("%s%c", buffer, key)
        }
        d.inputArea.Clear()
        d.inputArea.MovePrint(0, 0, buffer)
    } 
}


func main() {
    d := NewDisplay()
    d.MainLoop()
}
