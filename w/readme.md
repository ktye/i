```
web frontend (single-client, single-k, localhost:2019/k)

a.go main(server)
h.go html+js

files dropped to the browser are stored in-memory and can be read by k.
k cannot touch files on disk.
reloading the page resets the k instance on the server.

\c  clear terminal
\s  print stack trace (after errors)
ESC toggle hold mode(multiline)
RET execute current line or marked text
```
<img align="left" width="100" height="100" src="esd.png"/>

```
.e(edit)       editor content    ("line1";"line2";..)
.s(selection)  selected text     "text"
.d(display)    pixel buffer `i   w h#!w*h
```

```
TODO:
 write files (from k to memfs)
 \lf (list files)
 download files? e.g. by url localhost:2019/filename
 search with mark and right click
```
