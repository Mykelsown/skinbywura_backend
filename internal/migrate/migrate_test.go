package migrate

import (
	"testing"
	"testing/fstest"
)

func TestLoadOrdersAndIgnoresDownFiles(t *testing.T) {
	fsys := fstest.MapFS{
		"m/000002_b.up.sql":   {Data: []byte("SELECT 2;")},
		"m/000001_a.up.sql":   {Data: []byte("SELECT 1;")},
		"m/000001_a.down.sql": {Data: []byte("SELECT 0;")},
		"m/000010_c.up.sql":   {Data: []byte("SELECT 10;")},
	}
	got, err := Load(fsys, "m")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 || got[0].Version != 1 || got[1].Version != 2 || got[2].Version != 10 {
		t.Fatalf("unexpected order: %#v", got)
	}
}

func TestLoadRejectsDuplicateVersions(t *testing.T) {
	fsys := fstest.MapFS{
		"m/000001_a.up.sql": {Data: []byte("x")},
		"m/1_b.up.sql":      {Data: []byte("y")},
	}
	if _, err := Load(fsys, "m"); err == nil {
		t.Fatal("expected duplicate error")
	}
}

func TestPending(t *testing.T) {
	all := []Migration{{Version: 1}, {Version: 2}, {Version: 3}}
	if p := Pending(all, 2); len(p) != 1 || p[0].Version != 3 {
		t.Fatalf("unexpected pending %#v", p)
	}
	if p := Pending(all, 3); len(p) != 0 {
		t.Fatal("nothing should be pending")
	}
}

func TestParseVersionErrors(t *testing.T) {
	if _, err := ParseVersion("noprefix.up.sql"); err == nil {
		t.Fatal("expected error")
	}
	if _, err := ParseVersion("abc_x.up.sql"); err == nil {
		t.Fatal("expected error")
	}
}
