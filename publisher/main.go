package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nats-io/nats.go"
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
	defer js.DeleteStream(streamName)
	_, err = js.AddStream(&cfg)
	if err != nil {
		if err != nats.ErrStreamNameAlreadyInUse {
			fmt.Println("REC ERROR :", err)
			fmt.Println("JS ERROR :", nats.ErrStreamNameAlreadyInUse)
			panic(fmt.Errorf("Cant add stream: %s", err.Error()))
		} else {
			fmt.Println("Deleting Stream")
			js.DeleteStream(streamName)
			js.AddStream(&cfg)
		}

	}
	fmt.Println("created the stream")

	file, err := os.Open("../model.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fileSize, err := file.Seek(0, os.SEEK_END)
	file.Seek(0, os.SEEK_SET)
	if err != nil {
		panic("Cant seek file" + err.Error())
	}
	fmt.Println("FILESIZE IS ", fileSize)
	buf := make([]byte, fileSize)
	_, err = file.Read(buf)
	fmt.Println("BYTES READ ", string(buf))
	if err != nil {
		panic("Cant read from file : " + err.Error())
	}
	fmt.Println("SUCCESSFULLY READ")

	//data, err := json.Marshal(buf)
	if err != nil {
		panic("Cant Marshal file : " + err.Error())
	}
	fmt.Println("SUCCESSFULLY MARSHALLED")
	for {
		time.Sleep(time.Second * 5)
		fmt.Println("PUBLISHING")
		js.Publish("events.test", buf)
	}

}
