package protocol

// RESP: REDIS Serialization Protocol
type RESPValue struct {
	Type    byte
	String  string
	Integer int64
	Array   []RESPValue
}