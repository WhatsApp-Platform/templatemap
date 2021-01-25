package templatemap

import (
	"strings"
	"testing"
)

func TestLoadDir(t *testing.T) {
	tmap, err := LoadDir("templates")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tmap)

	// Test t1
	if _, ok := tmap["t1.tmpl"]; !ok {
		t.Fatal("t1.tmpl not loaded")
	}

	// Test t2
	b := new(strings.Builder)
	t2, ok := tmap["sub/t2.tmpl"]
	if !ok {
		t.Fatal("sub/t2.tmpl not loaded")
	}

	err = t2.Execute(b, nil)
	if err != nil {
		t.Fatal("Failed to execute sub/t2.tmpl: ", err)
	}

	if b.String() != "I'm base(d) too" {
		t.Fatal("Unexpected value for sub/t2.tmpl: ", b.String())
	}

	// Test t3
	b = new(strings.Builder)
	t3, ok := tmap["sub/t3.tmpl"]
	if !ok {
		t.Fatal("sub/t3.tmpl not loaded")
	}

	err = t3.Execute(b, nil)
	if err != nil {
		t.Fatal("Failed to execute sub/t3.tmpl: ", err)
	}

	if b.String() != "I'm also base(d)" {
		t.Fatal("Unexpected value for sub/t3.tmpl: ", b.String())
	}

	// Test t4
	b = new(strings.Builder)
	t4, ok := tmap["sub/sub/t4.tmpl"]
	if !ok {
		t.Fatal("sub/aub/t4.tmpl not loaded")
	}

	err = t4.Execute(b, nil)
	if err != nil {
		t.Fatal("Failed to execute sub/sub/t4.tmpl: ", err)
	}

	if b.String() != "I'm also double base(d)" {
		t.Fatal("Unexpected value for sub/sub/t4.tmpl: ", b.String())
	}
}
