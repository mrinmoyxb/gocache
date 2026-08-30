package protocol

type Response struct {
	Type    byte
	String  string
	Integer int64
}

func SimpleString(value string) Response {
	return Response{
		Type:   '+',
		String: value,
	}
}

func Error(message string) Response {
	return Response{
		Type:   '-',
		String: message,
	}
}

func Integer(value int64) Response {
	return Response{
		Type:    ':',
		Integer: value,
	}
}

func BulkString(value string) Response {
	return Response{
		Type:   '$',
		String: value,
	}
}

func Null() Response {
	return Response{
		Type:   '$',
		String: "",
	}
}
