# `goenv`
Basic .env file parsing

## Usage:
```go
import (
  "fmt"
  "github.com/zachshattuck/goenv"
)

func main() {
  // Will process `.env`
  err := goenv.ProcessEnv()
  if err != nil {
	fmt.Println("ProcessEnv: ", err)
	return
  }

  myExpectedVar := os.Getenv("MY_EXPECTED_VAR")
  if myExpectedVar == "" {
	fmt.Print("MY_EXPECTED_VAR not set")
	return
  }
}
```

## Drawbacks
- Whitespace and empty lines break it (yes I know this is dumb and bad)

## TODOs
- [ ] Specify other filenames
- [ ] Handle whitespace
- [ ] Should this populate a `goenv`-specific dictionary (e.g. `goenv.Get("VAR")`) instead of polluting the os environment?
