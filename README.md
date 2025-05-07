# slogflags

[![PkgGoDev](https://pkg.go.dev/badge/github.com/csmith/slogflags)](https://pkg.go.dev/github.com/csmith/slogflags)

Provides flags to configure the go structured logging package.

## Why?

How logs are formatted and what levels should be output depend on where and how
an application is deployed. For general purpose applications, they can't really
be set at compile time: some users will want JSON logs, others may want to
enable debug log levels to diagnose an issue, and so on.

The slog package doesn't provide a nice way of handling this, so `slogflags` was
created.

## How?

In the most basic case, import this package, call `flag.Parse()`, and then
retrieve your brand new logger by calling `slogflags.Logger()`:

```go
package main

import (
	"flag"
	"github.com/csmith/slogflags"
)

func main() {
	flag.Parse()
	l := slogflags.Logger()
	l.Warn("Danger", "user", "Will Robinson")
}
```

You can then run the app and specify `--log.level` (one of "debug", "info",
"warn" and "error") and `--log.format` (either "text" or "json").

## More advanced usage

You can pass options in to the `Logger` call to change the default behaviour.
See the godocs for a full list. The most useful one is perhaps `SetDefault`,
which sets the new logger as the default for the slog package. This also
enables bridging of calls from the `log` package.

```go
package main

import (
	"flag"
	"github.com/csmith/slogflags"
	"log"
)

func main() {
	flag.Parse()
	l := slogflags.Logger(slogflags.WithSetDefault(true))
	l.Warn("Danger", "user", "Will Robinson")
	log.Printf("I'll show up properly now too!")
}
```

## Licence/credits/contributions etc

Released under the MIT licence. See LICENCE for full details.

Contributions are welcome! Please feel free to open issues or send pull requests.