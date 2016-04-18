package redis

import "errors"

var (
	// ErrProtocol indicates unexpected protocol error
	ErrProtocol        error

	// ErrRedis indicates redis returned ERR
	ErrRedis           error

	// ErrRentTimeout indicates rented client is not returned in time.
	ErrRentTimeout           error
	ErrQuitWhileRented       error
	ErrRentedClientCannotQuit error
	InfoReturnedAfterTimeout error

	// user input errors
	ErrConversion     error
	ErrInvalidCmd     error
	ErrInvalidArgType error

	// network
	ErrCannotConnect   error
	ErrBeforeSending   error
	ErrWhileSending    error
	ErrCannotRead      error
	InfoConnected      error

	// about quitting
	ErrAlreadyClosed   error
	InfoAlreadyQuit    error
	InfoQuitStarted    error
	InfoQuitDone       error
)

/////////////////////////////////////////////////////////////////////

func init() {
	ErrProtocol = errors.New("cannot decode according to protocol.")

	ErrRedis = errors.New("redis returned ERR")

	ErrRentTimeout = errors.New("rent timeout")
	ErrQuitWhileRented = errors.New("quit signaled while rented")
	ErrRentedClientCannotQuit = errors.New("should not signal quit to rented client")
	InfoReturnedAfterTimeout = errors.New("returned after timeout")

	ErrConversion = errors.New("cannot convert.")
	ErrInvalidCmd = errors.New("command is invalid.")
	ErrInvalidArgType = errors.New("type is invalid.")

	ErrCannotConnect = errors.New("cannot connect to server.")
	ErrBeforeSending = errors.New("error before sending message")
	ErrWhileSending = errors.New("error in the middle of sending message.")
	ErrCannotRead = errors.New("error while receiving message.")
	InfoConnected = errors.New("connected")

	ErrAlreadyClosed = errors.New("client is already closed")
	InfoAlreadyQuit = errors.New("client is quitting/has quit already")
	InfoQuitStarted = errors.New("quit started")
	InfoQuitDone = errors.New("quit done")
}
