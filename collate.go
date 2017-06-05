package jumocsv

import (
	"encoding/csv"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

type Collator struct {
	r *Reader
}

func NewCollator(r io.Reader) *Collator {
	return &Collator{NewReader(r)}
}

// using a map to keep track of the (Network, Product, Month) tuple aggregations
// using the Record instance as a compound key so each entry only results in a single
// map lookup
func (c *Collator) Collate(w io.Writer) error {
	totals := make(map[Record]*aggregate)
	if _, err := c.r.Headers(); err != nil {
		return err
	}
	var tmp uint64
	for {
		rec, err := c.r.Read()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		tmp, rec.Amount = rec.Amount, 0
		c, ok := totals[rec]
		if !ok {
			c = &aggregate{}
			totals[rec] = c
		}
		c.update(tmp)
	}
	var aggs []*aggregate
	for rec, a := range totals {
		a.setRecord(rec)
		aggs = append(aggs, a)
	}
	sort.Sort(byMonth(aggs))
	headers := []string{"Network", "Product", "Month", "Total", "Count"}
	csvW := csv.NewWriter(w)
	if err := csvW.Write(headers); err != nil {
		return err
	}
	for _, a := range aggs {
		if err := csvW.Write(a.columns()); err != nil {
			return err
		}
	}
	csvW.Flush()
	return nil
}

type aggregate struct {
	total  uint64
	count  uint64
	record Record
	date   time.Time
}

func (a *aggregate) update(amount uint64) {
	a.total += amount
	a.count++
}

func (a *aggregate) setRecord(rec Record) {
	a.record = rec
	a.date = time.Date(rec.Year, rec.Month, 1, 1, 1, 1, 1, time.Local)
}

const outputDateFormat = "Jan-2006"

func (a *aggregate) columns() []string {
	rec := a.record
	dollar, cents := a.total/100, a.total%100
	return []string{
		fmt.Sprintf("'%s'", rec.Network),
		fmt.Sprintf("'%s'", rec.Product),
		fmt.Sprintf("'%s'", a.date.Format(outputDateFormat)),
		fmt.Sprintf("%d.%02d", dollar, cents),
		fmt.Sprintf("%d", a.count),
	}
}

type byMonth []*aggregate

func (a byMonth) Len() int      { return len(a) }
func (a byMonth) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byMonth) Less(i, j int) bool {
	l, r := a[i].record, a[j].record
	switch strings.Compare(l.Network, r.Network) {
	case -1:
		return true
	case 1:
		return false
	}
	switch strings.Compare(l.Product, r.Product) {
	case -1:
		return true
	case 1:
		return false
	}
	if a[i].date.Before(a[j].date) {
		return true
	}
	return false
}
