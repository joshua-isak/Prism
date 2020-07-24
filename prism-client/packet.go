// Code for the prism tcp packet protocol
package main

import (
	//"bufio"
	//"encoding/binary"
	"fmt"
	"net"
	"encoding/hex"
)

// PacketType :  enumeration for packet types
type PacketType int

// see above
const (
	Initial PacketType = 1
	Welcome PacketType = 2
	ClientConnect PacketType = 3
	ClientDisconnect PacketType = 4
	GeneralMessage PacketType = 5
	Received PacketType = 255
)


// Packet : a byte buffer with some vars to aid preparing a packet to send
type Packet struct{
	data []byte
	seek int
}



// NewPacket : Packet constructor
func NewPacket(t PacketType) Packet {
	var newPacket Packet

	// If this is a received packet, don't change anything
	if t == Received{
		return newPacket
	}

	// Write the packet type
	newPacket.data = append(newPacket.data, uint8(t))

	return newPacket
}


// PrepInitial : prepare the "Inital" packet type
func (p *Packet) PrepInitial(username string) {
	// Write in the len of the username then the username
	p.data = append(p.data, uint8(len(username)))
	p.data = append(p.data, []byte(username)...)
}


// PrepGeneralMessage : prepare the "GeneralMessage" packet type
func (p *Packet) PrepGeneralMessage(username string, message []byte, encrypted bool) {
	// Write in the len of the username then the username
	p.data = append(p.data, uint8(len(username)))
	p.data = append(p.data, []byte(username)...)

	// Write in whether the message is encrypted
	buf := make([]byte, 21 - len(username))
	var i uint8
	if encrypted { i = 1 }
	p.data = append(p.data, buf...)
	p.data = append(p.data, i)

	// Write in the len of the message then the message
	p.data = append(p.data, uint8(len(message)))
	p.data = append(p.data, message...)
}


// PrintData : prints out packet data
func (p *Packet) PrintData() {
	fmt.Println( p.data)
}


// PrintDataHex : prints out packet data as a hex dump
func (p *Packet) PrintDataHex() {
	fmt.Println(hex.Dump(p.data))	// DEBUG
}


// Send : writes packet buffer to tcp socket to send
func (p *Packet) Send(c net.Conn) {
	c.Write(p.data)
}


// Broadcast : sends packets to all connections in map m
func (p *Packet) Broadcast(m map[int]net.Conn ) {
	for _, connection := range m {
		p.Send(connection)
	}
}


// ReadBytes : reads in a slice from p.data and updates p.seek, n must be positive .. Maybe add errs??
func (p *Packet) ReadBytes(n int) []byte {
	start := p.seek
	end := start + n
	output := p.data[start:end]

	p.seek = end

	return output
}


// ReadUint8 : output a uint8 from p.data and update the seek
func (p *Packet) ReadUint8() uint8 {
	output := p.data[p.seek]
	p.seek++

	return output
}


// ReadString : Read len bytes from p.data and output them as a string then update p.seek
func (p *Packet) ReadString(len int) string {
	if len < 1 {
		return ""
	}

	start := p.seek
	end := start + len
	output := p.data[start:end]

	p.seek = end

	return string(output)
}


// ReadBool : output a bool from p.data and update the seek
func (p *Packet) ReadBool() bool {
	var output bool

	b := p.data[p.seek]
	p.seek++

	if b == 1 {				// TODO, there is SOOO a better way
		output = true
	} else {
		output = false
	}

	return output
}
