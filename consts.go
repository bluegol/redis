package redis

const (
	MarkSimpleString byte = byte('+')
	MarkError = byte('-')
	MarkInteger = byte(':')
	MarkBulkString = byte('$')
	MarkArray = byte('*')

	// markNil is not in redis protocol, but used internally in this package
	markNil = byte('N')
)




////////////////////////////////////////////////////////////////////

const (
	byteCr = byte('\r')
	byteLf = byte('\n')
)

var bytesMarkArray = []byte{ MarkArray }
var bytesMarkBulkString = []byte{ MarkBulkString }
var bytesCrlf []byte = []byte{ byteCr, byteLf }
