package stacktrace

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCapture(t *testing.T) {
	stack := Capture(0, Full)
	assert.NotNil(t, stack, "Expected captured stack to be non-nil")
	assert.Greater(t, stack.Count(), 0, "Expected captured stack to contain frames")

	stackStr := stack.String()
	assert.True(t, strings.Contains(stackStr, "TestCapture"), "Expected stack trace to contain function name")
}

func TestFormatter(t *testing.T) {
	stack := Capture(0, Full)
	formatter := GetFormatter()
	defer formatter.Free()

	formatter.FormatStack(stack)
	formattedStr := formatter.String()

	assert.True(t, strings.Contains(formattedStr, "TestFormatter"), "Expected formatted stack to contain function name")
	assert.True(t, strings.Contains(formattedStr, "stack_test.go"), "Expected formatted stack to contain file name")
}

func TestTake(t *testing.T) {
	trace := Take(0)
	assert.True(t, strings.Contains(trace, "TestTake"), "Expected stack trace to contain function name")
}
