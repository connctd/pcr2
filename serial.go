package pcr2

import (
	"errors"
	"fmt"
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
	port *serial.Port
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
	_, err := s.port.Write([]byte(in))
	if err != nil {
		return "", fmt.Errorf("Failed to write to serial port: %w", err)
	}

	buf := make([]byte, 2048)
	nc := 0
	for {
		n, err := s.port.Read(buf[nc:])
		if err != nil {
			return "", fmt.Errorf("Failed to read data from serial port: %w", err)
		}
		nc += n
		if buf[nc-2] == 0x0D && buf[nc-1] == 0x0A {
			break
		}
	}

	outStr := string(buf[:nc])
	outStr = strings.TrimRight(outStr, "\r\n")
	if outStr == strCommandNotFound {
		return "", ErrorCommandNotFound
	}
	return outStr, nil
}
