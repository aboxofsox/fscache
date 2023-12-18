package fscache

import (
	"os"
	"testing"
)

func TestFsCache(t *testing.T) {
	c := NewCache()

	c.Set("foo", "bar")

	if v, ok := c.Get("foo"); !ok || v != "bar" {
		t.Error("Get failed")
	}

	c.Delete("foo")

	if _, ok := c.Get("foo"); ok {
		t.Error("Delete failed")
	}

	os.Remove("test.gob")
}

func TestIo(t *testing.T) {
	c := NewCache()
	c.Set("foo", "bar")

	err := c.Save("test.gob")
	if err != nil {
		t.Error(err)
	}

	cc, err := Load("test.gob")
	if err != nil {
		t.Error(err)
	}

	if v, ok := cc.Get("foo"); !ok || v != "bar" {
		t.Error("Load failed")
	}

	cc.Set("baz", "qux")

	if v, ok := cc.Get("baz"); !ok || v != "qux" {
		t.Error("Set failed")
	}

	cc.Save("test.gob")

	ccc, err := Load("test.gob")
	if err != nil {
		t.Error(err)
	}

	if v, ok := ccc.Get("baz"); !ok || v != "qux" {
		t.Error("Load failed")
	}

	os.Remove("test.gob")

}
