raster:{[wh;a]
 s:`o`w;`h;`ix`iy`fg`pa`cl!((*/wh)#0;w;h;!w:wh 0;!h:wh 1;irgb 3#255;,/(0;0;wh))
 s`o}

strokeStyle:{[s;c]s[`fg]:icolor c;s}
beginPath:  {[s;u]moveTo[s;()]} /pa:(x0 y0;x1 y1;..)
moveTo:     {[s;u]s[`pa]:u;s}
lineTo:     {[s;u]s[`pa],:u;s}
moveTo:     {[s;u]s[`pa]:u;s}
clip:       {[s;u]s[`cl]:,/0 0-':(&/;|/)@\:s`p;s}
stroke:     {[s;u]f:iline[;;s`w;s`h];s[`o;,/1_f':s`pa]:s`fg;s}

lineWidth:textBaseline:{y;x}


ihex:  {[c]+/i*16\*\(n:#i:|"0123456789abcdef"?-c)#16}
irgb:  {[rgb]+/65536 256 1*rgb}
icolor:{c:`white`black`gray!(irgb 3#255;0;irgb 3#128);
 $[(#c)>i:(!c)?k:`$x;c k;ihex 1_x]}

iline: {[a;b;w;h;cl]dx:*d:-a-b;dy:d 1;m:dx%0.+dy
 $[(+dx)<+dy;(w*_x)+y:_(a 1)+m*x:(a 0)+(1+2*dx<0)*!dx
  y+w*x:_(a 0)+m\y:(a 1)+(1+2*dy<0)*!dy]}

