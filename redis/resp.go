package redis

type RESP struct {
	*RESPSerializer
	*RESPDeserializer
}

func NewRESP(data []byte) *RESP {
	return &RESP{RESPSerializer: NewRespSerializer(), RESPDeserializer: NewRESPDeserializer(data)}
}
