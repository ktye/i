```
k\ui

the ui is written in k:
.u.k /key   handler
.u.m /mouse handler
.u.s /(re)size handler

k receives an event, thinks and returns an updated screen.

many possible backends. only requirements:
- provide a screen (pixel array)
- send resize event
- send key event
- send mouse event
optionally
- clipboard (paste)
- drop files
```

#Backend
```
possible backend implementations
pure k:
 - web (via .Z.G?) built-in to ../a.go
 - /dev/fb0.. (linux)
 - /dev/draw.. (plan9)
 - sixel (DEC, xterm..)
native application:
 - embed k in main or
 - connect to k
 
This directory implements a backend based on gioui.org which should cover
- windows
- wayland
- maybe X
- mobile(android;apple)

The executable will do both: provide a built-in k, or connect to a server.
```

#Demo application
```
The demo application (which is built into ../a.go) provides a k terminal with a custom font (f/f3.k).
Planned: plot built-in (x plot y) or (plot x) or (`bar plot y) with pure-k rasterizers.
```

#Interface
```
k\ui                                                      web(.Z.G)
.u.k[key;(shift;alt;cntrl)]                               /k,97,0,0,0
 key:printable ascii
     bs(8),tab(9),ret(13),esc(27)del(46),
     pageUp,pageDown,end,home,left,up,right,down 14..21
.u.m[button;(x0;y0;x1;y1;(shift;alt;cntrl))]              /m,0,50,60,50,60,0,0,0
 button: left,middle,right,wheelUp,wheelDown 0..4
 x0 y0 x1 y1: press and release positions (no motion)
.u.s[w;h]                                                 /s,1440,1080
 resize/layout..
 
events respond with nothing or a frame (row-major):
 (w*h)#0       / black
 (w*h)#255*256 / green.. (alpha byte is ignored)
```
