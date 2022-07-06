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

	done := make(chan time.Time)

	options := NewOptions()
	options.Timeout = time.Second * 2
	options.Handler = func(chunk []interface{}) {
		if len(chunk) == 1 {
			done <- chunk[0].(time.Time)
			return
		}

		done <- time.Now()
	}

	bi := NewBufferedInput(options)
	defer bi.Close()

	bi.Push(time.Now())

	mt := <-done
	if !mt.Add(time.Second * 2).Before(time.Now()) {
		t.Fail()
	}
}
