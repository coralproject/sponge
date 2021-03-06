
# log
    import "github.com/ardanlabs/kit/log"

Package log provides a simple layer above the standard library logging package.
The base API provides two logging levels, DEV and USER. The DEV level logs things
developers need and can be verbose. The USER level logs things for users need
and should not be verbose. There is an Error call which falls under USER.

To initialize the logging system from your application, call Init:


	logLevel := func() int {
		ll, err := cfg.Int("LOGGING_LEVEL")
		if err != nil {
			return log.DEV
		}
		return ll
	}
	
	log.Init(os.Stderr, logLevel)

To write to the log you can make calls like this:


	log.Dev(context, "CreateUser", "Started : Email[%s]", nu.Email)
	log.Error(context, "CreateUser", err, "Completed")

The API for Dev and User follow this convention:


	log.User(context, "funcName", "formatted message %s", values)

context

A context is a unique id that can be used to trace an entire session or
request. This value should be generated as early as possible and passed
through out the different calls.

funcName

Provide the name of the function the log statement is being declared in. This
can take on a type name in the case of the method (type.method).

formatted message, values

Any string can be provided but it does support a formatted message. Values
would be substituted if provided. This messaging is up to you.




## Constants
``` go
const (
    NONE int = iota
    DEV
    USER
)
```
Level constants that define the supported usable LogLevel.



## func Dev
``` go
func Dev(context interface{}, funcName string, format string, a ...interface{})
```
Dev logs trace information for developers.


## func Error
``` go
func Error(context interface{}, funcName string, err error, format string, a ...interface{})
```
Error logs trace information that are errors.


## func Fatal
``` go
func Fatal(context interface{}, funcName string, format string, a ...interface{})
```
Fatal logs trace information for users and terminates the app.


## func Init
``` go
func Init(w io.Writer, level func() int) error
```
Init must be called to initialize the logging system. This function should
only be called once.


## func User
``` go
func User(context interface{}, funcName string, format string, a ...interface{})
```
User logs trace information for users.









- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)