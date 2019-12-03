/k\ui linux framebuffer driver

/start with
/ stdbuf -oL -eL showkey -a | k fb.k 
/ (unbuffer pipe, send an input line for each key stroke)
/quit with Cntrl-D

/load ui application (github.com/ktye/i/u/t.k)
\l t.k 

/this assumes 32 bit pixels        (cat /sys/class/graphics/fb0/bits_per_pixel)
(udw;udh):1920 1080 /screen size must match /sys/class/graphics/fb0/virtual_size

draw:{`"/dev/fb0"1:8_`@x}
draw us[udw;udh]
ky:{v:" "\:x;k:*|v;$[5<#v;;256>k:0+. k;draw uk[k;!0];]}
-1{ky@1:x;`}/`
