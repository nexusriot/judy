package main

import j "github.com/nexusriot/judy/server"

func main() {
	judy, err := j.NewJudyServer()
	if err != nil {
		panic(err)
	}
	judy.Run("0.0.0.0:1337")
}
