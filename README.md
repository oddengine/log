# logger

Factory-Pattern Logger with DefaultWriter (file) in Go.

## Example

```xml
<Logger>
    <!-- Directory in which log files will be placed.                             -->
    <Directory>logs/</Directory>

    <!-- Log file name, which will be formated with datetime.                     -->
    <FileName>2006-01-02 15-04-05.000.log</FileName>

    <!-- Level could be one of "trace", "debug[N]", "info", "warn", "error".      -->
    <!-- N could be 0-7, and "debug0" equals to "debug".                          -->
    <Level>trace|debug1|info|warn|error</Level>

    <Rotation>
        <!-- Maximum file size in kilobytes (1024 bytes).                         -->
        <MaxSize>32768</MaxSize>

        <!-- Rotation Time, there are 2 types:                                    -->
        <!-- If type="daily", rotation only occurs every 24 hours, and the format -->
        <!-- is hh:mm:ss, for example 00:00:00 will rotates every midnight;       -->
        <!-- If type="duration", rotation occurs when the duration of the log     -->
        <!-- exceed a certain length. A duration string is a possibly signed      -->
        <!-- sequence of decimal numbers, each with optional fraction and a unit  -->
        <!-- suffix, such as "300ms", "1.5h" or "2h45m". Valid time units are     -->
        <!-- "ns", "us" (or "µs"), "ms", "s", "m", "h".                           -->
        <Schedule type="daily">00:00:00</Schedule>

        <!-- Max number of log files to keep.                                     -->
        <History>32</History>
    </Rotation>
</Logger>
```

```go
constraints := new(log.DefaultWriterConstraints)
constraints.Directory = "logs/"
constraints.FileName = "2006-01-02 15-04-05.000.log"
constraints.Level = "trace|debug1|info|warn|error"
constraints.Rotation.MaxSize = 32 * 1024 * 1024
constraints.Rotation.Schedule.Type = "daily"
constraints.Rotation.Schedule.Duration = "00:00:00"
constraints.Rotation.History = 32

factory := log.NewDefaultLoggerFactory(constraints)
logger := factory.NewLogger("CORE")

logger.Infof("Hello, %s", "world")
logger.Debugf(1, "Hello, %s", "world")
```
