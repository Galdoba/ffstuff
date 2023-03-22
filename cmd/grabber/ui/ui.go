package main

import (
	"fmt"
	"time"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

// func main() {
// 	inputinfo.CleanScanData()
// }

type allProc struct {
	procs []*proc
	//total    int
	//endEvent bool
	//activeSelection int
}

func (ap *allProc) String() string {
	str := ""
	for i := range ap.procs {
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
	//ap.total = q
	return &ap
}

func (ap *allProc) update() {
	if len(ap.procs) == 0 {
		//	ap.endEvent = true
		return
	}

	switch ap.procs[0].done {
	case false:
		ap.procs[0].run()
	case true:
		if len(ap.procs) > 1 {
			ap.procs = ap.procs[1:]
			//ap.activeSelection--
			ap.sync()
		} else {
			ap.procs = nil
			//ap.endEvent = true
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
	//ap.activeSelection = 0
}

func (ap *allProc) bumpUpByOne(i int) {
	//func namesort.BumpToTopIndex(slInt []int, index int) []int
	if i < 1 || i > len(ap.procs)-1 {
		return
	}
	ap.procs[i-1], ap.procs[i] = ap.procs[i], ap.procs[i-1]
}

func (ap *allProc) sync() {
	for i := range ap.procs {
		ap.procs[i].pos = i
	}
	//ap.total = len(ap.procs)
	//active = ap.total
	// if ap.activeSelection < 0 {
	// 	ap.activeSelection = 0
	// }
	// if ap.activeSelection > len(ap.procs)-1 {
	// 	ap.activeSelection = len(ap.procs) - 1
	// }
}

type proc struct {
	name string
	prog int
	pos  int
	done bool
}

func (p *proc) run() {
	if p.prog < 100 {
		p.prog++
	} else {
		p.done = true
	}
}

func (p *proc) String() string {

	return fmt.Sprintf("%v %v %v", prc(p.prog), p.name, p.pos)
}

func prc(i int) string {
	n := fmt.Sprintf("%v ", i)
	for len(n) < 4 {
		n = " " + n
	}
	return fmt.Sprintf("[%v", n) + "%]"
}

type InfoBox struct {
	data   []string
	active int
}

func (ib *InfoBox) Draw(ac *allProc) {
	fg := termbox.ColorWhite
	bg := termbox.ColorBlack
	tbprint(0, 1, fg, bg, "Grabber Dowloading:")
	for i, data := range ib.data {

		if i == ib.active {
			fg = termbox.ColorBlack
			bg = termbox.ColorWhite
		}
		tbprint(2, i+2, fg, bg, data)
		fg = termbox.ColorWhite
		bg = termbox.ColorBlack
	}
}

// This function is often useful:
func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {

	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}

func (ib *InfoBox) Update(ac *allProc) error {
	newData := []string{}
	for _, pr := range ac.procs {
		newData = append(newData, pr.String())
	}
	ib.data = newData
	//ib.active = ac.activeSelection

	return nil
}

func main() {
	ap := newAllProcs(10)
	ib := &InfoBox{}
	ib.data = []string{}
	ib.active = 3
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)
	event_queue := make(chan termbox.Event)
	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()
	ib.Draw(ap)
	termbox.Flush()
	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
	ap.update()
	ib.Update(ap)
	ib.Draw(ap)
	draw_tick := time.NewTicker(100 * time.Millisecond)
	i := 0

loop:
	for {
		if i > 1000 {
			break
		}
		i++
		//ib.Update(ap)
		if len(ap.procs) == 0 {
			ib.Update(ap)
			termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
			tbprint(0, 1, 0, 0, "DONE")
			termbox.Flush()
			fmt.Println("All GRABBED!")
			break
		}
		select {
		case ev := <-event_queue:
			switch ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyEsc:
					//os.Exit(1)
					break loop
				case termbox.KeyArrowDown:
					Action_MoveSelectionDOWN(ap, ib)
				case termbox.KeyArrowUp:
					Action_MoveSelectionUP(ap, ib)
				case termbox.KeyEnter:
					//ap.bumpToTop(ib.active)
					Action_BumpStreamToTOP(ap, ib)
				}
			}
		case <-draw_tick.C:
			ap.update()
		}
		ib.Update(ap)
		termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
		ib.Draw(ap)
		termbox.Flush()

	}
	fmt.Println("END")

}

func (ap *allProc) Options() []string {
	opt := []string{}
	for i := range ap.procs {
		opt = append(opt, ap.procs[i].String())
	}
	return opt
}

var active int

func qqq() {
	event_queue := make(chan termbox.Event)
	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()
	posx, posy := -1, 0
	color := termbox.ColorDefault
	color_change_tick := time.NewTicker(1 * time.Second)
	draw_tick := time.NewTicker(30 * time.Millisecond)

loop:
	for {
		select {
		case ev := <-event_queue:
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
				break loop
			}
		case <-color_change_tick.C:
			color++
			if color >= 8 {
				color = 0
			}
		case <-draw_tick.C:
			w, h := termbox.Size()
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			posx++
			if posx >= w {
				posx = 0
				posy++
				if posy >= h {
					posy = 0
				}
			}
			termbox.SetCell(posx, posy, '_',
				termbox.ColorDefault, color)
			termbox.Flush()
		}
	}
}

/*












 */
////////////////////

func Action_MoveSelectionUP(ap *allProc, ib *InfoBox) {
	ib.active--
	if ib.active < 0 {
		ib.active = 0
	}
}

func Action_MoveSelectionDOWN(ap *allProc, ib *InfoBox) {
	ib.active++
	for ib.active >= len(ap.procs) {
		ib.active--
	}
}

func Action_MoveSelectionTO(ap *allProc, ib *InfoBox, pos int) {
	ib.active = pos
	if ib.active < 0 {
		ib.active = 0
	}
	for ib.active >= len(ap.procs) {
		ib.active--
	}
}

func Action_BumpStreamToTOP(ap *allProc, ib *InfoBox) {
	ap.bumpToTop(ib.active)
	Action_MoveSelectionTO(ap, ib, 0)
}

func Action_BumpStreamByOne(ap *allProc, ib *InfoBox) {
	ap.bumpUpByOne(ib.active)
	Action_MoveSelectionUP(ap, ib)
}
