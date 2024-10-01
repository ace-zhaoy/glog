package slog_test

import (
	"context"
	"log/slog"

	"github.com/ace-zhaoy/glog"
	slog2 "github.com/ace-zhaoy/glog/slog"
)

func ExampleHandler() {
	logger, _ := glog.NewDefault(
		glog.AddCallerSkip(3),
		glog.WithContextHandlers(
			glog.BuildContextHandler("req-id", "req_id"),
			glog.BuildContextHandler("user-id", "user_id"),
		),
	)
	slog.SetDefault(slog.New(slog2.NewHandler(logger)))

	ctx := context.WithValue(context.Background(), "req-id", "123")
	ctx = context.WithValue(ctx, "user-id", "456")

	slog.DebugContext(ctx, "Logging is enabled for level: debug")
	slog.InfoContext(ctx, "Logging is enabled for level: info")
	slog.WarnContext(ctx, "Logging is enabled for level: warn")
	slog.ErrorContext(ctx, "Logging is enabled for level: error")

	// Example Output:
	// {"level":"debug","ts":"2024-10-01T13:15:31.872349024Z","caller":"slog/example_test.go:24","msg":"Logging is enabled for level: debug","req_id":"123","user_id":"456"}
	// {"level":"info","ts":"2024-10-01T13:15:31.872422976Z","caller":"slog/example_test.go:25","msg":"Logging is enabled for level: info","req_id":"123","user_id":"456"}
	// {"level":"warn","ts":"2024-10-01T13:15:31.872431512Z","caller":"slog/example_test.go:26","msg":"Logging is enabled for level: warn","req_id":"123","user_id":"456"}
	// {"level":"error","ts":"2024-10-01T13:15:31.872436361Z","caller":"slog/example_test.go:27","msg":"Logging is enabled for level: error","req_id":"123","user_id":"456","stacktrace":"github.com/ace-zhaoy/glog/slog_test.ExampleHandler\n\t/go/src/github.com/ace-zhaoy/glog/slog/example_test.go:27\ntesting.runExample\n\t/usr/local/go/src/testing/run_example.go:63\ntesting.runExamples\n\t/usr/local/go/src/testing/example.go:44\ntesting.(*M).Run\n\t/usr/local/go/src/testing/testing.go:1927\nmain.main\n\t_testmain.go:61\nruntime.main\n\t/usr/local/go/src/runtime/proc.go:267"}
}
