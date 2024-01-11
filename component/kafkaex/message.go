package kafkaex

type Message struct {
	Id         string `json:"id"`         // identify id
	Topic      string `json:"topic"`      // topic
	Partition  int32  `json:"partition"`  // partition
	ConsumerId string `json:"consumerId"` // consumerId
	Offset     int64  `json:"offset"`     // offset
	Headers    string `json:"headers"`    // headers
	Key        []byte `json:"key"`        // key
	Value      []byte `json:"value"`      // value
	Ts         int64  `json:"ts"`         // ts
	BlockTs    int64  `json:"blockts"`    // blockts
}
