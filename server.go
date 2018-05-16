// Copyright 2015 Nevio Vesic
// Please check out LICENSE file for more information about what you CAN and what you CANNOT do!
// Basically in short this is a free software for you to do whatever you want to do BUT copyright must be included!
// I didn't write all of this code so you could say it's yours.
// MIT License

package goesl

import (
	"net"
)

type (
	// HandlerFunc hadler for incomming connection
	HandlerFunc func(*ESLConnection) bool
)

// ESLConnection wrapper for incoming connection
type ESLConnection struct {
	*SocketConnection
}

func (c *ESLConnection) process(aHandler HandlerFunc) {
	logger.Debug("Got new connection from: %s", c.OriginatorAddr())
	defer logger.Debug("Finish connection handler")

	if err := c.connect(); err != nil {
		logger.Error(errorWhileAccepConnection, err)
		c.Close()
		return
	}

	// process events fron Freeswitch
	go c.handle()

	shouldExit := aHandler(c)
	if shouldExit {
		// Close connection
		c.exit()
	} else {
		c.Close()
	}

}

// Connect - Helper designed to help you handle connection. Each outbound server when handling needs to connect e.g. accept
// connection in order for you to do answer, hangup or do whatever else you wish to do
func (c *ESLConnection) connect() error {
	return c.Send("connect")
}

// Exit - Used to send exit signal to ESL. It will basically hangup call and close connection
func (c *ESLConnection) exit() error {
	return c.Send("exit")
}

// ESLServer - In case you need to start server, this Struct have it covered
type ESLServer struct {
	listener net.Listener
	stop     chan struct{}
}

// Start - Will start new outbound server
func (s *ESLServer) Start(aListenAddress string, aHandler HandlerFunc) error {
	logger.Info("Starting Freeswitch Outbound Server @ (address: %s) ...", aListenAddress)

	var err error

	s.listener, err = net.Listen("tcp", aListenAddress)

	if err != nil {
		logger.Error(eCouldNotStartListener, err)
		return err
	}

	go s.runServer(aHandler)

	return err
}

func (s *ESLServer) runServer(aHandler HandlerFunc) {
	for {
		logger.Debug("Waiting for incoming connections ...")

		c, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.stop:
			default:
				logger.Error(eListenerConnection, err)
			}
			return
		}

		conn := ESLConnection{
			SocketConnection: newConnection(c),
		}

		go conn.process(aHandler)
	}
}

// Stop - Will close server connection once SIGTERM/Interrupt is received
func (s *ESLServer) Stop() {
	logger.Debug("Stopping Outbound Server ...")
	close(s.stop)
	s.listener.Close()
}

// NewESLServer - Will instanciate new outbound server
func NewESLServer() *ESLServer {
	return &ESLServer{
		stop: make(chan struct{}),
	}
}
