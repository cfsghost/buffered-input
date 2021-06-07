package buffered_input

import "time"

type DataHandlerFunc func([]interface{})

type Options struct {
	ChunkSize  int
	ChunkCount int
	Timeout    time.Duration
	Handler    DataHandlerFunc
}

func NewOptions() *Options {
	return &Options{
		ChunkSize:  1000,
		ChunkCount: 1000,
		Timeout:    50 * time.Millisecond,
		Handler:    func(chunk []interface{}) {},
	}
}
