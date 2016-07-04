package main

import (
	"log"
	"net"
	"os"
	"sync"
	"testing"
)

func TestMain(m *testing.M) {
	go main()
	os.Exit(m.Run())
}

const len144 = len("XX144XX")
const len2500 = len("XX2500XX")
const len24A7 = len("EE24A7EE")
const lenZZ = len("ZZ")

// Test normal client operation
func TestRequestSquare(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:4716")
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		buf := make([]byte, 0)
		pad := make([]byte, 1)
		for i := 0; i < len144; i++ {
			conn.Read(pad)
			buf = append(buf, pad[0])
		}

		if string(buf) != "YY144YY" {
			t.Error(`Expected YY144YY, got: `, string(buf))
		}

		buf = buf[:0]
		for i := 0; i < len2500; i++ {
			conn.Read(pad)
			buf = append(buf, pad[0])
		}

		if string(buf) != "YY2500YY" {
			t.Error(`Expected YY2500YY, got: `, string(buf))
		}
	}()

	conn.Write([]byte("XX12XX"))
	conn.Write([]byte("XX50XX"))
	wg.Wait()
}

// Test client sending a malformed number
func TestErrorRequestSquare(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:4716")
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		buf := make([]byte, 0)
		pad := make([]byte, 1)
		for i := 0; i < len24A7; i++ {
			conn.Read(pad)
			buf = append(buf, pad[0])
		}

		if string(buf) != "EE24A7EE" {
			t.Error(`Expected EE24A7EE, got: `, string(buf))
		}
	}()

	conn.Write([]byte("XX24A7XX"))
	wg.Wait()
}

// Test request to close the connection
func TestRequestClose(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:4716")
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		buf := make([]byte, 0)
		pad := make([]byte, 1)
		for i := 0; i < lenZZ; i++ {
			conn.Read(pad)
			buf = append(buf, pad[0])
		}

		if string(buf) != "ZZ" {
			t.Error(`Expected ZZ, got: `, string(buf))
		}
	}()

	conn.Write([]byte("ZZ"))
	wg.Wait()
	// The client must close connection after
	// receiving server "ZZ"
	conn.Close()
}
