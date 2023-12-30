# EctoEnv

The ectoenv package is designed to simplify the process of binding environment variables to struct fields in Go. It uses struct tags to map environment variables to struct fields, allowing for easy configuration of your applications using environment variables.

## Installation

To use the ectoenv package, first install it using:

```bash Copy code
go get github.com/Gobusters/ectoenv
```

## Usage

To use ectoenv, import it into your Go file:

```go Copy code
import "github.com/Gobusters/ectoenv"
```

### Defining Structs

Define your configuration struct with the env and env-default struct tags to specify which environment variables should be bound to which struct fields. The env tag is used to specify the name of the environment variable, and env-default is used for a default value if the environment variable is not set.

Example:

```go Copy code
type Config struct {
    Host string `env:"HOST"`
    Port int `env:"PORT" env-default:"8080"`
    Debug bool `env:"DEBUG" env-default:"false"`
    Database string `env:"DATABASE_URL"`
}
```

### Using BindEnv

To bind environment variables to your struct, create an instance of your struct and pass a pointer to it to the BindEnv function.

Example:

```go Copy code
func main() {
    var cfg Config
    err := ectoenv.BindEnv(&cfg)
    if err != nil {
        log.Fatalf("Failed to bind environment variables: %v", err)
    }

    // Use cfg...
}
```

### Supported Types

The ectoenv package currently supports the following field types:

- `string`
- `int`
- `bool`
- `float64`
- Slices of the above types (e.g., `[]string`, `[]int`)
- Nested structs

### Error Handling

The BindEnv function will return an error if:

- The provided value is not a non-nil pointer to a struct.
- An environment variable is set with a value that cannot be converted to the field type.
- Any other reflection-related error occurs during the process.

## Using BindEnvWithAutoRefresh

BindEnvWithAutoRefresh extends the functionality of BindEnv by adding automatic refreshing of environment variables at a specified interval. This is particularly useful for long-running applications where environment variables might change over time.

### Functionality

`BindEnvWithAutoRefresh` sets the values of the provided struct based on the values of the environment variables defined in the struct's tags and periodically refreshes these values.

### Parameters

- `v`: A non-nil pointer to a struct.

### AUTO_REFRESH_INTERVAL

This variable sets the frequency that variables are refreshed. The default is 60 seconds.

### Usage

To use `BindEnvWithAutoRefresh`, pass your configuration struct and the desired refresh interval:

```go Copy code
func main() {
    var cfg Config
    ectoenv.AUTO_REFRESH_INTERVAL := 30 // refresh every 30 seconds (defaults to 60 seconds)

    err := ectoenv.BindEnvWithAutoRefresh(&cfg)
    if err != nil {
        log.Fatalf("Failed to bind and auto-refresh environment variables: %v", err)
    }

    // Your application logic...

}
```

## Contributing

Contributions to the ectoenv package are welcome! Please feel free to submit issues and pull requests to the repository.

## License

This package is released under the MIT License.
