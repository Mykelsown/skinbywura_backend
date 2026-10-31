package main

import (
	"strings"
	"testing"
)

const header = "name,category,price,compare_at_price,rating,review_count,skin_type,badge,image_url,description,volume\n"

func TestParseValidRow(t *testing.T) {
	csv := header + "Serum,Serum,12500,15000,4.5,10,Dry,New,/images/a.jpg,Nice serum,30ml\n"
	got, err := parse(strings.NewReader(csv))
	if err != nil || len(got) != 1 {
		t.Fatalf("err=%v got=%v", err, got)
	}
	p := got[0]
	if p.price != 12500 || !p.hasCompare || p.compareAt != 15000 || p.rating != 4.5 || p.badge != "New" {
		t.Fatalf("unexpected product %#v", p)
	}
}

func TestParseOptionalFieldsMayBeBlank(t *testing.T) {
	csv := header + "Cream,Moisturiser,9000,,,,Dry,,/images/b.jpg,Rich cream,50ml\n"
	got, err := parse(strings.NewReader(csv))
	if err != nil {
		t.Fatal(err)
	}
	if got[0].hasCompare || got[0].rating != 0 || got[0].reviewCount != 0 {
		t.Fatalf("unexpected %#v", got[0])
	}
}

func TestParseRejectsBadRows(t *testing.T) {
	bad := []string{
		"name,category\nx,y\n",
		header + "A,B,abc,,,,Dry,,/i.jpg,d,1ml\n",
		header + "A,B,-5,,,,Dry,,/i.jpg,d,1ml\n",
		header + "A,B,100,,9,,Dry,,/i.jpg,d,1ml\n",
		header + ",B,100,,,,Dry,,/i.jpg,d,1ml\n",
		header,
	}
	for i, c := range bad {
		if _, err := parse(strings.NewReader(c)); err == nil {
			t.Errorf("case %d should fail", i)
		}
	}
}
