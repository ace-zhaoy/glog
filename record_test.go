package glog

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRecordWithCapacity(t *testing.T) {
	record := NewRecordWithCapacity(5)
	assert.NotNil(t, record)
	assert.Equal(t, 0, len(record.Fields()))
}

func TestAddFields(t *testing.T) {
	record := NewRecordWithCapacity(2)
	field1 := String("field1", "value1")
	field2 := String("field2", "value2")

	record.AddFields(field1, field2)
	fields := record.Fields()

	assert.Equal(t, 2, len(fields))
	assert.Equal(t, "field1", fields[0].Key)
	assert.Equal(t, "field2", fields[1].Key)
}

func TestAdd(t *testing.T) {
	record := NewRecordWithCapacity(3)
	record.Add("key", "value", String("key1", "value1"))

	fields := record.Fields()
	assert.Equal(t, 2, len(fields))
	assert.Equal(t, "key", fields[0].Key)
	assert.Equal(t, "value", fields[0].String)
	assert.Equal(t, "key1", fields[1].Key)
	assert.Equal(t, "value1", fields[1].String)
}

func TestBuildContextHandler(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "key", "value")
	handler := BuildContextHandler("key")
	record := NewRecordWithCapacity(1)
	handler(ctx, record)
	assert.Equal(t, 1, len(record.fields), "Expected one field to be added when key is provided and alias is empty")
	assert.Equal(t, "key", record.fields[0].Key, "Expected field key to be key when key is provided and alias is empty")
	assert.Equal(t, "value", record.fields[0].String, "Expected field value to be value when key is provided and alias is empty")

	ctx = context.Background()
	ctx = context.WithValue(ctx, "key", "value")
	handler = BuildContextHandler("key", "alias")
	record = NewRecordWithCapacity(1)
	handler(ctx, record)
	assert.Equal(t, 1, len(record.fields), "Expected one field to be added when key and alias are provided")
	assert.Equal(t, "alias", record.fields[0].Key, "Expected field key to be alias when key and alias are provided")
	assert.Equal(t, "value", record.fields[0].String, "Expected field value to be value when key and alias are provided")
}

func TestArgsToFields(t *testing.T) {
	args := []any{}
	expectedFields := []Field{}
	resultFields := argsToFields(args)
	assert.Equal(t, expectedFields, resultFields, "Expected empty fields for empty arguments")

	field := String("key", "value")
	args = []any{field}
	expectedFields = []Field{field}
	resultFields = argsToFields(args)
	assert.Equal(t, expectedFields, resultFields, "Expected same field for field argument")

	strArg := "key"
	valueArg := "value"
	args = []any{strArg, valueArg}
	expectedFields = []Field{Any(strArg, valueArg)}
	resultFields = argsToFields(args)
	assert.Equal(t, expectedFields, resultFields, "Expected field with key and value for string and value arguments")

	fieldArg := String("key2", "value2")
	args = []any{strArg, valueArg, fieldArg}
	expectedFields = []Field{Any(strArg, valueArg), fieldArg}
	resultFields = argsToFields(args)
	assert.Equal(t, expectedFields, resultFields, "Expected fields with key and value for string and value arguments, and the field argument")
}
