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

	context["foo"] = &ContextNode{result: strResult}
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
		"foo": &ContextNode{child: &Context{
			"bar": &ContextNode{child: &Context{
				"baz": &ContextNode{result: expectedStr},
			}},
			"fizz": &ContextNode{result: expectedInt},
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

func TestContextsCanInsertNestedKeys(t *testing.T) {
	keys := []string{"foo", "bar", "baz"}
	value := StringResult("the result")
	context := &Context{}

	context.Insert(keys, value)
	current := &ContextNode{child: context}
	for _, key := range keys {
		if node, ok := current.At(key); !ok {
			t.Errorf("expected context to include key %q but it does not", key)
		} else if key == keys[len(keys)-1] && node.result != value {
			t.Errorf("expected node to store result %q, got %q", value, node.result)
		} else {
			current = node
		}
	}
}

func TestContextAtNestedReturnsNestedNode(t *testing.T) {
	expectedStr := StringResult("test")
	context := Context{
		"foo": &ContextNode{child: &Context{
			"bar": &ContextNode{child: &Context{
				"baz": &ContextNode{result: expectedStr},
			}},
		}},
	}

	if at, ok := context.AtNested([]string{"foo", "bar", "baz"}); !ok {
		t.Errorf("expected nested keys to be found but they were not")
	} else if at.result != expectedStr {
		t.Errorf("expected nested result to equal %q, got %q", expectedStr, at.result)
	}
}

func TestContextAtReturnsNestedNode(t *testing.T) {
	expectedStr := StringResult("test")
	context := Context{
		"foo": &ContextNode{child: &Context{
			"bar": &ContextNode{child: &Context{
				"baz": &ContextNode{result: expectedStr},
			}},
		}},
	}

	if at, ok := context.At("foo.bar.baz"); !ok {
		t.Errorf("expected nested keys to be found but they were not")
	} else if at.result != expectedStr {
		t.Errorf("expected nested result to equal %q, got %q", expectedStr, at.result)
	}
}

func TestContextAtReturnsUnnestedNode(t *testing.T) {
	expectedStr := StringResult("test")
	context := Context{"foo": &ContextNode{result: expectedStr}}

	if at, ok := context.At("foo"); !ok {
		t.Errorf("expected nested keys to be found but they were not")
	} else if at.result != expectedStr {
		t.Errorf("expected nested result to equal %q, got %q", expectedStr, at.result)
	}
}

func TestContextAtNestedReturnsUnnestedNode(t *testing.T) {
	expectedStr := StringResult("test")
	context := Context{"foo": &ContextNode{result: expectedStr}}

	if at, ok := context.AtNested([]string{"foo"}); !ok {
		t.Errorf("expected nested keys to be found but they were not")
	} else if at.result != expectedStr {
		t.Errorf("expected nested result to equal %q, got %q", expectedStr, at.result)
	}
}
