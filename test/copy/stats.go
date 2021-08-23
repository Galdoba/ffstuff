package copy

import (
	"fmt"
	"sort"
	"time"
)

type tstats struct {
	valid bool
	data []time.Duration
	t time.Time
	min, max, avr, med, sum time.Duration
	li, ri int
	cnt int
}

func NewStats(datasize int) *tstats {
	return &tstats{data: make([]time.Duration, 0, datasize)}
}

func (o *tstats) next() {
	o.t = time.Now()
}

func (o *tstats) done() {
	o.valid = false
	dt := time.Now().Sub(o.t)
	o.data = append(o.data, dt)

	if o.cnt == 0 {
		o.min = dt
		o.max = dt
		o.avr = dt
		o.sum = dt
		o.cnt++
		return
	}

	if dt < o.min {
		o.min = dt
	}
	if dt > o.max {
		o.max = dt
	}
	o.sum += dt
	o.avr =  o.sum / time.Duration(o.cnt)
	o.cnt++
}

func (o *tstats) calc() {
	sort.Sort(durSlice(o.data))

	o.li = (len(o.data)-1) / 2
	o.ri = len(o.data) / 2
	o.med = (o.data[o.li]+o.data[o.ri]) / 2

	o.valid = true
}

func (o tstats) String() string {
	if !o.valid {
		o.calc()
	}
	return fmt.Sprintf(
` counts: %v
average: %v
 median: %v
    min: %v
    max: %v
    sum: %v
`,
		o.cnt, o.avr,  o.med, o.min, o.max, o.sum)
}

type durSlice []time.Duration

func (p durSlice) Len() int {
    return len(p)
}

func (p durSlice) Less(i, j int) bool {
    return p[i] < p[j]
}

func (p durSlice) Swap(i, j int) {
    p[i], p[j] = p[j], p[i]
}

