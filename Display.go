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

func NewDisplay(eventChan <-chan *WindowEvent) *Display {
    disp := new(Display)
    var err error
    disp.mainScreen, err = gc.Init()
    if err != nil {
        panic(err)
    }

    gc.Echo(false)
    gc.CBreak(true)
    gc.Raw(true)
    gc.Cursor(0)

    windowListWidth := 10

    rows, cols := disp.mainScreen.Maxyx()
    disp.windowList = NewWindowList(
        disp.mainScreen.Derived(rows-1, windowListWidth, 0, 0))

    disp.chatArea = NewChatArea(
        disp.mainScreen.Derived(
        rows-1, cols - windowListWidth, 0, windowListWidth))
 
    disp.inputArea = NewInputField(
        disp.mainScreen.Derived(1, cols, rows-1, 0))
    
    go disp.windowEventConsumer(eventChan)

    return disp
}

func (disp *Display) windowEventConsumer( eventChan <-chan *WindowEvent) {
    for event := range(eventChan) {
        disp.windowList.GetWindowEventChan() <- event
        if disp.windowList.SelectedWindowName() == event.window.name {
            disp.chatArea.renderWindow(event.window)
        }
    }
}

func (d *Display) exit() {
    gc.End()
}

func (d *Display) UserInputChannel() <-chan string {
    return d.inputArea.GetLineChan()
}
