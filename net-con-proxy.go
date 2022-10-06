package Network

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

/*
TODO:
	- Write a more complex error system
	- Write a more complex wait group system
	- Secure payloads more than they are

*/

var swg sync.WaitGroup

func proxy(from io.Reader, to io.Writer) error {
	fromw, fromiw := from.(io.Writer)
	tor, toir := to.(io.Reader)
	if toir && fromiw {
		go func() {
			_, _ = io.Copy(fromw, tor)
		}()
	}
	_, x := io.Copy(to, from)
	return x
}

func Proxy() {
	// proxy server
	server, x := net.Listen("tcp", net.JoinHostPort("127.0.0.1", "8080")) // Random port
	if x != nil {
		log.Fatal(x)
	}
	swg.Add(1)
	go func() {
		defer swg.Done()
		for {
			conn, x := server.Accept()
			if x != nil {
				return
			}
			go func(conn net.Conn) {
				defer conn.Close()
				for {
					buffer := make([]byte, 1024)
					n, x := conn.Read(buffer)
					if x != nil {
						if x != io.EOF {
							fmt.Println("SERVER ERROR: ", x)
						}
						return
					}
					switch msg := string(buffer[:n]); msg {
					case "ping":
						fmt.Println("Server -> Got ping")
						_, x = conn.Write([]byte("pong"))
					default:
						fmt.Printf("Server -> Got %s\n", msg)
						_, x = conn.Write(buffer[:n])
					}
					if x != nil {
						if x != io.EOF {
							fmt.Println("SERVER ERROR: ", x)
						}
						return
					}
				}

			}(conn)

		}
	}()
	// Client
	proxy_client, x := net.Listen("tcp", "127.0.0.1:")
	if x != nil {
		fmt.Println("SERVER [PROXY CLIENT] ERROR: port=8080", x)
		os.Exit(0)
	}
	swg.Add(1)
	go func() {
		for {
			connection, x := proxy_client.Accept()
			if x != nil {
				return
			}
			go func(from net.Conn) {
				to, x := net.Dial("tcp", net.JoinHostPort(server.Addr().String(), "8080"))
				if x != nil {
					fmt.Println("SERVER [PROXY CLIENT] ERROR: ", x)
					return
				}
				defer to.Close()
				x = proxy(from, to)
				if x != nil && x != io.EOF {
					fmt.Println("SERVER [PROXY CLIENT] ERROR: ", x)

				}
			}(connection)
		}
	}()
	// Client Server
	connect, x := net.Dial("tcp", server.Addr().String())
	if x != nil {
		fmt.Println("SERVER [PROXY CLIENT] ERROR: ", x)
	}
	payloads := []struct{ Message, Reply string }{
		{"ping", "reply1"},
	}
	for i, m := range payloads {
		_, x = connect.Write([]byte(m.Message))
		if x != nil {
			fmt.Println("SERVER [ CLIENT SERVER ] ERROR: ", x)
			os.Exit(1)
		}
		buffer := make([]byte, 1024)
		n, x := connect.Read(buffer)
		if x != nil {
			fmt.Println("")
		}
		actual := string(buffer[:n])
		fmt.Printf("%q -> Proxy -> %q\n", m.Message, m.Reply)
		if actual != m.Reply {
			fmt.Printf("%d: Expected reply: %q, actual reply = %q", i, m.Reply, actual)
		}
	}
	_ = connect.Close()
	fmt.Println("[+] CLosed connection on client - 0x00")
	_ = proxy_client.Close()
	fmt.Println("[+] Closed connection on client server - 0x00")
	_ = server.Close()
	fmt.Println("[+] Closed connection on proxy server")

}
