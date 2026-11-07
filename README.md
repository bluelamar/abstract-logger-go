# abstract-logger-go
Represents logging abstraction with typical logging methods.
It includes a default implementation that allows various output targets.
The default output is `stderr`, however any target that supports io.Writer
can be used.
The interface can be used for different user defined implementations.

The defined interface is LoggerI,
LoggerI provides leveled logging with configurable labels, time labels, and tags.
Set method implementations are expected to mutate the logger. 
With method implementations are expected to return a new configured logger without changing the original.
Fatal methods log at fatal severity but may or may not exit according to implementation.

## Installation

```sh
go get github.com/bluelamar/abstract-logger-go
```


## Usage

```go
package main

import (
    "os"

    "github.com/bluelamar/abstract-logger-go/alogger"
)

func main() {
    // Example using the default target `stderr` and date/time is *local*.
    // Output: `2022/12/31 10:35:32 [lcl] ERROR:[status][schedule] bad thing happened`
    //
    logger := alogger.New(nil, true).WithTag("status").WithTag("schedule")
    logger.Errorln("bad thing happened")

    // Example using the target `stdout` with date/time being *UTC* and *tag's* are added.
    // It also overrides the defaults for date/time output.
    // Output: `2022/12/31 10:04:42 :utc: ERROR:[status][schedule] bad thing happened`
    //
    logger := New(os.Stdout, false).WithTag("status").WithTag("schedule")
    logger.SetTimeLabels(":lcl:", ":utc:")
    logger.Errorln("bad thing happened")
}
```

