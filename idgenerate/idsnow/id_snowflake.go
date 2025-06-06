package idsnow

import (
	"errors"
	"net"
	"sync"
	"time"

	"github.com/illidaris/aphrodite/idgenerate/dep"
)

type Settings struct {
	BitsSequence   int
	BitsMachineID  int
	TimeUnit       time.Duration
	StartTime      time.Time
	MachineID      func() (int, error) // 改造成实时计算(入参基因1（0~31），基因2（0~3）) 1-默认机器IP组成 2-自定义: NodeId（节点Id） 2^8 FrameId（主备帧数） 2^2 Gene2（基因2位取模） 2^2 Gene（基因4位取模） 2^4
	CheckMachineID func(int) bool
}

type Sonyflake struct {
	mutex *sync.Mutex

	bitsTime     int
	bitsSequence int
	bitsMachine  int

	timeUnit    int64
	startTime   int64
	elapsedTime int64

	sequence int
	machine  int
}

var (
	ErrInvalidBitsTime      = errors.New("bit length for time must be 32 or more")
	ErrInvalidBitsSequence  = errors.New("invalid bit length for sequence number")
	ErrInvalidBitsMachineID = errors.New("invalid bit length for machine id")
	ErrInvalidTimeUnit      = errors.New("invalid time unit")
	ErrInvalidSequence      = errors.New("invalid sequence number")
	ErrInvalidMachineID     = errors.New("invalid machine id")
	ErrStartTimeAhead       = errors.New("start time is ahead")
	ErrOverTimeLimit        = errors.New("over the time limit")
	ErrNoPrivateAddress     = errors.New("no private ip address")
)

const (
	defaultTimeUnit = 1e7 // nsec, i.e. 10 msec

	defaultBitsTime     = 39
	defaultBitsSequence = 8
	defaultBitsMachine  = 16
)

var defaultInterfaceAddrs = net.InterfaceAddrs

func New(st Settings) (*Sonyflake, error) {
	if st.BitsSequence < 0 || st.BitsSequence > 30 {
		return nil, ErrInvalidBitsSequence
	}
	if st.BitsMachineID < 0 || st.BitsMachineID > 30 {
		return nil, ErrInvalidBitsMachineID
	}
	if st.TimeUnit < 0 || (st.TimeUnit > 0 && st.TimeUnit < time.Millisecond) {
		return nil, ErrInvalidTimeUnit
	}
	if st.StartTime.After(time.Now()) {
		return nil, ErrStartTimeAhead
	}

	sf := new(Sonyflake)
	sf.mutex = new(sync.Mutex)

	if st.BitsSequence == 0 {
		sf.bitsSequence = defaultBitsSequence
	} else {
		sf.bitsSequence = st.BitsSequence
	}

	if st.BitsMachineID == 0 {
		sf.bitsMachine = defaultBitsMachine
	} else {
		sf.bitsMachine = st.BitsMachineID
	}

	sf.bitsTime = 63 - sf.bitsSequence - sf.bitsMachine
	if sf.bitsTime < 32 {
		return nil, ErrInvalidBitsTime
	}

	if st.TimeUnit == 0 {
		sf.timeUnit = defaultTimeUnit
	} else {
		sf.timeUnit = int64(st.TimeUnit)
	}

	if st.StartTime.IsZero() {
		sf.startTime = sf.toInternalTime(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
	} else {
		sf.startTime = sf.toInternalTime(st.StartTime)
	}

	sf.sequence = 1<<sf.bitsSequence - 1

	var err error
	if st.MachineID == nil {
		sf.machine, err = lower16BitPrivateIP(defaultInterfaceAddrs)
	} else {
		sf.machine, err = st.MachineID()
	}
	if err != nil {
		return nil, err
	}

	if sf.machine < 0 || sf.machine >= 1<<sf.bitsMachine {
		return nil, ErrInvalidMachineID
	}

	if st.CheckMachineID != nil && !st.CheckMachineID(sf.machine) {
		return nil, ErrInvalidMachineID
	}

	return sf, nil
}

// NextID generates a next unique ID as int64.
// After the Sonyflake time overflows, NextID returns an error.
func (sf *Sonyflake) NextID() (int64, error) {
	maskSequence := 1<<sf.bitsSequence - 1

	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	current := sf.currentElapsedTime()
	if sf.elapsedTime < current {
		sf.elapsedTime = current
		sf.sequence = 0
	} else {
		sf.sequence = (sf.sequence + 1) & maskSequence
		if sf.sequence == 0 {
			sf.elapsedTime++
			overtime := sf.elapsedTime - current
			sf.sleep(overtime)
		}
	}

	return sf.toID()
}

func (sf *Sonyflake) toInternalTime(t time.Time) int64 {
	return t.UTC().UnixNano() / sf.timeUnit
}

func (sf *Sonyflake) currentElapsedTime() int64 {
	return sf.toInternalTime(time.Now()) - sf.startTime
}

func (sf *Sonyflake) sleep(overtime int64) {
	sleepTime := time.Duration(overtime*sf.timeUnit) -
		time.Duration(time.Now().UTC().UnixNano()%sf.timeUnit)
	time.Sleep(sleepTime)
}

func (sf *Sonyflake) toID() (int64, error) {
	if sf.elapsedTime >= 1<<sf.bitsTime {
		return 0, ErrOverTimeLimit
	}

	return sf.elapsedTime<<(sf.bitsSequence+sf.bitsMachine) |
		int64(sf.sequence)<<sf.bitsMachine |
		int64(sf.machine), nil
}

func privateIPv4(interfaceAddrs dep.InterfaceAddrs) (net.IP, error) {
	as, err := interfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, a := range as {
		ipnet, ok := a.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}

		ip := ipnet.IP.To4()
		if isPrivateIPv4(ip) {
			return ip, nil
		}
	}
	return nil, ErrNoPrivateAddress
}

func isPrivateIPv4(ip net.IP) bool {
	// Allow private IP addresses (RFC1918) and link-local addresses (RFC3927)
	return ip != nil &&
		(ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168 || ip[0] == 169 && ip[1] == 254)
}

func lower16BitPrivateIP(interfaceAddrs dep.InterfaceAddrs) (int, error) {
	ip, err := privateIPv4(interfaceAddrs)
	if err != nil {
		return 0, err
	}

	return int(ip[2])<<8 + int(ip[3]), nil
}

// ToTime returns the time when the given ID was generated.
func (sf *Sonyflake) ToTime(id int64) time.Time {
	return time.Unix(0, (sf.startTime+sf.timePart(id))*sf.timeUnit)
}

// Compose creates a Sonyflake ID from its components.
// The time parameter should be the time when the ID was generated.
// The sequence parameter should be between 0 and 2^BitsSequence-1 (inclusive).
// The machineID parameter should be between 0 and 2^BitsMachineID-1 (inclusive).
func (sf *Sonyflake) Compose(t time.Time, sequence, machineID int) (int64, error) {
	elapsedTime := sf.toInternalTime(t.UTC()) - sf.startTime
	if elapsedTime < 0 {
		return 0, ErrStartTimeAhead
	}
	if elapsedTime >= 1<<sf.bitsTime {
		return 0, ErrOverTimeLimit
	}

	if sequence < 0 || sequence >= 1<<sf.bitsSequence {
		return 0, ErrInvalidSequence
	}

	if machineID < 0 || machineID >= 1<<sf.bitsMachine {
		return 0, ErrInvalidMachineID
	}

	return elapsedTime<<(sf.bitsSequence+sf.bitsMachine) |
		int64(sequence)<<sf.bitsMachine |
		int64(machineID), nil
}

// Decompose returns a set of Sonyflake ID parts.
func (sf *Sonyflake) Decompose(id int64) map[string]int64 {
	time := sf.timePart(id)
	sequence := sf.sequencePart(id)
	machine := sf.machinePart(id)
	return map[string]int64{
		"id":       id,
		"time":     time,
		"sequence": sequence,
		"machine":  machine,
	}
}

func (sf *Sonyflake) timePart(id int64) int64 {
	return id >> (sf.bitsSequence + sf.bitsMachine)
}

func (sf *Sonyflake) sequencePart(id int64) int64 {
	maskSequence := int64((1<<sf.bitsSequence - 1) << sf.bitsMachine)
	return (id & maskSequence) >> sf.bitsMachine
}

func (sf *Sonyflake) machinePart(id int64) int64 {
	maskMachine := int64(1<<sf.bitsMachine - 1)
	return id & maskMachine
}
