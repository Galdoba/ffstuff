package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Galdoba/utils"
)

// func main() {
// 	inputinfo.CleanScanData()
// }

type allProc struct {
	procs    []*proc
	total    int
	endEvent bool
}

func (ap *allProc) String() string {
	str := ""
	for i := 0; i < ap.total; i++ {
		p := ap.procs[i]
		str += fmt.Sprintf("%v\n", p)
	}
	return str
}

func newAllProcs(q int) *allProc {
	ap := allProc{}
	for i := 0; i < q; i++ {
		pr := proc{fmt.Sprintf("pr%v", i), 0, i, false}
		ap.procs = append(ap.procs, &pr)
	}
	ap.total = q
	return &ap
}

func (ap *allProc) update() {
	if len(ap.procs) == 0 {
		ap.endEvent = true
		return
	}
	switch ap.procs[0].done {
	case false:
		ap.procs[0].run()
	case true:
		if len(ap.procs) > 1 {
			ap.procs = ap.procs[1:]
			ap.sync()
		} else {
			ap.procs = nil
			ap.endEvent = true
			return
		}
	}
}

func (ap *allProc) bumpToTop(i int) {
	//func namesort.BumpToTopIndex(slInt []int, index int) []int
	if i < 1 || i > len(ap.procs)-1 {
		return
	}
	newSl := []*proc{}
	for j := range ap.procs {
		switch {
		case j == 0:
			newSl = append(newSl, ap.procs[i])
		case j <= i:
			newSl = append(newSl, ap.procs[j-1])
		case j > i:
			newSl = append(newSl, ap.procs[j])
		}
	}
	ap.procs = newSl
	for i := range ap.procs {
		ap.procs[i].pos = i
	}
}

func (ap *allProc) sync() {
	for i := range ap.procs {
		ap.procs[i].pos = i
	}
	ap.total = len(ap.procs)
}

type proc struct {
	name string
	prog int
	pos  int
	done bool
}

func (p *proc) run() {
	if p.prog < 10 {
		p.prog++
	} else {
		p.done = true
	}
}

func (p *proc) String() string {
	return fmt.Sprintf("%v %v %v", p.pos, p.name, p.prog)
}

func main() {
	ap := newAllProcs(40)
	ch := make(chan string)
	go func(ch chan string) {
		// Uncomment this block to actually read from stdin
		reader := bufio.NewReader(os.Stdin)
		for {
			s, err := reader.ReadString('\n')
			if err != nil { // Maybe log non io.EOF errors, if you want
				close(ch)
				return
			}
			ch <- s
		}
		// Simulating stdin
		//ch <- "A line of text"
		//close(ch)
	}(ch)

stdinloop:
	for {
		ap.update()
		if ap.endEvent {
			fmt.Println("Graceful return here")
			os.Exit(0)
		}
		select {
		case stdin, ok := <-ch:
			if !ok {
				break stdinloop
			} else {
				strIn := strings.TrimSpace(stdin)
				i, err := strconv.Atoi(strIn)
				if err != nil {
					fmt.Println(err.Error())
				} else {
					ap.bumpToTop(i)
				}
				fmt.Println("stdin was:", strIn)

			}
		case <-time.After(150 * time.Millisecond):
			// Do something when there is nothing read from stdin
			utils.ClearScreen()
			fmt.Println(ap)
		}
	}
	fmt.Println("Done, stdin must be closed")

}
