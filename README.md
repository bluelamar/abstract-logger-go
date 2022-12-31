# abstract-logger-go
Represents logging abstraction with typical logging methods.
It includes a default implementation that allows various output targets.
The default output is `stderr`, however any target that support io.Writer
can be used.
The interface can be used for different user defined implementations.

## Usage

Example using the default target `stderr` and date/time is *local*.
Output: `2022/12/31 10:35:32 [lcl] ERROR: bad thing happened`

```go
logger := New(nil, true)

logger.Errorln("bad thing happened")
```

Example using the target `stdout` with date/time being *UTC* and *tag's* are added.
Output: `2022/12/31 10:04:42 :utc: ERROR:[status][schedule] bad thing happened`

```go
logger := New(os.Stdout, false).WithTag("status").WithTag("schedule")

logger.Errorln("bad thing happened")
```

