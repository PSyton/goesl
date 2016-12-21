// Copyright 2015 Nevio Vesic
// Please check out LICENSE file for more information about what you CAN and what you CANNOT do!
// Basically in short this is a free software for you to do whatever you want to do BUT copyright must be included!
// I didn't write all of this code so you could say it's yours.
// MIT License

package goesl

import (
	"fmt"
	"sort"
)

// Message - Freeswitch Message that is received by GoESL. Message struct is here to help with parsing message
// and dumping its contents. In addition to that it's here to make sure received message is in fact message we wish/can support
type Message struct {
	Headers map[string]string
	Body    []byte
}

// String - Will return message representation as string
func (m *Message) String() string {
	return fmt.Sprintf("%v body=%s", m.Headers, m.Body)
}

// GetCallUUID - Will return Caller-Unique-Id
func (m *Message) GetCallUUID() string {
	return m.GetHeader("Caller-Unique-Id")
}

// GetHeader - Will return message header value, or "" if the key is not set.
func (m *Message) GetHeader(key string) string {
	return m.Headers[key]
}

// Dump - Will return message prepared to be dumped out. It's like prettify message for output
func (m *Message) Dump() (resp string) {
	var keys []string

	for k := range m.Headers {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		resp += fmt.Sprintf("%s: %s\r\n", k, m.Headers[k])
	}

	resp += fmt.Sprintf("BODY: %v\r\n", string(m.Body))

	return
}
