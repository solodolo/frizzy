package processor

import (
	"testing"
)

func TestAtMissingKeyReturnsFalse(t *testing.T) {
	context := Context{}
	_, ok := context.At("foo")

	if ok {
		t.Errorf("expected key \"foo\" not to exist but it does")
	}
}

func TestContextStoresResults(t *testing.T) {
	context := Context{}
	strResult := StringResult("test")

	context["foo"] = ContextNode{result: strResult}
	val, ok := context.At("foo")
	if !ok {
		t.Errorf("expected key \"foo\" to exist but it does not")
	} else if val.result.GetResult() != strResult {
		t.Errorf("expected %q, got %q", strResult.String(), context["foo"].result.String())
	}
}

func TestContextsCanBeNested(t *testing.T) {
	expectedStr := StringResult("test")
	expectedInt := IntResult(5)
	context := Context{
		"foo": ContextNode{child: &Context{
			"bar": ContextNode{child: &Context{
				"baz": ContextNode{result: expectedStr},
			}},
			"fizz": ContextNode{result: expectedInt},
		}},
	}
	a, _ := context.At("foo")
	b, _ := a.At("bar")
	c, _ := b.At("baz")
	d, _ := a.At("fizz")

	if c.result.GetResult() != expectedStr {
		t.Errorf("expected %q, got %q", expectedStr, c.result.String())
	} else if d.result.GetResult() != expectedInt {
		t.Errorf("expected %d, got %d", expectedInt, d.result.GetResult())
	}
}
