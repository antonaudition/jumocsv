package jumocsv

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"os"
	"testing"
)

/*
MSISDN,Network,Date,Product,Amount
27729554427,'Network 1','12-Mar-2016','Loan Product 1',1000.00
27722342551,'Network 2','16-Mar-2016','Loan Product 1',1122.00
27725544272,'Network 3','17-Mar-2016','Loan Product 2',2084.00
27725326345,'Network 1','18-Mar-2016','Loan Product 2',3098.00
27729234533,'Network 2','01-Apr-2016','Loan Product 1',5671.00
27723453455,'Network 3','12-Apr-2016','Loan Product 3',1928.00
27725678534,'Network 2','15-Apr-2016','Loan Product 3',1747.00
27729554427,'Network 1','16-Apr-2016','Loan Product 2',1801.00
*/

func Test1(t *testing.T) {
	f, err := os.Open("./Loans.csv")
	if err != nil {
		panic(err)
	}
	r := csv.NewReader(bufio.NewReader(f))

	first := true
	for {
		rec, err := r.Read()
		if err != nil {
			t.Logf("read error: %v", err)
			break
		}
		if first {
			first = false
			continue
		}
		parsed, err := parseRecord(rec)
		t.Logf("%v, %v", parsed, err)
	}

}

func TestReader(t *testing.T) {
	r, err := NewReader("./Loans.csv")
	if err != nil {
		t.Fatalf("error opening file: %v", err)
	}
	if _, err = r.Headers(); err != nil {
		t.Fatalf("error opening file: %v", err)
	}
	for {
		rec, err := r.Read()
		if err != nil {
			t.Logf("read error: %v", err)
			break
		}
		t.Logf("%v", rec)
	}
}

func TestCollate(t *testing.T) {
	r, err := NewReader("./Loans.csv")
	if err != nil {
		t.Fatalf("error opening file: %v", err)
	}
	w := bytes.NewBuffer(nil)
	c := NewCollator(r)
	if err := c.Collate(w); err != nil {
		t.Fatalf("out err: %v", err)
	}
	t.Log(w.String())
}
