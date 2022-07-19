[![Go Reference](https://pkg.go.dev/badge/github.com/joakimofv/find.svg)](https://pkg.go.dev/github.com/joakimofv/find/v2)

find
====

Utility functions for pattern matching. For use in tools that do text editing or command-line automation.

# Import

```go
    "github.com/joakimofv/find/v2"
```

# Version 2

Patterns given to `Replace` and `LongestFixedPart` now treat "\\\\" as "\\" and "\\\*" as "\*".
This makes it possible to match on asterixes, which before could not be done.
