// Copyright 2015 Nevio Vesic
// Please check out LICENSE file for more information about what you CAN and what you CANNOT do!
// Basically in short this is a free software for you to do whatever you want to do BUT copyright must be included!
// I didn't write all of this code so you could say it's yours.
// MIT License

package goesl

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// SocketConnection main connection against ESL
type SocketConnection struct {
	connection net.Conn
	err        chan error
	m          chan *Message
	reader     *bufio.Reader
	textreader *textproto.Reader
}

// create SocketConnection instance
func newConnection(c net.Conn) *SocketConnection {
	result := &SocketConnection{
		connection: c,
		err:        make(chan error),
		m:          make(chan *Message),
		reader:     bufio.NewReaderSize(c, ReadBufferSize),
	}
	result.textreader = textproto.NewReader(result.reader)
	return result
}

// Will establish timedout dial against specified address. In this case, it will be freeswitch server
func dial(network string, addr string, timeout time.Duration) (*SocketConnection, error) {
	c, err := net.DialTimeout(network, addr, timeout)
	if err != nil {
		return nil, err
	}
	return newConnection(c), err
}

// Send - Will send raw message to open net connection
func (c *SocketConnection) Send(cmd string) error {
	if strings.Contains(cmd, "\r\n") {
		return fmt.Errorf(EInvalidCommandProvided, cmd)
	}

	fmt.Fprintf(c.connection, "%s\r\n\r\n", cmd)
	return nil
}

// SendMany - Will loop against passed commands and return 1st error if error happens
func (c *SocketConnection) SendMany(cmds []string) error {

	for _, cmd := range cmds {
		if err := c.Send(cmd); err != nil {
			return err
		}
	}

	return nil
}

// Execute - Helper fuck to execute commands with its args and sync/async mode
func (c *SocketConnection) Execute(command, args string, sync bool) (err error) {
	return c.SendMsg(map[string]string{
		"call-command":     "execute",
		"execute-app-name": command,
		"execute-app-arg":  args,
		"event-lock":       strconv.FormatBool(sync),
	}, "", "")
}

// ExecuteUUID - Helper fuck to execute uuid specific commands with its args and sync/async mode
func (c *SocketConnection) ExecuteUUID(uuid string, command string, args string, sync bool) (err error) {
	return c.SendMsg(map[string]string{
		"call-command":     "execute",
		"execute-app-name": command,
		"execute-app-arg":  args,
		"event-lock":       strconv.FormatBool(sync),
	}, uuid, "")
}

// SendMsg - Basically this func will send message to the opened connection
func (c *SocketConnection) SendMsg(msg map[string]string, uuid, data string) (err error) {

	b := bytes.NewBufferString("sendmsg")

	if uuid != "" {
		if strings.Contains(uuid, "\r\n") {
			return fmt.Errorf(EInvalidCommandProvided, msg)
		}

		b.WriteString(" " + uuid)
	}

	b.WriteString("\n")

	for k, v := range msg {
		if strings.Contains(k, "\r\n") {
			return fmt.Errorf(EInvalidCommandProvided, msg)
		}

		if v != "" {
			if strings.Contains(v, "\r\n") {
				return fmt.Errorf(EInvalidCommandProvided, msg)
			}

			b.WriteString(fmt.Sprintf("%s: %s\n", k, v))
		}
	}

	b.WriteString("\n")

	if msg["content-length"] != "" && data != "" {
		b.WriteString(data)
	}

	if _, err := b.WriteTo(c.connection); err != nil {
		return err
	}
	return nil
}

// OriginatorAddr - Will return originator address known as net.RemoteAddr()
// This will actually be a freeswitch address
func (c *SocketConnection) OriginatorAddr() net.Addr {
	return c.connection.RemoteAddr()
}

// ReadMessage - Will read message from channels and return them back accordingy.
// If error is received, error will be returned. If not, message will be returned back!
func (c *SocketConnection) ReadMessage() (*Message, error) {
	//logger.Debug("Waiting for connection message to be received ...")

	select {
	case err := <-c.err:
		return nil, err
	case msg := <-c.m:
		return msg, nil
	}
}

// Handle - Will handle new messages and close connection when there are no messages left to process
func (c *SocketConnection) handle() {
	for c.readOne() {
	}
	// Closing the connection now as there's nothing left to do ...
	c.Close()
}

// Close - Will close down net connection and return error if error happen
func (c *SocketConnection) Close() error {
	if err := c.connection.Close(); err != nil {
		return err
	}

	return nil
}

func (c *SocketConnection) readOne() bool {
	hdr, err := c.textreader.ReadMIMEHeader()
	if err != nil {
		c.err <- err
		return false
	}

	msg := &Message{}
	msg.Headers = make(map[string]string)
	if v := hdr.Get("Content-Length"); v != "" {
		length, err := strconv.Atoi(v)
		if err != nil {
			logger.Error(EInvalidContentLength, err)
			c.err <- err
			return false
		}
		msg.Body = make([]byte, length)
		if _, err := io.ReadFull(c.reader, msg.Body); err != nil {
			logger.Error(ECouldNotReadyBody, err)
			c.err <- err
			return false
		}
	}
	contentType := hdr.Get("Content-Type")
	if !StringInSlice(contentType, AvailableMessageTypes) {
		c.err <- fmt.Errorf(EUnsupportedMessageType, contentType, AvailableMessageTypes)
		return true
	}

	switch contentType {
	case "command/reply":
		reply := hdr.Get("Reply-Text")
		if reply[:2] == "-E" {
			c.err <- errors.New(reply[5:])
			return true
		}
		if reply[0] == '%' {
			copyHeaders(&hdr, msg, true)
		} else {
			copyHeaders(&hdr, msg, false)
		}
	case "api/response":
		if string(msg.Body[:2]) == "-E" {
			c.err <- errors.New(string(msg.Body)[5:])
			return true
		}
		copyHeaders(&hdr, msg, false)
	case "text/event-plain":
		reader := bufio.NewReader(bytes.NewReader(msg.Body))
		msg.Body = make([]byte, 0)
		textreader := textproto.NewReader(reader)
		hdr, err = textreader.ReadMIMEHeader()
		if err != nil {
			c.err <- err
			return false
		}
		if v := hdr.Get("Content-Length"); v != "" {
			length, err := strconv.Atoi(v)
			if err != nil {
				logger.Error(EInvalidContentLength, err)
				c.err <- err
				return false
			}
			msg.Body = make([]byte, length)
			if _, err = io.ReadFull(reader, msg.Body); err != nil {
				logger.Error(ECouldNotReadyBody, err)
				c.err <- err
				return false
			}
		}
		copyHeaders(&hdr, msg, true)
	case "text/event-json":
		decoded := make(map[string]interface{})
		if err := json.Unmarshal(msg.Body, &decoded); err != nil {
			logger.Error(ECouldNotUnmarshallJSON, err)
			c.err <- err
			return false
		}

		// Copy back in:
		for k, v := range decoded {
			switch v.(type) {
			case string:
				// capitalize header keys for consistency.
				msg.Headers[capitalize(k)] = v.(string)
			case int:
				msg.Headers[capitalize(k)] = strconv.Itoa(v.(int))
			default:
				logger.Warning(WRemoveNonStringProperty, k)
			}
		}
		if v, _ := msg.Headers["_body"]; v != "" {
			msg.Body = []byte(v)
			delete(msg.Headers, "_body")
		} else {
			msg.Body = []byte("")
		}
	case "text/disconnect-notice":
		copyHeaders(&hdr, msg, false)
	default:
		return true
	}

	eventName, ok := msg.Headers["Event-Name"]
	if ok && eventName == "CALL_UPDATE" {
		toPhone, ok := msg.Headers["Caller-Caller-Id-Number"]
		if !ok {
			toPhone = msg.Headers["Caller-Caller-Id-Name"]
		}
		logger.Debug("Got CALL_UPDATE for: %s", toPhone)
	}
	c.m <- msg
	return true
}

// copyHeaders copies all keys and values from the MIMEHeader to Event.Header,
// normalizing header keys to their capitalized version and values by
// unescaping them when decode is set to true.
//
// It's used after parsing plain text event headers, but not JSON.
func copyHeaders(src *textproto.MIMEHeader, dst *Message, decode bool) {
	var err error
	for k, v := range *src {
		k = capitalize(k)
		if decode {
			dst.Headers[k], err = url.QueryUnescape(v[0])
			if err != nil {
				dst.Headers[k] = v[0]
			}
		} else {
			dst.Headers[k] = v[0]
		}
	}
}

// capitalize capitalizes strings in a very particular manner.
// Headers such as Job-UUID become Job-Uuid and so on. Headers starting with
// Variable_ only replace ^v with V, and headers staring with _ are ignored.
func capitalize(s string) string {
	if s[0] == '_' {
		return s
	}
	ns := bytes.ToLower([]byte(s))
	if len(s) > 9 && s[1:9] == "ariable_" {
		ns[0] = 'V'
		return string(ns)
	}
	toUpper := true
	for n, c := range ns {
		if toUpper {
			if 'a' <= c && c <= 'z' {
				c -= 'a' - 'A'
			}
			ns[n] = c
			toUpper = false
		} else if c == '-' || c == '_' {
			toUpper = true
		}
	}
	return string(ns)
}
