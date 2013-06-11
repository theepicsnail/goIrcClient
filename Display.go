package main
import (
    gc "code.google.com/p/goncurses"
)
const (
    ChatAreaColor = iota
    InputAreaColor 
)

type Display struct {
    chatArea *ChatArea
    inputArea *InputField
    mainScreen gc.Window
    windowList *WindowList
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

    windowListWidth := 10

    rows, cols := disp.mainScreen.Maxyx()
    disp.windowList = NewWindowList(disp.mainScreen.Derived(rows-1, windowListWidth, 0, 0))

    disp.chatArea = NewChatArea(disp.mainScreen.Derived(rows-1, cols - windowListWidth, 0, windowListWidth))
 
    disp.inputArea = NewInputField(disp.mainScreen.Derived(1, cols, rows-1, 0))
    
    return disp
}
func (d *Display) exit() {
    gc.End()
}

func (d *Display) UserInputChannel() <-chan string {
    return d.inputArea.GetLineChan()
}
