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

    chatChan <- "Snail's go IRC client!"
    chatChan <- "While I think I have the kinks worked out, you might need to 'reset' after quitting."
    chatChan <- "Use /quit to quit"
    chatChan <- "Use /0 /1 /2... to switch to that window" 
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
