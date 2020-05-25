package pcr2

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/tarm/serial"
)

const (
	strCommandNotFound = "command not found"
)

var (
	ErrorCommandNotFound = errors.New("Command not found on device")
)

type SerialTransport struct {
	port  *serial.Port
	debug bool
}

func Open(portname string) (*SerialTransport, error) {
	config := &serial.Config{
		Name:        portname,
		Baud:        19200,
		ReadTimeout: 1,
		Size:        8,
		Parity:      serial.ParityNone,
		StopBits:    1,
	}
	port, err := serial.OpenPort(config)
	if err != nil {
		return nil, fmt.Errorf("Failed to open serial port %s: %w", portname, err)
	}

	return &SerialTransport{
		port: port,
	}, nil
}

func (s *SerialTransport) Write(in string) (string, error) {
	if s.debug {
		debugPrint("Sending %s", toHexArray([]byte(in)))
	}
	_, err := s.port.Write([]byte(in))
	if err != nil {
		return "", fmt.Errorf("Failed to write to serial port: %w", err)
	}

	buf := make([]byte, 2048)
	nc := 0
	for {
		n, err := s.port.Read(buf[nc:])
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", fmt.Errorf("Failed to read data from serial port: %w", err)
		}
		if s.debug {
			debugPrint("Received part %s", toHexArray(buf[nc:nc+n]))
		}
		nc += n
		if buf[nc-2] == 0x0D && buf[nc-1] == 0x0A {
			break
		}
	}
	if s.debug {
		debugPrint("Received complete response %s", toHexArray(buf[:nc]))
	}
	outStr := string(buf[:nc])
	outStr = strings.TrimRight(outStr, "\r\n")
	if outStr == strCommandNotFound {
		return "", ErrorCommandNotFound
	}
	return outStr, nil
}

func (s *SerialTransport) Debug(debug bool) {
	s.debug = debug
}

func toHexArray(in []byte) string {
	out := ""
	for _, b := range in {
		out += fmt.Sprintf("0x%X ", b)
	}
	return out
}

func debugPrint(format string, args ...interface{}) {
	fmt.Printf("[Serial] "+format+"\n", args)
}

func (s *SerialTransport) Close() error {
	return s.port.Close()
}
