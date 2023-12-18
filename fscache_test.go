package fscache

import (
	"os"
	"testing"
)

func TestFsCache(t *testing.T) {
	c := NewCache("test.gob")

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
	c := NewCache("test.gob")
	c.Set("foo", "bar")

	if _, err := os.Stat("test.gob"); os.IsNotExist(err) {
		t.Error("test.gob does not exist")
		return
	}

	err := c.Save()
	if err != nil {
		t.Error(err)
		return
	}

	cc, err := Load("test.gob")
	if err != nil {
		t.Error(err)
		return
	}

	if v, ok := cc.Get("foo"); !ok || v != "bar" {
		t.Error("Load failed")
	}

	cc.Set("baz", "qux")

	if v, ok := cc.Get("baz"); !ok || v != "qux" {
		t.Error("Set failed")
	}

	err = cc.Save()
	if err != nil {
		t.Error(err)
	}

	ccc, err := Load("test.gob")
	if err != nil {
		t.Error(err)

	}

	if v, ok := ccc.Get("baz"); !ok || v != "qux" {
		t.Error("Load failed")
	}

	os.Remove("test.gob")

}
