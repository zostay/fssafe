# github.com/zostay/fssafe

[![GoDoc](https://godoc.org/github.com/zostay/fssafe?status.svg)](https://godoc.org/github.com/zostay/fssafe)

This is a small go package that provides a tools for saving files atomically. It also provides tooling for testing loading and saving of files using the same interface.

# Installation

There's nothing to install, but to add it to your project:

```bash
go get github.com/zostay/fssafe
```

# Usage

Then to use it in your code:

```go
package main

import (
	"fmt"
	"io"
	
	"github.com/zostay/fssafe"
)

func main() {
    loaderSaver := fssafe.NewFileSystemLoaderSaver("file.txt")

    // save a file, this will create a file named file.txt.new temporarily
    w, err := loaderSaver.Saver()
    if err != nil {
        panic(err)
    }

    fmt.Fprintln(w, "Hello, World!")
    w.Close()

    // If file.txt already exists, it will be renamed to file.txt.old.  The temporary
    // file.txt.new will be renamed to file.txt

    // load a file
    r, err := loaderSaver.Loader()
    if err != nil {
        panic(err)
    }

    data, err := io.ReadAll(r)
    if err != nil {
        panic(err)
    }

    fmt.Println(string(data))
    // Output: Hello, World!
}
```

## Testing

Another loader/saver is provided that can be created using 
`NewTestingLoaderSaver()`. The saver just writes to a memory buffer and the 
loader reads from that memory buffer.

```go
package main

import (
	"fmt"
	"io"
	
	"github.com/zostay/fssafe"
)

func main() {
    loaderSaver := fssafe.NewTestingLoaderSaver()

    // save to a memory buffer
    w, err := loaderSaver.Saver()
    if err != nil {
        panic(err)
    }

    fmt.Fprintln(w, "Hello, World!")
    w.Close()

    // load from the shared memory buffer
    r, err := loaderSaver.Loader()
    if err != nil {
        panic(err)
    }
    
    data, err := io.ReadAll(r)
    if err != nil {
        panic(err)
    }

    fmt.Println(string(data))
    // Output: Hello, World!
}
```

# Copyright & License

Copyright 2020 Andrew Sterling Hanenkamp

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
