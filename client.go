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

	Proto   string `json:"freeswitch_protocol"`
	Addr    string `json:"freeswitch_addr"`
	Passwd  string `json:"freeswitch_password"`
	Timeout int    `json:"freeswitch_connection_timeout"`
}

// EstablishConnection - Will attempt to establish connection against freeswitch and create new SocketConnection
func (c *Client) establishConnection() (err error) {
	c.SocketConnection, err = dial(c.Proto, c.Addr, time.Duration(c.Timeout*int(time.Second)))
	return err
}

func (c *Client) authenticate() error {
	m, err := c.textreader.ReadMIMEHeader()
	if err != nil && err.Error() != "EOF" {
		return newErrorReadMIMEHeaders(err)
	}

	cType := m.Get("Content-Type")
	if cType != "auth/request" {
		logger.Error(unexpectedAuthHeader, cType)
		return newErrorUnexpectedAuthHeader(cType)
	}

	err = c.Send("auth " + c.Passwd)
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
func NewClient(host string, port uint, passwd string, timeout int) (Client, error) {
	client := Client{
		Proto:   "tcp", // Let me know if you ever need this open up lol
		Addr:    net.JoinHostPort(host, strconv.Itoa(int(port))),
		Passwd:  passwd,
		Timeout: timeout,
	}

	err := client.establishConnection()

	if err == nil {
		err = client.authenticate()

		if err != nil {
			client.Close()
		} else {
			go client.handle()
		}
	}

	return client, nil
}

// Errors - returns error channel
func (c *Client) Errors() chan error {
	return c.err
}

// Messages - returns messages channel
func (c *Client) Messages() chan *Message {
	return c.m
}
