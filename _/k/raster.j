raster:{[wh;a]w:wh 0;h:wh 1;o:(w*h)#0;r:fg:0;r:0;pa:();(.'!a)@'.a;o}

strokeStyle:{fg::icolor x;}
beginPath:  {[x]pa::();}
moveTo:     {pa::,_0.+x;}
lineTo:     {;pa::pa,,_0.+x;}
stroke:     {[x]f:iline[;;w];o[$[0~r;,/1_f':pa;&(r)>+iarc[r;pa;w;h]]]::fg;r::0;}
fill:       {[x]o[&0>iarc[r;pa;w;h]]::fg;r::0;}
rect:       {dx:,/(x 2;0);dy:,/(0;x 3);a:x 0 1;pa::(a;a+dx;a+dx+dy;a+dy;a);}
arc:        {pa::_x 0 1;r::_x 2;}
fillText:   {o[istr[x 1;x 2;w;x 0]]::fg;}
fillStyle:strokeStyle
clip:lineWidth:textBaseline:font:{x;}


ihex:  {[c]+/i*16\*\(n:#i:|"0123456789abcdef"?-c 4 5 2 3 0 1)#16}
irgb:  {[rgb]+/1 256 65536*rgb}
icolor:{c:`white`black`gray`red`green`blue!irgb'(255;0;128;255 0 0;0 255 0;0 0 255);$[(#c)>i:(!c)?k:`$x;c k;ihex 1_x]}
iline: {[a;b;w]$[(+*d)<+dy:(d:b-a)1;yline.$[0>dy;(b;a-b;w);(a;b-a;w)];xline.$[0>*d;(b;a-b;w);(a;b-a;w)]]}
xline: {[a;d;w]dx:*d;m:(d 1)%0.+dx;(x+a 0)+w*y:_(a 1)+m*x:!dx}
yline: {[a;d;w]dy:d 1;m:(*d)%0.+dy;(w*y+a 1)+x:_(a 0)+m*y:!dy}
iarc:  {[r;xy;w;h]`r`xy`w`h!(r;xy;w;h);x:(w/i:!w*h)-xy 0;y:(w\i)-xy 1;((x*x)+y*y)-r*r}
istr:  {[x;y;w;s],/(10*!#s)+(x+w*y-20)+(10/i)+w*10\i:glyphs@0+s} /font.k (10x20)
\font.k

rasterT:,/(`beginPath`strokeStyle`lineWidth`rect`stroke!(();"red";2;(10 10 80 50);())
 `fillText!,("text";10;30;40)
 `beginPath`strokeStyle`arc`stroke!(();"red";80 50 10 0 2p;())
 `beginPath`strokeStyle`arc`fill!(();"red";20 50 10 0 2p;()))
rasterTest:{[](*wh) 'draster[wh:200 100;rasterT]}

rasterT:`beginPath`lineWidth`strokeStyle`lineTo`lineTo`lineTo`lineTo`lineTo`lineTo`lineTo`stroke!(();2;"#8dd3c7";5.0 300.5379;25.52631 217.6333;46.05263 142.1343;66.57894 80.78503;87.10526 39.06561;107.6315 20.70275;128.1578 27.33673;())
/800 'draster[800 600;rasterT]
