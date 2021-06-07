package buffered_input

import (
	"sync"
	"time"
)

type BufferedInput struct {
	options *Options
	queue   []interface{}
	output  chan []interface{}
	closed  chan struct{}
	mutex   sync.RWMutex
}

func NewBufferedInput(options *Options) *BufferedInput {
	bi := &BufferedInput{
		options: options,
		queue:   make([]interface{}, 0, options.ChunkSize),
		output:  make(chan []interface{}, options.ChunkCount),
		closed:  make(chan struct{}, 1),
	}

	go bi.start()

	return bi
}

func (bi *BufferedInput) start() {
	for {
		select {
		case <-bi.closed:
			close(bi.closed)
			return
		case data := <-bi.output:
			bi.options.Handler(data)
		case <-time.After(bi.options.Timeout):
			bi.Flush()
		}
	}
}

func (bi *BufferedInput) Close() {
	bi.closed <- struct{}{}
	close(bi.output)
}

func (bi *BufferedInput) Push(data interface{}) {
	bi.mutex.Lock()
	bi.queue = append(bi.queue, data)
	bi.mutex.Unlock()

	// Flush immediately if the number of input exceeds queue size
	if len(bi.queue) == bi.options.ChunkSize {
		bi.Flush()
	}
}

func (bi *BufferedInput) Flush() error {

	bi.mutex.Lock()

	// Nothing's flushed
	if len(bi.queue) == 0 {
		bi.mutex.Unlock()
		return nil
	}

	// Allocate a new queue
	queue := bi.queue
	bi.queue = make([]interface{}, 0, bi.options.ChunkSize)
	bi.mutex.Unlock()

	bi.output <- queue

	return nil
}
