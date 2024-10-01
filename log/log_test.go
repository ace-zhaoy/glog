package log

import (
	"testing"

	"github.com/ace-zhaoy/glog"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	l := &glog.Logger{}
	SetLogger(l)
	assert.Equal(t, l, Logger(), "Expected logger to be set and retrieved correctly")
}

func TestSetLogger(t *testing.T) {
	l := &glog.Logger{}
	SetLogger(l)
	assert.Equal(t, l, Logger(), "Expected logger to be set and retrieved correctly")
}

func TestWithFormatEnable(t *testing.T) {
	l := &glog.Logger{}
	SetLogger(l)
	assert.NotNil(t, WithFormatEnable(), "Expected WithFormatEnable to return a logger")
}

func TestWithFormatDisable(t *testing.T) {
	l := &glog.Logger{}
	SetLogger(l)
	assert.NotNil(t, WithFormatDisable(), "Expected WithFormatDisable to return a logger")
}
