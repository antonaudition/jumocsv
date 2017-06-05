package jumocsv

import (
	"encoding/csv"
	"fmt"
	"io"
	"sort"
	"strings"
)

type Collator struct {
	r *Reader
}

func NewCollator(r *Reader) *Collator {
	return &Collator{r}
}

type counter struct {
	total  uint64
	count  uint64
	record Record
}

func (s *counter) update(amount uint64) {
	s.total += amount
	s.count++
}

func (c *Collator) Collate(w io.Writer) error {
	totals := make(map[Record]*counter)
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
			c = &counter{}
			totals[rec] = c
		}
		c.update(tmp)
	}
	var counters []*counter
	for rec, total := range totals {
		// go range variables get reused so cannot take its address. have to copy
		count := total
		count.record = rec
		counters = append(counters, count)
	}
	sort.Slice(counters, func(i, j int) bool {
		l, r := counters[i].record, counters[j].record
		x := l.Year - r.Year
		switch {
		case x < 0:
			return true
		case x > 0:
			return false
		}
		m := l.Month - r.Month
		switch {
		case m < 0:
			return true
		case m > 0:
			return false
		}
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
		return false
	})
	headers := []string{"Network", "Product", "Month", "Total", "Count"}
	csvW := csv.NewWriter(w)
	if err := csvW.Write(headers); err != nil {
		return err
	}
	for _, count := range counters {
		if err := csvW.Write(counterEntry(count)); err != nil {
			return err
		}
	}
	csvW.Flush()
	return nil
}

func counterEntry(count *counter) []string {
	rec := count.record
	return []string{
		rec.Network,
		rec.Product,
		fmt.Sprintf("%s-%d", rec.Month, rec.Year),
		currency(count.total),
		fmt.Sprintf("%d", count.count),
	}
}

func currency(amount uint64) string {
	dollar, cents := amount/100, amount%100
	return fmt.Sprintf("%d.%02d", dollar, cents)
}
