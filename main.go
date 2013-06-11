package main
import (
    "fmt"
    "strings"
)
func main() {
    d := NewDisplay()
    defer d.exit()

    chatChan := d.chatArea.GetChatChan()

    chatChan <- "Welcome to snails shitty chat thing."
    chatChan <- "Press esc to quit, it may or may not break stuff. "
    chatChan <- "If it does, do a 'reset' to fix it."
    chatChan <- "Use /quit to exit."
    chatChan <- "" 

    userInputChan := d.inputArea.GetLineChan()
    for msg := range(userInputChan) {
        if len(msg) >0 && msg[0] == '/' {
            parts := strings.Split(msg[1:], " ")
            
            chatChan <- fmt.Sprintf("Command [%s] %s", parts[0], parts[1:])
            if parts[0] == "quit" {
                return
            }
        } else {
            chatChan <- msg
        }
    }


}
