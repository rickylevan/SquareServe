package main

import (
	"log"
	"net"
	"strconv"
)

const (
	base    = 10
	bitSize = 64
)

func main() {
	ls, err := net.Listen("tcp", "localhost:4716")
	log.Println("Server listening @ localhost:4716")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ls.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	log.Println("Server handling a new connection")
	cmdHead := make([]byte, 2)
	cmdTail := make([]byte, 2)
	pad := make([]byte, 1)

	for {
		conn.Read(pad)
		cmdHead[0] = pad[0]
		conn.Read(pad)
		cmdHead[1] = pad[0]

		switch {
		case string(cmdHead) == "XX":
			numberBuf := make([]byte, 0)
			for {
				conn.Read(pad)
				if !(pad[0] >= '0' && pad[0] <= '9') {
					break
				}
				numberBuf = append(numberBuf, pad[0])
			}
			cmdTail[0] = pad[0]
			conn.Read(pad)
			cmdTail[1] = pad[0]

			if string(cmdTail) != "XX" {
				// robustness against stupid or malicious clients
				// is not part of the spec
				log.Fatal(`Error: No matching "XX" tail`)
			}

			square, err := strconv.ParseInt(string(numberBuf), base, bitSize)
			square *= square
			if err != nil {
				conn.Write([]byte("EE"))
				conn.Write(numberBuf)
				conn.Write([]byte("EE"))
			} else {
				conn.Write([]byte("YY"))
				conn.Write([]byte(strconv.Itoa(int(square))))
				conn.Write([]byte("YY"))
			}

		case string(cmdHead) == "ZZ":
			conn.Write([]byte("ZZ"))

		default:
			// robustness against stupid or malicious clients
			// is not part of the spec
			log.Fatal(`Client commands must begin with "XX" or "ZZ"`)
		}
	}
}
