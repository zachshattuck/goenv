# `goenv`
Basic .env file parsing

## Usage:
```go
func main() {
  // Will process `.env`
  err := ProcessEnv()
  if err != nil {
		fmt.Print("Error setting environment: ")
		fmt.Println(err)
		return
  }

  myExpectedVar := os.Getenv("MY_EXPECTED_VAR")
  if myExpectedVar == "" {
		fmt.Print("MY_EXPECTED_VAR not set!")
		fmt.Println(err)
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
