```
web frontend (single-client, single-k, localhost:2019/k)

a.go main(server)
h.go html+js(term,edit,image)

files dropped to the browser are stored in-memory and are accessible with ioverbs.
reloading the page resets the k instance on the server.

term(left)
\c  clear terminal
\s  print stack trace (after errors)
ESC toggle hold mode(multiline)
RET execute current line or marked text
\eFILE edit(right)
.d:..  draw(right)

```
<img align="left" width="100" height="100" src="esd.png"/>

```
in sync(trigger)
.e(dit)      editor content "line1\nline2\ntext"
.s(election) selected text  ("text";12;15)
.d(isplay)   pixel buffer   h w#0
```

```
TODO:
 \lf (list files)
 download files? e.g. by url localhost:2019/filename
 search with mark and right click
```
