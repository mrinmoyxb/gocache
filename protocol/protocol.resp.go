package protocol

import (
	"bufio"
	"fmt"
	"io"
)

// RESP: REDIS Serialization Protocol
type RESPValue struct {
	Type    byte
	String  string
	Integer int64
	Array   []RESPValue
}

type Reader struct {
	reader *bufio.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		reader: bufio.NewReader(r),
	}
}

func (r *Reader) Read() (RESPValue, error) {
	typeByte, err := r.reader.ReadByte()
	if err != nil {
		return RESPValue{}, err
	}

	switch typeByte {
	case '+':
		return r.readSimpleString()
	case '-':
		return r.readError()
	case ':':
		return r.readInteger()
	case '$':
		return r.readBulkString()
	case '*':
	return r.readArray()
	default:
		return RESPValue{}, fmt.Errorf("unknown RESP type: %q", typeByte)
	}
}

func (r *Reader) readSimpleString() (RESPValue, error) {
	line, err := r.reader.ReadString('\n')
	if err != nil {
		return RESPValue{}, err
	}

	line = line[:len(line)-2]

	return RESPValue{
		Type:   '+',
		String: line,
	}, nil
}

func (r *Reader) readError() (RESPValue, error) {
	line, err := r.reader.ReadString('\n')
	if err != nil {
		return RESPValue{}, err
	}

	line = line[:len(line)-2]

	return RESPValue{
		Type:   '-',
		String: line,
	}, nil
}

func (r *Reader) readInteger() (RESPValue, error) {
	line, err := r.reader.ReadString('\n')
	if err != nil {
		return RESPValue{}, err
	}

	line = line[:len(line)-2]

	var value int64

	_, err = fmt.Sscanf(line, "%d", &value)
	if err != nil {
		return RESPValue{}, err
	}

	return RESPValue{
		Type:    ':',
		Integer: value,
	}, nil
}

func (r *Reader) readBulkString() (RESPValue, error) {
	line, err := r.reader.ReadString('\n')
	if err != nil {
		return RESPValue{}, err
	}

	line = line[:len(line)-2]

	var length int

	_, err = fmt.Sscanf(line, "%d", &length)
	if err != nil {
		return RESPValue{}, err
	}

	data := make([]byte, length)

	_, err = io.ReadFull(r.reader, data)
	if err != nil {
		return RESPValue{}, err
	}

	_, err = r.reader.ReadString('\n')
	if err != nil {
		return RESPValue{}, err
	}

	return RESPValue{
		Type:   '$',
		String: string(data),
	}, nil
}

func (r *Reader) readArray() (RESPValue, error) {
	line, err := r.reader.ReadString('\n')
	if err != nil {
		return RESPValue{}, err
	}

	line = line[:len(line)-2]

	var length int 
	_, err = fmt.Sscanf(line, "%d", &length)
	if err != nil {
		return RESPValue{}, err
	}

	values := make([]RESPValue, length)
	for i:=0; i<length; i++ {
		value, err := r.Read()
		if err != nil {
			return RESPValue{}, err
		}
		values[i] = value
	}


	return RESPValue{
		Type: '*',
		Array: values,
	}, nil
}
