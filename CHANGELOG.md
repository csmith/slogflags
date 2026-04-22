# Changelog

## 1.2.0 - 2026-04-22

### Other changes

* A warning is now logged if an invalid log level is specified.
  Previously this was a silent fallback to the default.

## 1.1.0 - 2025-05-17

### Features

* If you use the `AddCustomLevels` option, the level attribute of log statements
  will automatically be replaced to use the right names. Previously the default
  behaviour of slog was used, which was to use a built-in level and modifier
  e.g. `DEBUG-4`.

## 1.0.0 - 2025-05-07

_Initial release._
