// Copyright 2015 Nevio Vesic
// Please check out LICENSE file for more information about what you CAN and what you CANNOT do!
// Basically in short this is a free software for you to do whatever you want to do BUT copyright must be included!
// I didn't write all of this code so you could say it's yours.
// MIT License

package goesl

import (
	"fmt"
)

const (
	eInvalidCommand             = "Invalid command provided. Command cannot contain \\r and/or \\n. Provided command is: %s"
	eCouldNotReadMIMEHeaders    = "Error while reading MIME headers: %s"
	eInvalidContentLength       = "Unable to get size of content-length: %s"
	eUnsuccessfulReply          = "Got error while reading from reply command: %s"
	eCouldNotReadBody           = "Got error while reading reader body: %s"
	eUnsupportedMessageType     = "Unsupported message type! We got '%s'. Supported types are: %v "
	eUnsupportedMessageTypeLite = "Unsupported message type! We got '%s'"
	eUnmarshallJSON             = "Error while unmarshal JSON event: %s"
	eCouldNotStartListener      = "Got error while attempting to start listener: %s"
	eListenerConnection         = "Listener connection error: %s"
	invalidServerAddr           = "Please make sure to pass along valid address. You've passed: \"%s\""
	unexpectedAuthHeader        = "Expected auth/request content type. Got %s"
	invalidPassword             = "Could not authenticate against freeswitch with provided password."
	wRemoveNonStringProperty    = "Removed non-string property (%s)"
	errorWhileAccepConnection   = "Got error while accepting connection: %s"
)

type errorImpl struct {
	message string
}

func (e *errorImpl) Error() string {
	return e.message
}

func newError(aMsg string) errorImpl {
	return errorImpl{
		message: aMsg,
	}
}

// ErrorInvalidCommand fired when try to send invalid command
type ErrorInvalidCommand struct {
	errorImpl
}

func newErrorInvalidCommand(aData interface{}) *ErrorInvalidCommand {
	return &ErrorInvalidCommand{
		errorImpl: newError(fmt.Sprint(eInvalidCommand, aData)),
	}
}

// ErrorReadMIMEHeaders ...
type ErrorReadMIMEHeaders struct {
	errorImpl
}

func newErrorReadMIMEHeaders(aError error) *ErrorReadMIMEHeaders {
	return &ErrorReadMIMEHeaders{
		errorImpl: newError(fmt.Sprintf(eCouldNotReadMIMEHeaders, aError.Error())),
	}
}

// ErrorInvalidContentLength ...
type ErrorInvalidContentLength struct {
	errorImpl
}

func newErrorInvalidContentLength(aError error) *ErrorInvalidContentLength {
	return &ErrorInvalidContentLength{
		errorImpl: newError(fmt.Sprintf(eInvalidContentLength, aError.Error())),
	}
}

// ErrorUnsuccessfulReply ...
type ErrorUnsuccessfulReply struct {
	errorImpl
}

func newErrorUnsuccessfulReply(aReply string) *ErrorUnsuccessfulReply {
	return &ErrorUnsuccessfulReply{
		errorImpl: newError(fmt.Sprintf(eUnsuccessfulReply, aReply)),
	}
}

// ErrorCouldNotReadBody ...
type ErrorCouldNotReadBody struct {
	errorImpl
}

func newErrorCouldNotReadBody(aError error) *ErrorCouldNotReadBody {
	return &ErrorCouldNotReadBody{
		errorImpl: newError(fmt.Sprintf(eCouldNotReadBody, aError.Error())),
	}
}

// ErrorUnsupportedMessageType ...
type ErrorUnsupportedMessageType struct {
	errorImpl
}

func newErrorUnsupportedMessageType(aType string) *ErrorUnsuccessfulReply {
	return &ErrorUnsuccessfulReply{
		errorImpl: newError(fmt.Sprintf(eUnsupportedMessageTypeLite, aType)),
	}
}

// ErrorInvalidServerAddr ...
type ErrorInvalidServerAddr struct {
	errorImpl
}

func newErrorInvalidServerAddr(aAddress string) *ErrorInvalidServerAddr {
	return &ErrorInvalidServerAddr{
		errorImpl: newError(fmt.Sprintf(invalidServerAddr, aAddress)),
	}
}

// ErrorUnexpectedAuthHeader ...
type ErrorUnexpectedAuthHeader struct {
	errorImpl
}

func newErrorUnexpectedAuthHeader(aCType string) *ErrorUnexpectedAuthHeader {
	return &ErrorUnexpectedAuthHeader{
		errorImpl: newError(fmt.Sprintf(unexpectedAuthHeader, aCType)),
	}
}

// ErrorInvalidPassword ...
type ErrorInvalidPassword struct {
	errorImpl
}

func newErrorInvalidPassword() *ErrorInvalidPassword {
	return &ErrorInvalidPassword{
		errorImpl: newError(invalidPassword),
	}
}

// ErrorUnmarshallJSON ...
type ErrorUnmarshallJSON struct {
	errorImpl
}

func newErrorUnmarshallJSON(aError error) *ErrorUnmarshallJSON {
	return &ErrorUnmarshallJSON{
		errorImpl: newError(fmt.Sprintf(eUnmarshallJSON, aError.Error())),
	}
}

// ErrorSendEvent ...
type ErrorSendEvent struct {
	errorImpl
}

func newErrorSendEvent(aLen int) *ErrorSendEvent {
	return &ErrorSendEvent{
		errorImpl: newError(fmt.Sprintf("Must send at least one event header, detected `%d` header", aLen)),
	}
}
