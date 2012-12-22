package main

import (
	"bufio"
	"flag"
	"github.com/tarm/goserial"
	"io"
	"log"
	"strings"
	"fmt"
)

var (
	port = flag.String("port", "/dev/cu.HUAWEIMobile-Pcui", "port modem is connected on e.g. /dev/cu.HUAWEIMobile-Pcui on Mac OSX")
)

type GsmModem struct {
	port       string
	readWriter io.ReadWriteCloser
}

func NewGsmModem(port string) *GsmModem {
	return &GsmModem{port: port}
}
func (g *GsmModem) Connect() (io.ReadWriteCloser, error) {
	c := &serial.Config{Name: g.port, Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	//defer s.Close()
	g.readWriter = s
	return g.readWriter, nil
}

func (g *GsmModem) SendCommand(command string) (lines []string, response string, isOK bool) {
	_, err := g.readWriter.Write([]byte(command + "\n"))
	if err != nil {
		log.Fatal(err)
	}
	r := bufio.NewReader(g.readWriter)
	for {
		read, _, err := r.ReadLine()
		if err != nil {
			isOK = false
			return
		}
		strread := string(read)
		fmt.Println(strread)
		if strings.Contains(strread, "OK") {
			response = strread
			isOK = true
			return
		}
		if strings.Contains(strread, "ERROR") {
			response = strread
			isOK = false
			return
		}

		lines = append(lines, strread)
	}
	return 
}



func main() {
	flag.Parse()
	if len(*port) == 0 {
		log.Fatal("--port flag not provided")
	}

	modem := NewGsmModem(*port)
	modem.Connect()
	_, _, isOK := modem.SendCommand("AT+CMGF=1")
	if isOK != true {
		log.Fatal("failed to set modem on text mode")
	}
  for{
    resp, lines, _ := modem.SendCommand("AT+CMGL=ALL")
    fmt.Println(resp)
    fmt.Println(lines)
  }
}
