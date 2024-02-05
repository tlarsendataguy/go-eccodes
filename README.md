# go-eccodes
Go wrapper for [ecCodes](https://confluence.ecmwf.int/display/ECC)

Forked from [amsokol/go-eccodes](https://github.com/amsokol/go-eccodes) and modified in the following ways:

- Remove dependency to jasper
- Added Mac build
- Converted to Go Modules
- Changed Go function names to camel casing

### Install

`go get github.com/tlarsendataguy/go-eccodes`

You will need to install eccodes in your machine. This module is only a wrapper over the C bindings of eccodes, it is not a pure Go solution.

If you use a package manager to install eccodes, be sure to set `CGO_CFLAGS -I` to the path with the header files and `CGO_LDFLAGS -L` to the path containing the libraries. On an M2 Mac, that might look like the following:
```
CGO_CFLAGS -I/opt/homebrew/include
CGO_LDFLAGS -L/opt/homebrew/lib
```

This package has been tested with ecCodes version 2.33.0 on an M2 Mac using Homebrew.
