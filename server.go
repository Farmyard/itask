package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"./server"
)

func main() {
	ip := flag.String("ip", "0.0.0.0", "http listen ip")
	port := flag.Int("port", 84, "http listen port")
	flag.Parse()

	server.NewPush()

	addr := *ip + ":" + strconv.Itoa(*port)
	log.Println("Serving HTTP on " + addr)

	go startTCP()

	r := server.NewRoute()
	log.Fatal(http.ListenAndServe(addr, r))
}

func startTCP() {
	listener, err := net.Listen("tcp", "0.0.0.0:8484")
	log.Println("Serving TCP Starting ")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}
		listenClient(conn)
	}
}

var uuid string

func listenClient(conn net.Conn) {
	server.Stop = false
	buf := make([]byte, 50)

	defer conn.Close()

	for {
		n, err := conn.Read(buf)
		if err != nil {
			m := server.NewModel()
			task, _ := m.GetTaskInfoByKey(uuid)

			device, _ := m.GetDeviceInfoByKey(uuid)

			params := make(map[string]interface{})

			params["key"] = task.Uuid
			params["title"] = task.Title
			params["category"] = task.Category
			params["body"] = task.Body
			params["url"] = task.Url

			p := server.NewPush()
			for !server.Stop {
				p.PostPush(task.Category, task.Title, task.Body, device.Token, params)
				time.Sleep(10 * time.Second)
			}

			return
		}
		uuid = string(buf[0:n])
	}
}
