package base

import (
	"net/url"
	"testing"
)

func TestGenOrderID(t *testing.T) {
	id := GenOrderID("abcdefghijklmnopqrstuvwxyz1234567890")
	if len(id) > orderIDMaxLen {
		t.FailNow()
	}
	id = GenOrderID("abcdefg")
	if len(id) > orderIDMaxLen {
		t.FailNow()
	}
}

func TestURLScheme(t *testing.T) {
	g, err := url.Parse("http://www.google.com")
	if err != nil || g.Scheme != "http" {
		t.Log(err)
		t.FailNow()
	}
}
