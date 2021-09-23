package main

import j "github.com/nexusriot/judy/client"

func main() {
	client, err := j.NewClient()
	if err != nil {
		panic(err)
	}
	client.Run()
}
