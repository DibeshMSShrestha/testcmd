/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	socketio_client "github.com/zhouhui8915/go-socket.io-client"
)

func main() {
	opts := &socketio_client.Options{
		Transport: "websocket",
		Query:     make(map[string]string),
	}
	opts.Query["user"] = "user"
	opts.Query["pwd"] = "pass"
	uri := "http://localhost:8000"

	client, err := socketio_client.NewClient(uri, opts)
	if err != nil {
		log.Printf("NewClient error:%v\n", err)
		return
	}

	client.On("error", func() {
		log.Printf("on error\n")
	})
	client.On("connection", func() {
		log.Printf("on connect\n")
	})
	client.On("reply", func(msg string) {
		log.Printf("%v\n", msg)
	})
	client.On("room", func(msg string) {
		log.Printf("%v\n", msg)
	})
	client.On("disconnection", func() {
		log.Printf("on disconnect\n")
	})

	var cmdPrint = &cobra.Command{
		Use:   "login",
		Short: "enter the login name",
		Long:  `enter the login name`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("enter name:")
			reader := bufio.NewReader(os.Stdin)
			name, _, _ := reader.ReadLine()
			client.Emit("name", string(name))

			client.Emit("room", "")

			room, _, _ := reader.ReadLine()
			client.Emit("join", string(room))

			for {
				data, _, _ := reader.ReadLine()
				command := string(data)
				client.Emit("notice", command)
			}

		},
	}

	var rootCmd = &cobra.Command{Use: "testcmd"}
	rootCmd.AddCommand(cmdPrint)
	rootCmd.Execute()
}
