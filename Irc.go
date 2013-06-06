package main
import (
    "fmt"
    "bufio"
    "strings"
    "net"
)

type IrcClient struct {
    conn net.Conn
    input chan string
    output chan string
}

func NewIrcClient(hostport string) *IrcClient {
    con, err := net.Dial("tcp", hostport)
    if err != nil {
        panic(err)
    }
    client := new(IrcClient)
    client.conn = con
     
    client.input = make(chan string, 100)
    go func() {
        for line := range(client.input) {
            fmt.Fprint(client.conn, line + "\r\n")
        }
    }()

    client.output = make(chan string, 100)
    go func() {
        bufRead := bufio.NewReader(client.conn)
        for {
            line, err := bufRead.ReadString('\n')
            if err != nil {
                panic(err)
            }
            line = strings.TrimRight(line,"\r\n")
            if strings.HasPrefix(line, "PING") {
                client.input <- strings.Replace(line, "I", "O", 1)
            }
            client.output <- line
        }
    }()
    return client
}

func main() {
    c := NewIrcClient("localhost:6667")
    c.input <- "USER localhost localhost localhost :realname"
    c.input <- "NICK testBot"
    var from, msgType, target, data string
    for line := range(c.output) {
        fmt.Sscanf(line, ":%s %s %s :%s", &from, &msgType, &target, &data)
        if msgType == "376" {
            c.input <- "JOIN #test"
        }
        fmt.Println("────┼───────")
        fmt.Println("from│ ", from)
        fmt.Println("type│ ", msgType)
        fmt.Println("to  │ ", target)
        fmt.Println("data│ ", data)
    }
}
