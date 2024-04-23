package main

import (
    "github.com/nats-io/nats.go"
    "log"
    "time"
    "encoding/json"
    "sync"
    "fmt"
    "math/rand"
)

type Message struct {
	InterfaceName   string
	PLR             float64
}

func main() {
    var wg sync.WaitGroup
    nc, err := nats.Connect("nats://localhost:4222")
    if err != nil {
        log.Fatal(err)
    }
    defer nc.Close()

    for i := 0; i < 400; i++ {
        wg.Add(1)
        rand.Seed(time.Now().UnixNano())
		randomNumber := rand.Intn(20)
        name := fmt.Sprintf("veth%d", randomNumber)
        msg := Message{
            InterfaceName: name,
            PLR:    0.001,
        }
        reqData, err := json.Marshal(msg)
	    if err != nil {
		    log.Fatal(err)
	    }
        _, err = nc.Request("foo", reqData, 10*time.Second)
        if err != nil {
            log.Fatal(err)
        }
        wg.Done()
    }
    wg.Wait()
}