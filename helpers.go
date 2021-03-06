// Copyright 2015 Nevio Vesic
// Please check out LICENSE file for more information about what you CAN and what you CANNOT do!
// Basically in short this is a free software for you to do whatever you want to do BUT copyright must be included!
// I didn't write all of this code so you could say it's yours.
// MIT License

package goesl

import "fmt"

// ExecuteSet - Helper that you can use to execute SET application against active ESL session
func (sc *SocketConnection) ExecuteSet(key string, value string, sync bool) error {
	return sc.Execute("set", key+"="+value, sync)
}

// ExecuteAnswer - Helper desgned to help with executing Answer against active ESL session
func (sc *SocketConnection) ExecuteAnswer(args string, sync bool) (err error) {
	return sc.Execute("answer", args, sync)
}

// ExecuteHangup - Helper desgned to help with executing Hangup against active ESL session
func (sc *SocketConnection) ExecuteHangup(uuid string, args string, sync bool) (err error) {
	if uuid != "" {
		return sc.ExecuteUUID(uuid, "hangup", args, sync)
	}

	return sc.Execute("hangup", args, sync)
}

// Api - Helper designed to attach api in front of the command so that you do not need to write it
func (sc *SocketConnection) Api(command string) error {
	return sc.Send(fmt.Sprintf("api " + command))
}

// BgApi - Helper designed to attach bgapi in front of the command so that you do not need to write it
func (sc *SocketConnection) BgApi(command string) error {
	return sc.Send(fmt.Sprintf("bgapi " + command))
}
