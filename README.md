# go-logger

Structured logger for Uniplaces Go projects.

## Usage

#### Import
```go
import logger "github.com/uniplaces/go-logger"
```

#### Setup

Initializing the logger requires a config as parameter where you can 
set which is the current environment and the desired log level.

```go
logger.Init(logger.NewConfig("staging", "warning"))
```

Note: make sure you call `Init` only once and as early as possible in your program.

#### Use
Since the logger is implemented as a singleton, you can access it anywhere.

##### Simple use
```go
logger.Debug("debug message")
logger.Error(errors.New("error message"))
```

##### Structured use
```go
logger.
    Builder().
    AddField("test", "value").
    AddContextField("foo", "bar").
    Debug("debug message")
```

## Guidelines
This project's main focus is to bring a common standard to our Go projects's logs and, as such, 
please try to follow these guidelines when using this logger:

#### 1 - Add contextual information when logging an error

###### Use `AddContextField(key string, value interface{})`.

Example: when trying to fetch something from dynamodb and it fails.

In this case, contextual logging would include the ID of the object you're trying to fetch. 


#### 2 - Log as soon as possible

Logs should be specific and easy to trace back, so generic error logging should be avoided.
