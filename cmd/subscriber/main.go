package main

import (
	cmd "github.com/sabrs0/L0_WB/cmd"
	srv "github.com/sabrs0/L0_WB/server"
	sub "github.com/sabrs0/L0_WB/subscriber"
)

func main() {
	cache, err := sub.RecoverCache()
	if err != nil {
		panic(err)
	}
	sub, err := sub.NewSubscriber(&cache)
	if err != nil {
		panic(err)
	}
	go srv.StartServer(cmd.Addr, &cache)
	err = sub.Subscribe(cmd.StreamName)
	if err != nil {
		panic(err)
	}

}
