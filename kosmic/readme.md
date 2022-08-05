# kosmic

[kosmic.com](https://github.com/ktye/i/releases/download/latest/kosmic.com)
is a build of `k.c` with a webserver.

The binary runs on linux/mac/windows/bsds thanks to [cosmopolitan](https://justine.lol/cosmopolitan/index.html).

The c side of the webserver is minimalistic, single-threaded and forwards every request to k.

# arguments

```
kosmic [file.k *] [1234] [-e][kstr]
 files.k  are executed
 int   is used as a port
 -e    terminates (no repl)
 -e X  evaluates X before terminating
```

- without arguments (e.g. double-click) executes `a.k`, if present in the same directory
  and listens on `port` (a k variable) if it is non-zero. otherwise the server is not started.
  
- without explicit `-e`, the terminal always remains a k-console.

# minimal web application
- k has a variable `port` (a nonzero integer)
- k has a function `serve` that receives the request as `C` and responds with `C`

# applications

- the example application `a.k` handles GET requests:
- requests starting with `/?` are evaluated as k expressions, e.g. `GET /?fn[1;2] HTTP/1.1`
- other get requests return files relative to the current directory

# install/run
- copy kosmic.com anywhere (e.g. a directory of data files)
- add/edit a.k, html/js files for a front-end
- double-clock kosmic.com and open a web browser with the address written in the terminal.

# c-extensions

`kosmic` by itself can only execute k functions (e.g. read and write files, ktye/k does not include fork/exec).

you can add additional c code to `kosmic.c` and make it available to k by the [c-api](../+/k.h)

as an example a `clock` function is included.

see (mk)[mk] how to compile.




