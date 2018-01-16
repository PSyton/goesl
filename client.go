// Copyright 2015 Nevio Vesic
// Please check out LICENSE file for more information about what you CAN and what you CANNOT do!
// Basically in short this is a free software for you to do whatever you want to do BUT copyright must be included!
// I didn't write all of this code so you could say it's yours.
// MIT License

package goesl

import (
	"net"
	"strconv"
	"time"
)

// Client - In case you need to do inbound dialing against freeswitch server in order to originate call or see
// sofia statuses or whatever else you came up with
type Client struct {
	*SocketConnection
}

func (c *Client) authenticate(password string) error {
	m, err := c.textreader.ReadMIMEHeader()
	if err != nil && err.Error() != "EOF" {
		return newErrorReadMIMEHeaders(err)
	}

	cType := m.Get("Content-Type")
	if cType != "auth/request" {
		logger.Error(unexpectedAuthHeader, cType)
		return newErrorUnexpectedAuthHeader(cType)
	}

	err = c.Send("auth " + password)
	if err != nil {
		return err
	}

	m, err = c.textreader.ReadMIMEHeader()
	if err != nil && err.Error() != "EOF" {
		return newErrorReadMIMEHeaders(err)
	}

	if m.Get("Reply-Text") != "+OK accepted" {
		logger.Error(invalidPassword)
		return newErrorInvalidPassword()
	}

	return nil
}

// NewClient - Will initiate new client that will establish connection and attempt to authenticate
// against connected freeswitch server
func NewClient(host string, port uint, passwd string, timeout int) (*Client, error) {

	address := net.JoinHostPort(host, strconv.Itoa(int(port)))
	socketConnection, err := dial("tcp", address, time.Duration(timeout*int(time.Second)))

	if err != nil {
		return nil, err
	}

	client := &Client{
		SocketConnection: socketConnection,
	}

	err = client.authenticate(passwd)
	if err != nil {
		client.Close()
		return nil, err
	}

	go client.handle()
	return client, nil
}
