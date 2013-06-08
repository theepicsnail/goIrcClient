package main
import (
    "fmt"
    "bufio"
    "strings"
    "net"
    "regexp"
)
type IrcClient struct {
    conn net.Conn
    input chan string
    output chan *IrcMessage
}

//Example private message
//:snail!snail@airc-BD88CA3C PRIVMSG testbot :Test :)
// 0 - Full line            (:snail!snail@airc-BD88CA3C PRIVMSG testbot :test)
// 1 - Full prefix          (:snail!snail@airc-BD88CA3C)
//   2 - Prefix without :   (snail!snail@airc-BD88CA3C)
// 3 - command              (PRIVMSG)
// 4 - args                 ( testbot) NOTE, leading space IS included (unless there are 0 args)
// 5 - full tail            ( :Test :))
//   6 - Tail without :     (Test :))
//
//Example ping
//0 (PING :og.udderweb.com
//1 ()
//2 ()
//3 (PING)
//4 ()
//5 (:og.udderweb.com)
//6 (og.udderweb.com)
var ircLineRe = regexp.MustCompile("^(:([^ !]+)[^ ]* )?([^ ]*)(.*?)( :(.*))?$")

type IrcMessage struct {
    source string
    command string
    params []string
    trailing string
    raw string
}
//The above two examples return:
//{source: "snail", command: "PRIVMSG", params:["testbot"], trailing: "Test :)"}
//{source: "", command: "PING", params:[], trailing: "og.udderweb.com"}
//:og.udderweb.com 004 testBot og.udderweb.com Unreal3.2.9 iowghraAsORTVSxNCWqBzvdHtGp lvhopsmntikrRcaqOALQbSeIKVfMCuzNTGjZ
func parseMessage(line string) *IrcMessage {
    msg := new(IrcMessage)
    
    match := ircLineRe.FindAllStringSubmatch(line, -1)
    if match != nil {
        msg.source = match[0][2]
        msg.command= match[0][3]
        msg.params = strings.Split(match[0][4]," ")[1:]
        msg.trailing=match[0][6]
        msg.raw = line
        return msg
    } 
    fmt.Println("Didn't match:")
    fmt.Println(line)
    return msg
}
func NewIrcClient(hostport, password string) *IrcClient {
    con, err := net.Dial("tcp", hostport)
    if err != nil {
        panic(err)
    }
    client := new(IrcClient)
    client.conn = con
    
    if password != "" {
        fmt.Fprint(client.conn, "PASS %s\r\n", password) 
    }

    client.input = make(chan string, 100)
    go func() {
        for line := range(client.input) {
            fmt.Fprint(client.conn, line + "\r\n")
        }
    }()

    client.output = make(chan *IrcMessage, 100)
    go func() {
        bufRead := bufio.NewReader(client.conn)
        for {
            line, err := bufRead.ReadString('\n')
            if err != nil {
                panic(err)
            }
            line = strings.TrimRight(line,"\r\n")
            client.output <- parseMessage(line)
            if strings.HasPrefix(line, "PING") {
                client.input <- strings.Replace(line, "I", "O", 1)
            }
        }
    }()
    return client
}

func main() {
    name := "testBot"
    c := NewIrcClient("localhost:6667", "")
    c.input <- "USER " + name + " 0 * :" + name
    c.input <- "NICK " + name
    for msg := range(c.output) {
        fmt.Println("────┼───────")
        fmt.Println("raw │ ", msg.raw)
        fmt.Println("src │ ", msg.source)
        fmt.Println("cmd │ ", msg.command)
        fmt.Println("args│ ", msg.params)
        fmt.Println("tail│ ", msg.trailing)
        switch msg.command {
            case "376":
                c.input <- "JOIN #test"
        }
    }
}
