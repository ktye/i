(w;h):(12;12)
(r;g;b):|2 (256*)\1

/example image
f:{(y+(-/(*x)**x;2f**/*x);1+*|x)}              /mandelbrot step
l:{{{(4f>+/(*x)**x)&255>*|x}f[;y]/x}[(x;0);x]} /x y(`f) to n(0.255)
c:{+/(r*x;g*x;255-((x-127)*255*x-128)/16384)}  /colormap
m:{c@*|l x}                                    /pixel x y(`f) to color(`i)
d:0.01%w
x:d*!w;y:0.645-d*!h
p:,/x,/:\:y                                    /xy grid
i:m'p                                          /image(`i)

/sixel encoding of image i
rgb:{(g\x/r;g\x/g;g\x)}                        /`i to r g b
o:"" //o:0x1b50303b303b387122313b31            /sixel header
p:{$[255<#x;^?256#x;^?x]}i                     /palette used in image
i:p?i                                          /paletted image
o,:,/{"#",(*$!#x),";2;",";"/:$_(100%255)*rgb x}'p /write palette
o
{wc}'((_h%6);6;w)#i
