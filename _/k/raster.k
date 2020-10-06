raster:{[wh;a]
 f:u!.'u:`strokeStyle`beginPath`moveTo`lineTo`clip`stroke`rect`lineWidth`textBaseline
 s:`f`a`o`w`h`ix`iy`fg`pa`cl!(f(!a);.a;(*/wh)#0;w;h;!w:wh 0;!h:wh 1;irgb 3#255;,/(0;0;wh))
 s:rasterF/:s
 s`o}

rasterF:{[s]f:*s`f;a:*s`a;s:$[#s`f;f[s;a];s];s[`f]:1_s`f;s[`a]:1_s`a;s}

strokeStyle:{[s;c];s[`fg]:icolor c;s}
beginPath:  {[s;u];moveTo[s;()]} /pa:(x0 y0;x1 y1;..)
moveTo:     {[s;u];s[`pa]:u;s}
lineTo:     {[s;u];s[`pa],:u;s}
clip:       {[s;u];s[`cl]:,/0 0-':(&/;|/)@\:s`p;s}
stroke:     {[s;u];f:iline[;;s`w;s`h];s[`o;,/1_f':s`pa]:s`fg;s}
rect:       {[s;u];dx:,/(u 2;0);dy:,/(0;u 3);o:u 0 1;s[`pa]:(o;o+dx;o+dx+dy;o+dy;o);s}

lineWidth:textBaseline:{y;x}


ihex:  {[c]+/i*16\*\(n:#i:|"0123456789abcdef"?-c)#16}
irgb:  {[rgb]+/65536 256 1*rgb}
icolor:{c:`white`black`gray!(irgb 3#255;0;irgb 3#128);
 $[(#c)>i:(!c)?k:`$x;c k;ihex 1_x]}

iline: {[a;b;w;h;cl]dx:*d:-a-b;dy:d 1;m:dx%0.+dy
 $[(+dx)<+dy;(w*_x)+y:_(a 1)+m*x:(a 0)+(1+2*dx<0)*!dx
  y+w*x:_(a 0)+m\y:(a 1)+(1+2*dy<0)*!dy]}

istr:{[x;y;w;s](x+w*y-20)+i+w*10\i:,/font@0+s} /font.k (10x20)

rasterT:`beginPath`strokeStyle`lineWidth`rect`stroke!(();"white";2;(10 10 80 50);())
rasterTest:{[](*wh) 'draster[wh:200 100;rasterT]}
