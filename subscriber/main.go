package subscriber

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	ents "github.com/sabrs0/L0_WB/entities"
)

var streamName = "EVENTS"

func main() {
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
	}

	// Create an unauthenticated connection to NATS.
	nc, err := nats.Connect(url)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected")
	// Drain is a safe way to to ensure all buffered messages that were published
	// are sent and all buffered messages received on a subscription are processed
	// being closing the connection.
	defer nc.Drain()
	js, err := nc.JetStream()
	if err != nil {
		panic(err)
	}
	fmt.Println("Jet Stream Created")
	// We will declare the initial stream configuration by specifying
	// the name and subjects. Stream names are commonly uppercased to
	// visually differentiate them from subjects, but this is not required.
	// A stream can bind one or more subjects which almost always include
	// wildcards. In addition, no two streams can have overlapping subjects
	// otherwise the primary messages would be persisted twice. There
	// are option to replicate messages in various ways, but that will
	// be explained in later examples.
	/*cfg := nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{"events.>"},
	}
	cfg.Storage = nats.FileStorage*/

	sub, err := js.PullSubscribe("", "", nats.BindStream(streamName))
	if err != nil {
		panic("Cant subscribe : " + err.Error())
	}
	defer sub.Unsubscribe()
	for {
		/*msg, err := sub.NextMsg(time.Second * 1)
		if err == nats.ErrTimeout {
			//break
			msg.Ack()
			continue
		}*/
		/*if err != nil && err != nats.ErrTimeout {
			panic(err)

		}*/
		msgs, err := sub.Fetch(1, nats.MaxWait(time.Second*6))
		if err != nil {
			panic("Cant fetch msg: " + err.Error())
		}
		ords := &ents.Orders{}
		json.Unmarshal(msgs[0].Data, ords)
		fmt.Println("ORDS : ", *ords)
		fmt.Println("MSG : ", string(msgs[0].Data))
		msgs[0].Ack()
	}
}
