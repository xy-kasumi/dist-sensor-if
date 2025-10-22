package main

import (
	"fmt"
	"sync"

	"go.bug.st/serial"
)

const (
	STX = 0x02
	ETX = 0x03
	ACK = 0x06
	NAK = 0x15

	CMD_COMMAND = 0x43 // 'C'
	CMD_WRITE   = 0x57 // 'W'
	CMD_READ    = 0x52 // 'R'

	// Device error codes
	ERR_INVALID_ADDR    = 0x02
	ERR_INVALID_BCC     = 0x04
	ERR_INVALID_COMMAND = 0x05
	ERR_INVALID_PARAM   = 0x06
	ERR_OUT_OF_RANGE    = 0x07
)

type CD22 struct {
	port   serial.Port
	portMu sync.Mutex
}

func NewCD22(portName string, baudRate int) (*CD22, error) {
	mode := &serial.Mode{BaudRate: baudRate}
	p, err := serial.Open(portName, mode)
	if err != nil {
		return nil, err
	}
	return &CD22{port: p}, nil
}

func (dev *CD22) Send(command byte, data uint16) (uint16, error) {
	dev.portMu.Lock()
	defer dev.portMu.Unlock()

	// Send command
	if err := dev.writeAll(buildFrame(command, data)); err != nil {
		return 0, err
	}

	// Read response (6 bytes: STX + type + data1 + data2 + ETX + BCC)
	buf, err := dev.readAll(6)
	if err != nil {
		return 0, err
	}
	data, err = parseResponse(buf)
	if err != nil {
		return 0, err
	}
	return data, nil
}

func (dev *CD22) writeAll(data []byte) error {
	total := 0
	for total < len(data) {
		n, err := dev.port.Write(data[total:])
		if err != nil {
			return fmt.Errorf("write failed: %w", err)
		}
		total += n
	}
	return nil
}

func (dev *CD22) readAll(n int) ([]byte, error) {
	buf := make([]byte, n)
	total := 0
	for total < n {
		nRead, err := dev.port.Read(buf[total:])
		if err != nil {
			return nil, fmt.Errorf("read failed: %w", err)
		}
		total += nRead
	}
	return buf, nil
}

func buildFrame(command byte, data uint16) []byte {
	data1 := byte(data >> 8)
	data2 := byte(data & 0xFF)

	payload := []byte{command, data1, data2}
	bcc := calculateBCC(payload)
	return []byte{STX, command, data1, data2, ETX, bcc}
}

func calculateBCC(data []byte) byte {
	bcc := byte(0)
	for _, b := range data {
		bcc ^= b
	}
	return bcc
}

func parseResponse(frame []byte) (uint16, error) {
	if len(frame) < 6 {
		return 0, fmt.Errorf("response too short: %d bytes", len(frame))
	}
	if frame[0] != STX || frame[4] != ETX {
		return 0, fmt.Errorf("missing STX/ETX")
	}

	// Verify BCC
	payload := frame[1 : len(frame)-2]
	expectedBCC := calculateBCC(payload)
	actualBCC := frame[len(frame)-1]
	if expectedBCC != actualBCC {
		return 0, fmt.Errorf("BCC mismatch: expected %02X, got %02X", expectedBCC, actualBCC)
	}

	responseType := frame[1]
	data1 := frame[2]
	data2 := frame[3]

	if responseType == ACK {
		return (uint16(data1) << 8) | uint16(data2), nil
	} else if responseType == NAK {
		return 0, fmt.Errorf("device error: %s", deviceErrorToString(data1))
	}
	return 0, fmt.Errorf("invalid response type: %02X", responseType)
}

func deviceErrorToString(errCode byte) string {
	switch errCode {
	case ERR_INVALID_ADDR:
		return "invalid address"
	case ERR_INVALID_BCC:
		return "invalid BCC value"
	case ERR_INVALID_COMMAND:
		return "invalid command"
	case ERR_INVALID_PARAM:
		return "invalid setting value (parameter specification invalid)"
	case ERR_OUT_OF_RANGE:
		return "invalid setting value (out-of-range specified)"
	default:
		return fmt.Sprintf("unknown error code: %02X", errCode)
	}
}
