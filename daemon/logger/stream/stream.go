package stream

import (
	"github.com/docker/docker/daemon/logger"
)

// RingBuffer as per http://pivotallabs.com/a-concurrent-ring-buffer-for-go/
type RingBuffer struct {
	in  chan *logger.Message
	out chan *logger.Message
}

func NewRingBuffer(in chan *logger.Message, out chan *logger.Message) *RingBuffer {
	return &RingBuffer{in, out}
}

func (r *RingBuffer) Run() {
	for v := range r.in {
		select {
		case r.out <- v:
		default:
			<-r.out
			r.out <- v
		}
	}
	close(r.out)
}

type Stream struct {
	r *RingBuffer
}

func New() (logger.Logger, error) {
	in := make(chan *logger.Message)
	out := make(chan *logger.Message, 100)
	rb := NewRingBuffer(in, out)
	go rb.Run()
	return &Stream{
		r: rb,
	}, nil
}

func (s *Stream) Log(msg *logger.Message) error {
	s.r.in <- msg
	return nil
}

func (s *Stream) Close() error {
	close(s.r.in)
	return nil
}

func (s *Stream) Name() string {
	return "Stream"
}

func (s *Stream) Output() <-chan *logger.Message {
	return s.r.out
}
