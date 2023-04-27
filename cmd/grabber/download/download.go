package download

import (
	"fmt"
	"io"
	"os"
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

func (dj *downloadHandler) Listen() chan Response {
	return dj.resp_chan
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
	com_chan  chan Command
	resp_chan chan Response
}

type Handler interface {
	Listen() chan Response
	Continue()
	Pause()
	Kill()
}

/*


handler, err := download.NewHandler(source, dest string)
handler.Listen()
handler.Pause()
handler.Kill()

stream, err := download.StartNew(source, dest string)
stream.Listen()
stream.Pause()
stream.Kill()

*/

func StartNew(source, dest string) (*downloadHandler, error) {
	dj := downloadHandler{}
	in, err := os.Open(source)
	if err != nil {
		return &dj, err
	}
	out, err := os.Create(dest)
	if err != nil {
		in.Close()
		return &dj, err
	}
	dj.com_chan = make(chan Command)
	dj.resp_chan = make(chan Response)
	go func() {
		transferData(out, in, dj.com_chan, dj.resp_chan)
		in.Close()
		out.Close()
	}()
	return &dj, nil
}

func (ds *downloadHandler) Pause() {
	ds.com_chan <- com_Pause
}

func (ds *downloadHandler) Continue() {
	ds.com_chan <- com_Continue
}

func (ds *downloadHandler) Kill() {
	ds.com_chan <- com_Stop
}

//func transferData(out io.Writer, in io.Reader, com_chan <-chan Command, resp_chan chan<- Response) {
func transferData(out io.Writer, in io.Reader, com_chan <-chan Command, resp_chan chan<- Response) {
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
	var total int64
	var err error
	//blockwriteLoop:
	for b := range datach {
		var n int
		select {
		default: //не блокирует если не получает команды
		case cmd := <-com_chan:
			switch cmd {
			case com_Pause:
			blockwriteLoop:
				for com := range com_chan { //блокирующее ожидание
					switch com {
					case com_Pause:
						continue
					case com_Continue:
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
			total += int64(n)
			select {
			case resp_chan <- NewResponseProgress(total, 0):
			default:
			}
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
			resp_chan <- NewResponseTerminated()
		} else {
			resp_chan <- NewResponseError(err)
		}
	} else {
		resp_chan <- NewResponseCompleted()
	}
	close(resp_chan)
}

/*


 */
