package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

var streamName = "EVENTS"

func main() {
	// Use the env variable if running in the container, otherwise use the default.
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

	// Access `JetStreamContext` which provides methods to create
	// streams and consumers as well as convenience methods for publishing
	// to streams and implicitly creating consumers through `*Subscribe*`
	// methods (which will be discussed in examples focused on consumers).
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
	cfg := nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{"events.>"},
	}

	// JetStream provides both file and in-memory storage options. For
	// durability of the stream data, file storage must be chosen to
	// survive crashes and restarts. This is the default for the stream,
	// but we can still set it explicitly.
	cfg.Storage = nats.FileStorage
	// Finally, let's add/create the stream with the default (no) limits.
	_, err = js.AddStream(&cfg)
	if err != nil {
		if err != nats.ErrStreamNameAlreadyInUse {
			fmt.Println("REC ERROR :", err)
			fmt.Println("JS ERROR :", nats.ErrStreamNameAlreadyInUse)
			panic(fmt.Errorf("Cant add stream: %s", err.Error()))
		} else {
			js.DeleteStream(streamName)
			js.AddStream(&cfg)
		}

	}
	fmt.Println("created the stream")
	// Let's publish a few messages which are received by the stream since
	// they match the subject bound to the stream. The `js.Publish` method
	// is a convenience for sending a `nc.Request` and waiting for the
	// acknowledgement.
	js.Publish("events.page_loaded", []byte("Hello world 0"))
	js.Publish("events.mouse_clicked", []byte("Hello world 1"))
	js.Publish("events.mouse_clicked", []byte("Hello world 2"))
	js.Publish("events.page_loaded", []byte("Hello world 3"))
	js.Publish("events.mouse_clicked", []byte("Hello world 4"))
	js.Publish("events.input_focused", []byte("Hello world 5"))
	fmt.Println("published 6 messages")

	// There is also is an async form in which the client batches the
	// messages to the server and then asynchronously receives the
	// the acknowledgements.
	js.PublishAsync("events.input_changed", []byte("Hello world 0 sync"))
	js.PublishAsync("events.input_blurred", []byte("Hello world 1 sync"))
	js.PublishAsync("events.key_pressed", []byte("Hello world 2 sync"))
	js.PublishAsync("events.input_focused", []byte("Hello world 3 sync"))
	js.PublishAsync("events.input_changed", []byte("Hello world 4 sync"))
	js.PublishAsync("events.input_blurred", []byte("Hello world 5 sync"))

	// For a given batch, we select on a channel returned from
	// `js.PublishAsyncComplete`.
	select {
	case <-js.PublishAsyncComplete():
		fmt.Println("published 6 messages")
	case <-time.After(time.Second):
		log.Fatal("publish took too long")
	}

	// Checking out the stream info, we can see how many messages we
	// have.
	printStreamState(js, cfg.Name)

	RecMsg(js)
	// Stream configuration can be dynamically changed. For example,
	// we can set the max messages limit to 10 and it will truncate the
	// two initial events in the stream.
	time.Sleep(time.Second * 5)
	cfg.MaxMsgs = 10
	js.UpdateStream(&cfg)
	fmt.Println("set max messages to 10")

	// Checking out the info, we see there are now 10 messages and the
	// first sequence and timestamp are based on the third message.
	printStreamState(js, cfg.Name)

	// Limits can be combined and whichever one is reached, it will
	// be applied to truncate the stream. For example, let's set a
	// maximum number of bytes for the stream.
	cfg.MaxBytes = 300
	js.UpdateStream(&cfg)
	fmt.Println("set max bytes to 300")

	// Inspecting the stream info we now see more messages have been
	// truncated to ensure the size is not exceeded.
	printStreamState(js, cfg.Name)

	// Finally, for the last primary limit, we can set the max age.
	cfg.MaxAge = time.Second
	js.UpdateStream(&cfg)
	fmt.Println("set max age to one second")

	// Looking at the stream info, we still see all the messages..
	printStreamState(js, cfg.Name)

	// until a second passes.
	fmt.Println("sleeping one second...")
	time.Sleep(time.Second)

	printStreamState(js, cfg.Name)
}

// This is just a helper function to print the stream state info 😉
func printStreamState(js nats.JetStreamContext, name string) {
	info, _ := js.StreamInfo(name)
	b, _ := json.MarshalIndent(info.State, "", " ")
	fmt.Println("inspecting stream info")
	fmt.Println(string(b))
}
func RecMsg(js nats.JetStreamContext) {
	sub, err := js.SubscribeSync("events.>") //js.PullSubscribe("", "", nats.BindStream(streamName))
	if err != nil {
		panic("Cant subscribe : " + err.Error())
	}

	for {
		msg, err := sub.NextMsg(time.Second * 1)
		if err == nats.ErrTimeout {
			break
		}
		if err != nil {
			panic(err)

		}
		fmt.Println(string(msg.Data))
		msg.Ack()
	}
}
