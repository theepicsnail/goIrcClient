package main
import (
    gc "code.google.com/p/goncurses"
    "fmt"
)

type InputField struct {
    display gc.Window
    lineChan chan string
    buffer string
}

func NewInputField(window gc.Window) *InputField {
    f := new(InputField)
    f.display = window
    f.buffer = ""
    window.Keypad(true)
    f.lineChan = make(chan string)
    go func() {
        for {
            key := f.display.GetChar()
            f.handleKey(key)
        }
    }()
    return f
}

func (f *InputField) handleKey(key gc.Key) {
    switch key {
        case gc.KEY_RETURN:
            f.lineChan <- f.buffer
            f.buffer = ""
        case gc.Key(127):
            l := len(f.buffer)
            if l > 0 {
                f.buffer = f.buffer[:l-1]
            }
        default:
            f.buffer = fmt.Sprintf("%s%c", f.buffer, key)
    }
    go func(){CURSES.Lock()
    f.display.Clear()
    f.display.MovePrint(0, 0, f.buffer)
    f.display.Refresh()
    CURSES.Unlock()}()
}

func (f *InputField) GetLineChan() <-chan string {
    return f.lineChan
}
