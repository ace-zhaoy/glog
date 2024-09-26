package stacktrace

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPool_GetAndPut(t *testing.T) {
	p := NewPool(func() int {
		return 42
	})

	val := p.Get()
	assert.Equal(t, 42, val, "Expected initial value to be 42")

	p.Put(100)
	val = p.Get()
	assert.Equal(t, 100, val, "Expected value to be 100 after Put")
}

func TestPool_WithCustomType(t *testing.T) {
	type customStruct struct {
		Field string
	}

	p := NewPool(func() *customStruct {
		return &customStruct{Field: "test"}
	})

	obj := p.Get()
	assert.Equal(t, "test", obj.Field, "Expected initial field value to be 'test'")

	p.Put(&customStruct{Field: "updated"})
	obj = p.Get()
	assert.Equal(t, "updated", obj.Field, "Expected field value to be 'updated' after Put")
}
