# glog

A user-friendly logging library based on `zapcore.Core` that automatically loads values from `context.Context`.

## Features

- Supports extracting field values from `context.Context` and recording them in logs
- User-friendly API
- Supports printf-style formatting
- Based on `zapcore.Core`

## Installation

To install the package, use the following command:

```sh
go get github.com/ace-zhaoy/glog
```

## Usage

### Quick Usage
#### Default Logger
```go
package main

import (
	"github.com/ace-zhaoy/glog/log"
)

func main() {
	log.Info("This is an info message with context")
}

// Output:
// {"level":"info","ts":"2024-10-01T21:44:08.2737888+08:00","caller":"glog/main.go:9","msg":"This is an info message with context"}
```

#### Set Default Logger
```go
package main

import (
	"context"
	"github.com/ace-zhaoy/glog"
	"github.com/ace-zhaoy/glog/log"
)

func main() {
	logger := log.Logger()
	logger = logger.WithOptions(
		glog.WithContextHandlers(
			glog.BuildContextHandler("request-id", "req_id"),
		),
	)
	log.SetLogger(logger)

	ctx := context.WithValue(context.Background(), "request-id", "123456")
	log.InfoContext(ctx, "This is an info message with context")
}

// Output:
// {"level":"info","ts":"2024-10-01T21:41:36.825259+08:00","caller":"glog/main.go:19","msg":"This is an info message with context","req_id":"123456"}
```
  

### Basic Usage

```go
package main

import (
	"github.com/ace-zhaoy/glog"
	"go.uber.org/zap/zapcore"
)

func main() {
	logger, err := glog.NewDefault()
	if err != nil {
		panic(err)
	}

	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")
}
```

### Contextual Logging

```go
package main

import (
	"context"
	"github.com/ace-zhaoy/glog"
)

func main() {
	logger, err := glog.NewDefault(
		glog.WithContextHandlers(
			glog.BuildContextHandler("request_id"),
		),
	)
	if err != nil {
		panic(err)
	}

	ctx := context.WithValue(context.Background(), "request_id", "12345")
	logger.WithContext(ctx).Info("This is an info message with context")
	logger.InfoContext(ctx, "This is an info message")
}
```

### Customizing Logger

```go
package main

import (
	"github.com/ace-zhaoy/glog"
	"go.uber.org/zap/zapcore"
	"os"
)

func main() {
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{}),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	logger := glog.NewLogger(core)
	logger.Info("This is a custom logger message")
}
```

## Testing

To run the tests, use the following command:

```sh
go test ./...
```

## License

This project is licensed under the MIT License.