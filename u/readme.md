```
k\ui

the ui is written in k:
uk /key handler
um /mouse handler
us /(re)size handler

k receives an event, thinks and returns an updated screen.

many possible backends. only requirements:
- provide a screen (pixel array)
- send resize, key and mouse events to k
optionally
- clipboard (paste)
- drop files
```

# backend
```
possible backend implementations
pure k:                          file
 - web (via .Z.G or websocket?)  s.k
 - /dev/fb0.. (linux)
 - /dev/draw.. (plan9)
 - sixel (DEC, xterm..)
native application:              e.g. ./*.go
 - embed k in main or
 - connect to unmodified k binary
 
The go progam in this directory implements a backend based on gioui.org which should cover
- windows, X, wayland, mobile(android;apple)

The executable will do both: provide a built-in k, or connect to a server (k s.k t.k)
```

# demo application t.k (TODO)
```
The demo application provides a k terminal with a custom font (f/f3.k).
planned: plot built-in (x plot y) or (plot x) or (`bar plot y) with pure-k rasterizers.
```

# k interface
```
k\ui                                                      web(.Z.G)
uk[key;(shift;alt;cntrl)]                                /k,97,0,0,0
 key:printable ascii
     bs(8),tab(9),ret(13),esc(27)del(46),
     pageUp,pageDown,end,home,left,up,right,down(14..21)
um[button;(x0 x1);(y0;y1);(shift;alt;cntrl))]            /m,0,50,50,60,60,0,0,0
 button: left,middle,right,wheelUp,wheelDown 0..4
 x0 y0 x1 y1: press and release positions (no motion)
us[w;h]                                                  /s,1440,1080
 resize/layout..
 
events respond with nothing or a frame (row-major):
 (w*h)#0       / black
 (w*h)#255*256 / green.. (alpha byte is ignored)
```
