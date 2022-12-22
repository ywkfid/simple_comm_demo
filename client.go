package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.conn = conn
	return client
}

func (client *Client) menu() bool {
	var flag int
	fmt.Println("1.public")
	fmt.Println("2.private")
	fmt.Println("3.rename")
	fmt.Println("4.exit")

	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("please input 0~3")
		return false
	}
}

func (client *Client) PublicChat() {
	var chatMsg string

	fmt.Println("please input message, input exit can return")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn Write err:", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println("please input message, input exit can return")
		fmt.Scanln(&chatMsg)
	}
}

func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err:", err)
		return
	}
}

func (client *Client) PrivateChat() {
	var toUserName string
	var chatMsg string

	client.SelectUsers()
	fmt.Println("please input name who you want to talk")
	fmt.Scanln(&toUserName)

	for toUserName != "exit" {
		fmt.Println("please input message, input exit can return")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + toUserName + "|" + chatMsg + "\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err:", err)
					break
				}
			}
			chatMsg = ""
			fmt.Println("please input message, input exit can return")
			fmt.Scanln(&chatMsg)
		}
		client.SelectUsers()
		fmt.Println("please input name who you want to talk")
		fmt.Scanln(&toUserName)
	}

}

func (client *Client) UpdateName() bool {
	fmt.Println("please input username")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}
	return true
}

//deal response from server
func (client *Client) DealResponse() {
	//when client.conn buf exist
	io.Copy(os.Stdout, client.conn)
	/*for {
		buf := make()
		client.conn.Read(buf)
		fmt.Println(buf)
	}*/

}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}
		switch client.flag {
		//public
		case 1:
			client.PublicChat()
			break
		//private
		case 2:
			client.PrivateChat()
			break
		//rename
		case 3:
			client.UpdateName()
			break
		}
	}
}

var serverIp string
var serverPort int

// ./client -ip 127.0.0.1 -port 8888
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "set ip (Default:127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "set port (Default:8888)")
}

func main() {
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("connect server failed")
		return
	}

	go client.DealResponse()

	fmt.Println("connect server succeed")

	client.Run()

}
