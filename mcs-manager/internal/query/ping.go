package query

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"
)

// ServerStatus represents the JSON response from Minecraft server
type ServerStatus struct {
	Version Version `json:"version"`
	Players Players `json:"players"`
	Description DescriptionData `json:"description"`
}

type Version struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type Players struct {
	Max    int `json:"max"`
	Online int `json:"online"`
}

type DescriptionData struct {
	Text string `json:"text"`
}

// PingServer sends a status ping to a Minecraft server and gets player info
func PingServer(host string, port int, timeout time.Duration) (*ServerStatus, error) {
	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server at %s: %w", addr, err)
	}
	defer conn.Close()

	// Set read deadline
	conn.SetReadDeadline(time.Now().Add(timeout))

	// Handshake packet
	handshakeData := bytes.NewBuffer([]byte{})
	
	// Protocol version (use 762 for 1.20.x compatibility)
	writeVarInt(handshakeData, 762)
	writeString(handshakeData, host)
	binary.Write(handshakeData, binary.BigEndian, uint16(port))
	writeVarInt(handshakeData, 1) // Next state: status
	
	// Wrap with packet ID
	packet := bytes.NewBuffer([]byte{})
	writeVarInt(packet, 0) // Handshake packet ID
	packet.Write(handshakeData.Bytes())
	
	// Write handshake packet
	if err := writePacket(conn, packet.Bytes()); err != nil {
		return nil, fmt.Errorf("failed to send handshake: %w", err)
	}

	// Status request packet (empty)
	statusPacket := bytes.NewBuffer([]byte{})
	writeVarInt(statusPacket, 0) // Status request packet ID
	if err := writePacket(conn, statusPacket.Bytes()); err != nil {
		return nil, fmt.Errorf("failed to send status request: %w", err)
	}

	// Read response
	statusResponse, err := readPacket(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to read status response: %w", err)
	}

	// Skip packet ID (first VarInt)
	reader := bytes.NewReader(statusResponse)
	readVarInt(reader) // Skip packet ID (0x00 for status response)

	// Read JSON string length and data
	jsonLen := readVarInt(reader)
	jsonData := make([]byte, jsonLen)
	if _, err := io.ReadFull(reader, jsonData); err != nil {
		return nil, fmt.Errorf("failed to read JSON data: %w", err)
	}

	// Parse JSON
	var status ServerStatus
	if err := json.Unmarshal(jsonData, &status); err != nil {
		return nil, fmt.Errorf("failed to parse status JSON: %w", err)
	}

	return &status, nil
}

// Helper functions for reading/writing VarInt and other types

func writeVarInt(buf *bytes.Buffer, value int32) {
	v := uint32(value)
	for {
		b := byte(v & 0x7F)
		v >>= 7
		if v != 0 {
			b |= 0x80
		}
		buf.WriteByte(b)
		if v == 0 {
			break
		}
	}
}

func readVarInt(r io.Reader) int32 {
	result := int32(0)
	shift := 0
	for {
		b := make([]byte, 1)
		if _, err := io.ReadFull(r, b); err != nil {
			return 0
		}
		result |= int32(b[0]&0x7F) << shift
		if b[0]&0x80 == 0 {
			break
		}
		shift += 7
	}
	return result
}

func writeString(buf *bytes.Buffer, s string) {
	data := []byte(s)
	writeVarInt(buf, int32(len(data)))
	buf.Write(data)
}

func writePacket(conn net.Conn, data []byte) error {
	// Write packet length as VarInt, then packet data
	lengthBuf := bytes.NewBuffer([]byte{})
	writeVarInt(lengthBuf, int32(len(data)))
	
	fullPacket := append(lengthBuf.Bytes(), data...)
	_, err := conn.Write(fullPacket)
	return err
}

func readPacket(conn net.Conn) ([]byte, error) {
	// Read packet length
	lengthBuf := make([]byte, 5) // Max 5 bytes for VarInt
	n, err := conn.Read(lengthBuf)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(lengthBuf[:n])
	length := readVarInt(reader)

	// Read packet data
	packetData := make([]byte, length)
	totalRead := 0
	for totalRead < int(length) {
		n, err := conn.Read(packetData[totalRead:])
		if err != nil {
			return nil, err
		}
		totalRead += n
	}

	return packetData, nil
}
