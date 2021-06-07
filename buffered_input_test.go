package buffered_input

import (
	"testing"
	"time"
)

func TestChunkSize(t *testing.T) {

	done := make(chan bool)

	options := NewOptions()
	options.ChunkSize = 100
	options.Handler = func(chunk []interface{}) {
		if len(chunk) == 100 {
			done <- true
			return
		}

		done <- false
	}

	bi := NewBufferedInput(options)
	defer bi.Close()

	for i := 1; i <= 100; i++ {
		bi.Push(i)
	}

	success := <-done
	if !success {
		t.Fail()
	}
}

func TestTimeout(t *testing.T) {

	done := make(chan bool)

	options := NewOptions()
	options.Timeout = time.Second
	options.Handler = func(chunk []interface{}) {
		if len(chunk) == 1 {
			done <- true
			return
		}

		done <- false
	}

	bi := NewBufferedInput(options)
	defer bi.Close()

	bi.Push(true)

	success := <-done
	if !success {
		t.Fail()
	}
}
