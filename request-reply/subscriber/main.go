package main

import (
    "github.com/nats-io/nats.go"
    "log"
    "os/exec"
	"sync"
    "fmt"
    "encoding/json"
    "strconv"

	"github.com/panjf2000/ants/v2"
)

type Message struct {
	InterfaceName   string
	PLR             float64
}

type Task struct {
	msg     Message
	wg      *sync.WaitGroup
}

func (t *Task) configureInterfaces() {
    s := strconv.FormatFloat(t.msg.PLR*100, 'f', -1, 64)
	cmd := exec.Command("sudo", "tc", "qdisc", "change", "dev", t.msg.InterfaceName, "root", "netem", "loss", s)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error configuring interface: %v\n", err)
	}
	t.wg.Done()
}

func main() {
    nc, err := nats.Connect("nats://localhost:4222")
    if err != nil {
        log.Fatal(err)
    }
    defer nc.Close()
    defer ants.Release()
	var wg sync.WaitGroup
	taskFunc := func(data interface{}) {
		task := data.(*Task)
		task.configureInterfaces()
	}
    p, _ := ants.NewPoolWithFunc(4, taskFunc)
    defer p.Release()
    nc.Subscribe("foo", func(msg *nats.Msg) {
        //log.Println("Request receive:", string(msg.Data))
        var resp Message
        err = json.Unmarshal(msg.Data, &resp)
	    if err != nil {
		    log.Fatal(err)
	    }
        wg.Add(1)
		task := &Task{
			msg:    resp,
			wg:      &wg,
		}
		p.Invoke(task)
        wg.Wait()
        msg.Respond([]byte("Finished"))
    })

    select {}
}