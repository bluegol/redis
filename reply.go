package redis

import (
	"bufio"
	"fmt"
	"strconv"

	"github.com/bluegol/errutil"
)



type Reply struct {
	Type  byte

	arr   []*Reply
	bytes []byte
}

func (r *Reply) String() (string, error) {
	if r.Type == MarkArray {
		return "", errutil.New(ErrConversion,
			"dst", "string", "src", "array")
	} else if r.Type == markNil {
		return "", errutil.New(ErrConversion,
			"dst", "string", "src", "nil")
	}

	return string(r.bytes), nil
}

func (r *Reply) Int() (int, error) {
	if r.Type == MarkArray {
		return 0, errutil.New(ErrConversion,
			"dst", "int", "src", "array")
	} else if r.Type == markNil {
		return 0, errutil.New(ErrConversion,
			"dst", "int", "src", "nil")
	}
	s := string(r.bytes)
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, errutil.Embed(ErrConversion, err,
			"dst", "int",
			"value", strconv.QuoteToASCII(s))
	}
	return i, nil
}

func (r *Reply) Int64() (int64, error) {
	if r.Type == MarkArray {
		return 0, errutil.New(ErrConversion,
			"dst", "int64", "src", "array")
	} else if r.Type == markNil {
		return 0, errutil.New(ErrConversion,
			"dst", "int64", "src", "nil")
	}
	s := string(r.bytes)
	i, err := strconv.ParseInt(string(r.bytes), 10, 64)
	if err != nil {
		return 0, errutil.Embed(ErrConversion, err,
			"dst", "int64",
			"value", strconv.QuoteToASCII(s))
	}
	return i, nil
}

func (r *Reply) Error() error {
	if r.Type == MarkError {
		errstr, _ := r.String()
		return errutil.New(ErrRedis, errutil.MoreInfo, errstr)
	}
	return nil
}

func IsNil(r []*Reply) bool {
	return len(r) == 1 && r[0].Type == markNil
}



/////////////////////////////////////////////////////////////////////

func parseReply(r *bufio.Reader) (*Reply, error) {
	t, err := r.ReadByte()
	if err != nil {
		return nil, errutil.Embed(ErrProtocol, err,
			errutil.MoreInfo, "cannot read marker.")
	}
	line, err2 := readUpToCrlf(r)
	if err2 != nil {
		return nil, errutil.Embed(ErrProtocol, err,
			errutil.MoreInfo, "cannot read marker.")
	}

	switch t {

	case MarkSimpleString, MarkError, MarkInteger:
		return &Reply{ Type: t, arr: nil, bytes: line }, nil

	case MarkBulkString:
		ll, err := strconv.Atoi(string(line))
		if err != nil {
			return nil, errutil.Embed(ErrProtocol, err,
				errutil.MoreInfo, "cannot find length",
				"read_so_far", string(line))
		}
		if ll == -1 {
			return &nilReply, nil
		} else {
			line, err := readUpToCrlf(r)
			if err != nil {
				return nil, err
			}
			return &Reply{ Type:t, arr: nil, bytes: line }, nil
		}

	case MarkArray:
		ll, err := strconv.Atoi(string(line))
		if err != nil {
			return nil, errutil.Embed(ErrProtocol, err,
				errutil.MoreInfo, "cannot find length",
				"read_so_far", string(line))
		}
		if ll == -1 {
			return &nilReply, nil
		} else {
			arr := make([]*Reply, ll)
			for i := 0; i < ll; i++ {
				arr[i], err = parseReply(r)
				if err != nil {
					return nil, errutil.Embed( ErrProtocol, err,
						errutil.MoreInfo,
							fmt.Sprintf("cannot read %d-th element of array", i ) )
				}
			}
			return &Reply{ Type:t, arr: arr, bytes: nil }, nil
		}

	default:
		return nil, errutil.New( ErrProtocol,
			errutil.MoreInfo,
				fmt.Sprintf("unknown mark %c, read so far: %v", t, line) )
	}
}

var nilReply Reply = Reply{ Type: markNil }

func readUpToCrlf(r *bufio.Reader) ([]byte, error) {
	line, err := r.ReadBytes(byteCr)
	if err != nil {
		return nil, errutil.Embed(ErrProtocol, err,
			errutil.MoreInfo, "cannot find cr",
			"read_so_far", string(line))
	}
	c, err := r.ReadByte()
	if err != nil {
		return nil, errutil.Embed(ErrProtocol, err,
			errutil.MoreInfo, "cannot read lf after cr",
			"read_so_far", string(line))
	} else if c != byteLf {
		return nil, errutil.New(ErrShouldNotHappen,
			errutil.MoreInfo, "lf doesn't follow cr",
			"read_so_far", string(line), "read", string(c) )
	}

	return line[0:len(line)-1], nil
}
