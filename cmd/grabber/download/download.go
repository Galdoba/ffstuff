package download

import (
	"fmt"
	"io"
	"os"
	"sync"
)

const (
	com_Pause = Command(iota)
	com_Continue
	com_Stop
)

type Command int
type Response struct {
	err        error
	progress   int64 //written bytes total
	dwSpeed    int   //kb/s
	completed  bool
	terminated bool
}

const (
	MSG_abortedByUser = "transfert aborted by user"
)

func NewResponseError(err error) Response {
	return Response{
		err:      err,
		progress: -1,
		dwSpeed:  -1,
	}
}

func NewResponseProgress(progr int64, speed int) Response {
	return Response{
		err:      nil,
		progress: progr,
		dwSpeed:  speed,
	}
}

func NewResponseCompleted() Response {
	return Response{
		completed: true,
	}
}

func NewResponseTerminated() Response {
	return Response{
		terminated: true,
	}
}

func (r *Response) String() string {
	if r.completed {
		return "completed"
	}
	if r.terminated {
		return "terminated"
	}
	if r.err != nil {
		return fmt.Sprintf("error responce:  %v", r.err.Error())
	}
	return fmt.Sprintf("%v", r.progress)
}

type downloadHandler struct {
	mutex        sync.Mutex
	com_chan     chan Command
	status       Status
	progress     int64
	fileSize     int64
	speedLimit   int64
	speedCurrent int64
	err          error
}

type Status int

const (
	STATUS_NIL Status = iota
	STATUS_TRANSFERING
	STATUS_PAUSED
	STATUS_ERR
	STATUS_TERMINATED
	STATUS_COMPLETED
)

type Handler interface {
	//send command
	Continue()
	Pause()
	Close()
	//getter
	Status() Status
	Error() error
	Progress() int64
	FileSize() int64
}

var _ Handler = (*downloadHandler)(nil)

func StartNew(source, dest string) *downloadHandler {
	dj := downloadHandler{}
	in, err := os.Open(source)
	if err != nil {
		dj.setError(err)
		return &dj
	}

	dj.com_chan = make(chan Command)

	go func() {
		defer func() {
			close(dj.com_chan)
			for range dj.com_chan {
			}
		}()
		defer in.Close()

		out, err := os.Create(dest)
		defer out.Close()
		if err != nil {
			dj.setError(err)
			return
		}

		//assert(!dj.isReady(), "MUST NOT BE HAPPENED")

		dj.transferData(out, in)

		//assert(!dj.isReady(), "MUST NOT BE HAPPENED")

		// close(dj.com_chan)
		// for range dj.com_chan {
		// }
	}()

	fi, err := in.Stat()
	if err != nil {
		dj.setError(err)
		return &dj
	}
	dj.fileSize = fi.Size()
	return &dj
}

func assert(ok bool, msg string) {
	if !ok {
		panic(msg)
	}
}

func (dj *downloadHandler) isReady() bool {
	status := dj.Status()
	if status < STATUS_ERR {
		return true
	}
	return false
}

func (dj *downloadHandler) Pause() {

	if dj.isReady() {
		dj.com_chan <- com_Pause
	}
}

func (dj *downloadHandler) Continue() {
	if dj.isReady() {
		dj.com_chan <- com_Continue
	}

}

func (dj *downloadHandler) Close() {
	if dj.isReady() {
		dj.com_chan <- com_Stop
	}
}

func (dj *downloadHandler) Status() Status {
	dj.mutex.Lock()
	defer dj.mutex.Unlock()
	return dj.status
}

func (dj *downloadHandler) Error() error {
	dj.mutex.Lock()
	defer dj.mutex.Unlock()
	return dj.err
}

func (dj *downloadHandler) Progress() int64 {
	dj.mutex.Lock()
	defer dj.mutex.Unlock()
	return dj.progress
}

func (dj *downloadHandler) FileSize() int64 {
	dj.mutex.Lock()
	defer dj.mutex.Unlock()
	return dj.fileSize
}

func (dj *downloadHandler) setStatus(s Status) {
	dj.mutex.Lock()
	dj.status = s
	dj.mutex.Unlock()
}

func (dj *downloadHandler) setError(e error) {
	dj.mutex.Lock()
	dj.err = e
	dj.status = STATUS_ERR
	dj.mutex.Unlock()
}

func (dj *downloadHandler) transferData(out io.Writer, in io.Reader) {
	const (
		minSize = 1024
		bufSize = 1024 * 1024
		numBufs = 8
	)
	type Chunk struct {
		buf [bufSize]byte
		len int
	}
	if minSize > bufSize {
		panic("bufSize must be >= minSize")
	}
	errch := make(chan error)
	datach := make(chan *Chunk, numBufs)
	reusech := make(chan *Chunk, numBufs)

	for range [numBufs]struct{}{} {
		reusech <- &Chunk{}
	}
	dj.setStatus(STATUS_TRANSFERING)
	go func() {
		defer close(datach)
		for {
			b, ok := <-reusech
			if !ok {
				break
			}
			var err error
			var n int
			for {
				n, err = in.Read(b.buf[b.len:])
				rest := len(b.buf[b.len:])
				if n < 0 || rest < n {
					if err == nil {
						err = fmt.Errorf("invalid read operation: 0 <= n:%v <= buflen:%v", n, rest)
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
	}() //go func
	//var total int64
	var err error
	//blockwriteLoop:
	for b := range datach {
		var n int
		select {
		default: //не блокирует если не получает команды
		case cmd := <-dj.com_chan:
			switch cmd {
			case com_Pause:
				dj.setStatus(STATUS_PAUSED)
			blockwriteLoop:
				for com := range dj.com_chan { //блокирующее ожидание
					switch com {
					case com_Pause:
						dj.setStatus(STATUS_PAUSED)
						continue
					case com_Continue:
						dj.setStatus(STATUS_TRANSFERING)
						break blockwriteLoop
					case com_Stop:
						//закрываем лавочку
						err = fmt.Errorf(MSG_abortedByUser)
						break blockwriteLoop
					}
				}
			case com_Continue:
			case com_Stop:
				//закрываем лавочку
				err = fmt.Errorf(MSG_abortedByUser)
			}
		} //select
		if err == nil {
			n, err = out.Write(b.buf[:b.len])
		}
		if err != nil {
			break
		}
		if n != b.len {
			err = fmt.Errorf("invalid write operation: n:%v == buflen:%v", n, b.len)
			break
		}
		b.len = 0
		reusech <- b
		if n != 0 {
			//total += int64(n)
			// select {
			// case resp_chan <- NewResponseProgress(total, 0):
			// default:
			// }
			dj.mutex.Lock()
			dj.progress += int64(n)
			dj.mutex.Unlock()
		}
	} // for b := range datach {
	close(reusech)
	for range reusech {
	}
	close(errch)

	e := <-errch
	if e == nil {
		e = err
	}
	// switch {
	// case e == nil:
	// 	resp_chan <- NewResponseCompleted()
	// case e.Error() == MSG_abortedByUser:
	// 	resp_chan <- NewResponseTerminated()
	// default:
	// 	resp_chan <- NewResponseError(err)
	// }
	if e != nil {
		if e.Error() == MSG_abortedByUser {
			//resp_chan <- NewResponseTerminated()
			dj.setStatus(STATUS_TERMINATED)
		} else {
			//resp_chan <- NewResponseError(err)
			dj.setError(e)
			dj.setStatus(STATUS_ERR)
		}
	} else {
		//resp_chan <- NewResponseCompleted()
		dj.setStatus(STATUS_COMPLETED)
	}
}
