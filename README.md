# `log`

Event logging package with configurable output sinks. Supported log levels: `Trace < Info < Error`. 

## Install

```sh
go get github.com/aj0strow/log
```

There are no dependencies outside of the standard library. You can also copy and paste the source code into `project/log` package. 

## Add Sinks

Start by creating a new logger. You should do this early in your command `main()`. 

```go
    logger := log.New()
```

Right now `logger` has no sinks; new messages are effectively ignored. Add sinks to specify where `logger` should write. 

```go
    if production {
        // In production, only log.Info messages or higher will be written
        // to process stderr. 
        logger.AddSink(log.Info, &log.Std{})
        
        // ErrorReporter must implement log.Sink
        logger.AddSink(log.Error, &ErrorReporter{})
    } else {
        logger.AddSink(log.Trace, &log.Std{})
    }
```

In the example, `ErrorReporter` is a custom sink. Learn about custom sinks below. 

## Write Events

Log messages include `Level`, `Time`, `Message`. I use `Level = Trace` for debugging the implementation, `Level = Info` for business events, and `Level = Error` for errors. 

You can write new messages using the verbose `Append` function. 

```go
    logger.Append(&log.Message{
        Level: log.Info,
        Time: time.Now(),
        Message: "Send complete log message events.",
    })
```

You can also use the following convenience functions. Format strings and arguments use `fmt.Sprintf`. 

```go
    logger.Appendf(log.Info, "readme %s progress", "aj0strow/log")
    
    logger.Trace("shortcut to Appendf with %s level", "TRACE")
    
    logger.Info("just like %s but at %s level", "TRACE", "INFO")
    
    logger.Error(fmt.Errorf("accepts an error for convenience"))
    logger.Errorf("you can pass %s format strings too", "ERROR")
    
    logger.Fatal(fmt.Errorf("logs error and then non-zero exit"))
    logger.Fatalf("logs error and then exits with %d", 1)
```

Sinks are responsible for formatting the structured messages, for example prepending time stamp and server info. 

## Time Source

Logger convenience functions set the time automatically. The default time source uses `time.Now()`. You can change the time source by implementing the interface.

```go
type TimeSource interface {
	  Now() time.Time
}
```

For example to use `UTC` time.

```go
type UTCTime struct {}

func (*UTCTime) Now() time.Time {
    return time.Now().UTC()
}

var _ log.TimeSource = (*UTCTime)(nil)

// Change time source
logger.TimeSource = &UTCTime{}
```

It's better to leave truncating time (if necessary) to log sinks. 

## Custom Sinks

The only sink included is `log.Std` which writes messages to `os.Stderr`. Custom log sinks need to implement the `log.Sink` interface. 

```go
type Sink interface {
	Append(*Message) error
}
```

For example, to write to `os.Stderr` with colors. 

```go
import (
    "github.com/aj0strow/log"
    "github.com/fatih/color"
)

var (
    red  = color.New(color.FgRed)
    blue = color.New(color.FgBlue)
    gray = color.New(color.FgBlack, color.Faint)
)

type StdWithColor struct {
    Std *log.Std
}

func (swc *StdWithColor) Append(msg *log.Message) error {
    switch msg.Level {
    case log.Error:
        swc.Std.Write(red.Sprintf(msg.Message))
    case log.Info:
        swc.Std.Write(blue.Sprintf(msg.Message))
    default:
        swc.Std.Write(gray.Sprintf(msg.Message))
    }
    return nil
}

var _ log.Sink = (*StdWithColor)(nil)
```

Add the sink to your logger.

```go
    logger.AddSink(log.Trace, &StdWithColor{
        Std: &log.Std{},
    })
```
