package jumocsv

import (
	"bytes"
	"testing"
	"strings"
	"io"
	"fmt"
)

var testData = []string{
	`27729554427,'Network 1','12-Mar-2016','Loan Product 1',1000.00`,
	`27722342551,'Network 2','16-Mar-2016','Loan Product 1',1122.00`,
	`27725544272,'Network 3','17-Mar-2016','Loan Product 2',2084.00`,
	`27725326345,'Network 1','18-Mar-2016','Loan Product 2',3098.00`,
	`27729234533,'Network 2','01-Apr-2016','Loan Product 1',5671.00`,
	`27723453455,'Network 3','12-Apr-2016','Loan Product 3',1928.00`,
	`27725678534,'Network 2','15-Apr-2016','Loan Product 3',1747.00`,
	`27729554427,'Network 1','16-Apr-2016','Loan Product 2',1801.00`,
}

const (
	csvHeader = "MSISDN,Network,Date,Product,Amount"
)

type mockReader struct {
	lines []string
	first bool
	i     int
}

func (r *mockReader) Read() (record []string, err error) {
	if r.first {
		r.first = false
		return strings.Split(csvHeader, ","), nil
	}
	if r.i >= len(r.lines) {
		return nil, io.EOF
	}
	defer func() { r.i++ }()
	return strings.Split(r.lines[r.i], ","), nil
}

func testReader(data []string) *Reader {
	r := NewReader(nil)
	r.csv = &mockReader{lines: data, first: true}
	return r
}

func TestReader(t *testing.T) {
	r := testReader(testData)
	if _, err := r.Headers(); err != nil {
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
	w := bytes.NewBuffer(nil)
	c := NewCollator(nil)
	c.r = testReader(append(testData, testData...))
	if err := c.Collate(w); err != nil {
		t.Fatalf("out err: %v", err)
	}
	t.Log(w.String())
}

func TestCollateSum(t *testing.T) {
	var td []string
	for i := 1; i < 10; i++ {
		td = append(td, fmt.Sprintf(`27729554427,'Network 1','%02d-Mar-2016','Loan Product 1',1000.01`, i))
	}
	c := NewCollator(nil)
	c.r = testReader(td)
	w := bytes.NewBuffer(nil)
	if err := c.Collate(w); err != nil {
		t.Fatal(err)
	}
	_, err := w.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	result, err := w.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	if strings.Compare(result, "'Network 1','Loan Product 1','Mar-2016',9000.09,9\n") != 0 {
		t.Fatalf("expected equal got: %q", result)
	}
}
