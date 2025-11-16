package storage

import "testing"

func TestMemtableBasic(t *testing.T) {
	mt := NewMemtable()

	mt.Set("foo", "bar")
	v, ok := mt.Get("foo")
	if !ok || v != "bar" {
		t.Fatalf("expected foo=bar, got %q, ok=%v", v, ok)
	}

	mt.Del("foo")
	_, ok = mt.Get("foo")
	if ok {
		t.Fatalf("expected foo to be deleted")
	}
}
