package embedded

import (
	"math/rand"
	"time"
)

type IItem any // client or sdk

// IComponent
type IComponent[T IItem] interface {
	NewWriter(id string, items ...T)
	GetWriter(id string) T
	NewReader(id string, items ...T)
	GetReader(id string) T
	SetWriterBalance(f func(ts ...IInstance[T]) IInstance[T])
	SetReaderBalance(f func(ts ...IInstance[T]) IInstance[T])
}

// IInstance
type IInstance[T IItem] interface {
	GetID() string
	GetWeight() float64
	GetValue() T
}

var _ = IComponent[IItem](&Component[IItem]{})           // check
var _ = IInstance[IItem](&Instance[IItem]{})             // check
var rd = rand.New(rand.NewSource(time.Now().UnixNano())) // rand seed

func defaultBalance[T IItem](ts ...IInstance[T]) IInstance[T] {
	switch len(ts) {
	case 0:
		return nil
	case 1:
		return ts[0]
	default:
		index := rd.Intn(len(ts))
		return ts[index]
	}
}

func NewComponent[T IItem]() IComponent[T] {
	return &Component[T]{
		Writers:       map[string][]IInstance[T]{},
		WriterBalance: defaultBalance[T],
		Readers:       map[string][]IInstance[T]{},
		ReaderBalance: defaultBalance[T],
	}
}

type Component[T IItem] struct {
	Writers       map[string][]IInstance[T]             // 写节点
	WriterBalance func(ts ...IInstance[T]) IInstance[T] // 选取算法
	Readers       map[string][]IInstance[T]             // 读节点
	ReaderBalance func(ts ...IInstance[T]) IInstance[T] //选取算法
}

func (c *Component[T]) SetWriterBalance(f func(ts ...IInstance[T]) IInstance[T]) {
	c.WriterBalance = f
}

func (c *Component[T]) SetReaderBalance(f func(ts ...IInstance[T]) IInstance[T]) {
	c.ReaderBalance = f
}

func (c *Component[T]) NewWriter(id string, items ...T) {
	for _, item := range items {
		c.Writers[id] = append(c.Writers[id], NewInstance(id, item))
	}
}

func (c *Component[T]) GetWriter(id string) T {
	return c.WriterBalance(c.Writers[id]...).GetValue()
}

func (c *Component[T]) NewReader(id string, items ...T) {
	for _, item := range items {
		c.Readers[id] = append(c.Readers[id], NewInstance(id, item))
	}
}

func (c *Component[T]) GetReader(id string) T {
	return c.ReaderBalance(c.Readers[id]...).GetValue()
}
func NewInstance[T IItem](id string, item T) IInstance[T] {
	return &Instance[T]{
		Id:    id,
		Value: item,
	}
}

type Instance[T IItem] struct {
	Id     string
	Weight float64
	Value  T
}

func (c *Instance[T]) GetID() string {
	return c.Id
}

func (c *Instance[T]) GetWeight() float64 {
	return c.Weight
}

func (c *Instance[T]) GetValue() T {
	return c.Value
}
