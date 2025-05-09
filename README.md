# maskjson

maskjson is a Go library for masking sensitive data in JSON. It allows you to easily hide sensitive information when marshaling structs to JSON by using struct tags.

## Installation

```bash
go get github.com/henry.chu/maskjson
```

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/henry.chu/maskjson"
)

type User struct {
    Username string `json:"username"`
    Password string `json:"password" mask:"true"` // This field will be masked
    Email    string `json:"email" mask:"true"`    // This field will be masked
}

func main() {
    user := User{
        Username: "johndoe",
        Password: "secret123",
        Email:    "john.doe@example.com",
    }
    
    // Create a new masker with default settings
    // Parameters: fullMask (bool), atLeastStar (uint)
    // - fullMask: if true, completely masks the value
    // - atLeastStar: minimum number of asterisks to use for masking
    masker := maskjson.NewMask(false, 3)
    
    // Marshal the struct with masked fields
    jsonData, err := masker.Marshal(user)
    if err != nil {
        panic(err)
    }
    
    fmt.Println(string(jsonData))
    // Output: {"username":"johndoe","password":"sec*****","email":"john**************"}
}
```

### Masking Options

The `NewMask` function accepts two parameters:

1. `fullMask bool`: When set to `true`, the entire value will be masked regardless of its length. When `false`, a portion of the beginning of the value will be visible.
2. `atLeastStar uint`: The minimum number of asterisks to use for masking. This ensures that sensitive fields are masked with at least this many asterisks.

## How Masking Works

When `fullMask` is set to `false`:
- For strings, a small portion at the beginning is kept visible, and the rest is masked with asterisks
- The number of visible characters depends on the string length and the `atLeastStar` value
- Non-string values are completely masked with asterisks

When `fullMask` is set to `true`:
- All values are completely masked with asterisks, regardless of their type or length

## License

[MIT License](LICENSE)
