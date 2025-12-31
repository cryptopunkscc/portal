# Building runners

```shell
./mage build:apps
```

## MacOS M1
When running adding this on M1 MacOS is required:

```shell
export CGO_LDFLAGS="-framework UniformTypeIdentifiers"
```

Or

```go
/*
#cgo darwin LDFLAGS: -framework UniformTypeIdentifiers
*/

import "C"
```




