```
input file: 10x20.png
glyphs are black on white, separated with 1 red pixel padding
for 10x20, width: 461 (41+42*10), height: 3+4*20 (83)

16 0123456789ABCDE (small double numbers for non-printable)
42 :+-*%&|<>=!~,^#_$?@.0123456789'/\;`"(){}[]
26 abcdefghijklmnopqrstuvwxyz
26 ABCDEFGHIJKLMNOPQRSTUVWXYZ

font.k is generated with
 go run gen.go 10x20.png 10 20 > font.k

/unpack (linear black-pixel indexes)
\l font.k
font:{&,/2\:'0+x}'font
```

![10x20](10x20.png)
