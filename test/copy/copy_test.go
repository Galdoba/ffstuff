package copy

import (
	"os"
	"fmt"
	"io"
	"time"

	"testing"
)

func TestCopy(t *testing.T) {
	in, err := os.Open(`\\192.168.31.4\root\EDIT\_amedia\Hily_s01\Hily_s01e01_HD.mp4`)
	if err != nil {
		t.Error(err)
		return
	}
	out, err := os.Create("test.mkv")
	if err != nil {
		t.Error(err)
		return
	}
	start := time.Now()
	err = Copy(t, out, in)
	t.Logf(" --------: %v\n", time.Now().Sub(start))
	if err != nil {
		t.Error(err)
		return
	}
}

func Copy(t *testing.T, out io.Writer, in io.Reader) error {
	const(
		minSize = 1024*1024*10
	 	bufSize = minSize
		numBufs = 8
	)
	type Chunk struct {
		buf [bufSize]byte
		len int
	}
	//var bufs [numBufs]Chunk

	startTime0 := time.Now()

	if minSize > bufSize {
		panic("bufSize must be >= minSize")
	}
	errch := make(chan error, 1)
	datach := make(chan *Chunk, numBufs)
	reusech := make(chan *Chunk, numBufs)

	for range [numBufs]struct{}{} {
		reusech <- &Chunk{}
	}
	//for i := range bufs {
	//	reusech <- &bufs[i]
	//}


	rds := NewStats(32000)
	wrs := NewStats(32000)


	go func() {
		defer close(datach)

		for {
			b := <-reusech
			var err error
			var n int
			for {
				rds.next()
				n, err = in.Read(b.buf[b.len:])
				if n > 0 {
					rds.done()
				}

				rest := len(b.buf[b.len:])
				if n < 0 || rest < n {
					if err == nil {
						err = fmt.Errorf("Invalid read operation: 0 <= n:%v <= buflen:%v", n, rest)
					}
					break
				}
				b.len += n

				if b.len >= minSize || err != nil {
					break
				}
			} // for

			if b.len > 0 {
				datach <- b
			}
			if err != nil {
				if err != io.EOF {
					errch <- err
				}
				return
			}
		} // for
	}()

	startTime := time.Now()

	var err error
	for b := range datach {
		var n int
		wrs.next()
		n, err = out.Write(b.buf[:b.len])
		wrs.done()
		if err != nil {
			break
		}
		if n != b.len {
			errch <- fmt.Errorf("Invalid write operation: n:%v == buflen:%v", n, b.len)
			break
		}
		b.len = 0
		reusech <- b
	}

	close(reusech)
	for range reusech {}

	close(errch)
	e := <- errch
	if e != nil {
		err = e
	}
	t.Logf("read stats:\n%v\n", rds)
	t.Logf("write stats:\n%v\n", wrs)
	t.Logf("overall : %v\n", time.Now().Sub(startTime))
	t.Logf("overall0: %v\n", time.Now().Sub(startTime0))
	return err
}
