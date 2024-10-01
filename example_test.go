package glog_test

import (
	"context"
	"github.com/ace-zhaoy/glog"
)

func ExampleLogger() {
	logger, err := glog.NewDefault()
	if err != nil {
		panic(err)
	}

	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")

	// Example Output:
	// {"level":"info","ts":"2024-10-01T12:29:51.329714215Z","caller":"glog/example_test.go:14","msg":"This is an info message"}
	// {"level":"warn","ts":"2024-10-01T12:29:51.329883511Z","caller":"glog/example_test.go:15","msg":"This is a warning message"}
	// {"level":"error","ts":"2024-10-01T12:29:51.329891144Z","caller":"glog/example_test.go:16","msg":"This is an error message","stacktrace":"github.com/ace-zhaoy/glog_test.ExampleLogger\n\t/app/github.com/ace-zhaoy/glog/example_test.go:16\ntesting.runExample\n\t/usr/local/go/src/testing/run_example.go:63\ntesting.runExamples\n\t/usr/local/go/src/testing/example.go:44\ntesting.(*M).Run\n\t/usr/local/go/src/testing/testing.go:1721\nmain.main\n\t_testmain.go:123\nruntime.main\n\t/usr/local/go/src/runtime/proc.go:250"}
}

func ExampleLogger_WithContext() {
	logger, err := glog.NewDefault(
		glog.WithContextHandlers(
			glog.BuildContextHandler("request-id"),
			glog.BuildContextHandler("uid"),
		),
	)
	if err != nil {
		panic(err)
	}

	ctx := context.WithValue(context.Background(), "request-id", "123456")
	ctx = context.WithValue(ctx, "uid", "ace-zhaoy")
	ctx = context.WithValue(ctx, "foo", "bar")
	logger.InfoContext(ctx, "This is an info message with context")

	// Example Output:
	// {"level":"info","ts":"2024-10-01T20:26:13.5240654+08:00","caller":"glog/example_test.go:39","msg":"This is an info message with context","request-id":"123456","uid":"ace-zhaoy"}
}
