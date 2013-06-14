package main
import (
    "fmt"
    "time"
    "strings"
    "strconv"
)
func main() {
    eventChan := make(chan *WindowEvent,100)
    wm := NewWindowManager(eventChan)
    d := NewDisplay(eventChan)
    defer d.exit()
    window1 := wm.GetWindowByName("Test")
    go func() {
        for {
            time.Sleep(5e9)
            window1.GetLineChan() <- "A"
        }
    }()
    window := wm.GetWindowByName("Main")
    chatChan := window.GetLineChan()

    chatChan <- "Welcome to snails shitty chat thing."
    chatChan <- "Press esc to quit, it may or may not break stuff. "
    chatChan <- "If it does, do a 'reset' to fix it."
    chatChan <- "Use /quit to exit."
    chatChan <- "" 
    userInputChan := d.inputArea.GetLineChan()
    for msg := range(userInputChan) {
        if len(msg) >0 && msg[0] == '/' {
            parts := strings.Split(msg[1:], " ")
            if len(parts) == 1 {
                if idx, err := strconv.Atoi(parts[0]); err == nil {
                    if win := wm.SelectWindowById(idx); win != nil {
                        chatChan = win.GetLineChan()
                    }
                } else if parts[0] == "quit" {
                    return 
                }
            } else {
                chatChan <- fmt.Sprintf("Command [%s] %s", parts[0], parts[1:])
            }
        } else {
            chatChan <- msg
        }
    }


}
