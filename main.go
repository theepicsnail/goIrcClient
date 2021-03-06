package main
import (
    "fmt"
    "strings"
    "strconv"
)
func main() {
    var client *IrcClient
    eventChan := make(chan *WindowEvent,100)
    wm := NewWindowManager(eventChan)
    d := NewDisplay(eventChan)
    defer d.exit()
    window := wm.GetWindowByName("Client")
    chatChan := window.GetLineChan()

    messageHandler := func(ircMessages chan *IrcMessage) {
        for msg := range(ircMessages) {
            switch msg.command {
                case "JOIN":
//               Msg: &{"testBot" "JOIN" [] "#test" ":testBot!testBot@airc-BD88CA3C JOIN :#test"}
                    chatChan = wm.GetWindowByName(msg.trailing).GetLineChan()
                case "PRIVMSG":
                    target := msg.source
                    if msg.params[0] != "testBot" {
                        target = msg.params[0]
                    }
                    
                    chatChan = wm.GetWindowByName(target).GetLineChan()
                    chatChan <-fmt.Sprintf("%s | %s", msg.source, msg.trailing)
//               Msg: &{"snail" "PRIVMSG" ["#test"] "test" ":snail!snail@airc-BD88CA3C PRIVMSG #test :test"}
//               Msg: &{"snail" "PRIVMSG" ["testBot"] "Privmsg" ":snail!snail@airc-BD88CA3C PRIVMSG testBot :Privmsg"}
                case "375", "372", "376": //MOTD Start, message, end
                    //wm.GetWindowByName("Client").GetLineChan() <- msg.trailing
                case "PING":
                    continue
                case "NOTICE":
                    wm.GetWindowByName("Notices").GetLineChan() <- fmt.Sprintf("%s | %s", msg.source, msg.trailing) 
                default:
                chatChan <- fmt.Sprintf("TODO:%q", msg)
            }
        }
    }

    
    chatChan <- "Snail's go IRC client!"
    chatChan <- "While I think I have the kinks worked out, you might need to 'reset' after quitting."
    chatChan <- "Use /quit to quit"
    chatChan <- "Use /0 /1 /2... to switch to that window" 
    userInputChan := d.inputArea.GetLineChan()
    go func(){d.inputArea.lineChan <- "/CONNECT localhost:6667"}()
    for msg := range(userInputChan) {
        target := d.windowList.SelectedWindowName()
        if len(msg) >0 && msg[0] == '/' {
            parts := strings.SplitN(msg[1:], " ", 2)
            parts[0] = strings.ToUpper(parts[0]) // Capitalize the command portion
            if len(parts) == 1 {
                if idx, err := strconv.Atoi(parts[0]); err == nil {
                    if win := wm.SelectWindowById(idx); win != nil {
                        chatChan = win.GetLineChan()
                    } else {
                        d.Alert()
                    }
                } else if parts[0] == "QUIT" {
                    if client != nil {
                        client.HandleCommand("QUIT", "Quit message here.")
                    }
                    return 
                }
            } else {
                switch parts[0] {
                    case "CONNECT":
                    if client != nil {
                        client.HandleCommand("QUIT", "Quit message here...")
                    }
                    parts = append(parts, "") //default password.
                    client = NewIrcClient(parts[1], parts[2])
                    client.HandleCommand("USER testBot 0 *", "testBot")
                    client.HandleCommand("NICK", "testBot")
                    go messageHandler(client.output)

                    case "QUIT":
                    if client != nil {
                        client.HandleCommand("QUIT", parts[1])
                    }
                    return

                    case "ME":
                        client.HandleCommand(fmt.Sprintf("PRIVMSG %s", target ), fmt.Sprintf("\001ACTION %s\001", parts[1]))
                    default:
                    client.HandleCommand(parts[0], parts[1])
                } 
            }
        } else {
            client.HandleCommand(fmt.Sprintf("PRIVMSG %s", target), msg)
            chatChan <- msg
        }
    }


}
