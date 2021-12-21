
b:256*g:256*r:255
d: (`rect;10;10;50;50;`co;r;`fill;`co;g;`stroke)
d,:(`co;b;`circ;70;20;20;`fill)
d,:(`lw;3;`co;g;`elli;70;70;20;10;`stroke)
d,:(`lw;2;`co;0;`poly;0 20 40 80 100.;20 80 30 70 10.;`stroke;`flush)
d,:(`lw;1;`tr;20;80;`circ;0;0;10;`stroke;`flush;`save)
d,:(`sc;2;2;`circ;0;0;10;`stroke;`flush)
d,:(`rest;`rect;0;0;10;10;`fill;`flush)
d,:(`ro;-30;`tr;20;0;`M;0;0;`L;20;0;`stroke;`flush)
d,:(`rest;`M;70;0;`h;20;`v;10;`stroke)
l:draw[d;100 100]
p:png l
`"draw.png"<p


d: (`M;10;80;`C;40;10;65;10;95;80;`S;150;150;180;80;`stroke)
d,:(`co;g;`M;80;80;`A;45;45;0;0;0;125;125;`L;125;80;`Z;`fill)
d,:(`co;r;`M;230;80;`A;45;45;0;1;0;275;125;`L;275;80;`Z;`fill)
d,:(`co;rgb 255 0 255;`M;80;230;`A;45;45;0;0;1;125;275;`L;125;230;`Z;`fill)
d,:(`co;b;`M;230;230;`A;45;45;0;1;1;275;275;`L;275;230;`Z;`fill)

l:draw[d;320 320]
p:png l
`"path.png"<p

`"arcs.png"<png svg@<`"img/arcs.svg"
