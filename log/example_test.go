package log_test

import (
	"context"
	"github.com/ace-zhaoy/glog"
	"github.com/ace-zhaoy/glog/log"
)

func ExampleLogger() {
	logger := log.Logger().WithOptions(glog.AddCallerSkip(-1))
	logger.Info("This is an info message from log package")

	// Example Output:
	// {"level":"info","ts":"2024-10-01T12:47:47.65254347Z","caller":"log/example_test.go:11","msg":"This is an info message from log package"}
}

func ExampleLogContext() {
	oldLogger := log.Logger()
	defer log.SetLogger(oldLogger)
	log.SetLogger(log.Logger().WithOptions(
		glog.WithContextHandlers(
			glog.BuildContextHandler("request-id"),
		),
	))
	ctx := context.WithValue(context.Background(), "request-id", "123456")
	log.LogContext(ctx, glog.LevelInfo, "This is an info message with context")

	// Example Output:
	// {"level":"info","ts":"2024-10-01T12:50:54.30641073Z","caller":"log/example_test.go:24","msg":"This is an info message with context","request-id":"123456"}
}

func ExampleInfo() {
	log.Info("This is an info message")

	// Example Output:
	// {"level":"info","ts":"2024-10-01T12:55:42.826440061Z","caller":"log/example_test.go:33","msg":"This is an info message"}
}

func ExampleWarn() {
	log.Warn("This is a warning message")

	// Example Output:
	// {"level":"warn","ts":"2024-10-01T12:56:29.939595639Z","caller":"log/example_test.go:40","msg":"This is a warning message"}
}

func ExampleError() {
	log.Error("This is an error message")

	// Example Output:
	// {"level":"error","ts":"2024-10-01T12:57:20.8502351Z","caller":"log/example_test.go:47","msg":"This is an error message","stacktrace":"github.com/ace-zhaoy/glog/log_test.ExampleError\n\t/app/github.com/ace-zhaoy/glog/log/example_test.go:47\ntesting.runExample\n\t/usr/local/go/src/testing/run_example.go:63\ntesting.runExamples\n\t/usr/local/go/src/testing/example.go:44\ntesting.(*M).Run\n\t/usr/local/go/src/testing/testing.go:1721\nmain.main\n\t_testmain.go:57\nruntime.main\n\t/usr/local/go/src/runtime/proc.go:250"}
}
