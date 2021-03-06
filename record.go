package jumocsv

import (
	"encoding/csv"
	"io"
	"strconv"
	"strings"
	"time"
)

const (
	msisdnIdx  = int(iota)
	networkIdx
	dateIdx
	productIdx
	amountIdx
)

const dateFormat = "_2-Jan-2006"

type Record struct {
	Network string
	Product string
	Month   time.Month
	Year    int
	Amount  uint64
}

type csvReader interface {
	Read() (record []string, err error)
}

type Reader struct {
	csv csvReader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		csv.NewReader(r),
	}
}

func (r *Reader) Headers() ([]string, error) {
	return r.csv.Read()
}

func (r *Reader) Read() (rec Record, err error) {
	row, err := r.csv.Read()
	if err != nil {
		return
	}
	return parseRecord(row)
}

func parseRecord(rec []string) (record Record, err error) {
	date, err := time.Parse(dateFormat, unqoute(rec[dateIdx]))
	if err != nil {
		return
	}
	// storing the currency in an uint64 to avoid float overflows
	record.Amount, err = strconv.ParseUint(strings.Replace(rec[amountIdx], ".", "", -1), 10, 64)
	if err != nil {
		return
	}
	record.Network, record.Product = unqoute(rec[networkIdx]), unqoute(rec[productIdx])
	record.Year, record.Month, _ = date.Date()
	return
}

func unqoute(s string) string {
	if s[0] == '\'' {
		return s[1: len(s)-1]
	}
	return s
}
