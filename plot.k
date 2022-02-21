/plot `x`y!(xF;yF)     line with titles x and y
/plot `x`y!(xF;yL)     multiple lines, individual axes(xL;yL)
/plot `t!,xZ           polar plot
/plot 3 4 5!10 1 22    bar plot xI!yI

PW:800;PH:600;FH:20;FONT:"monospace"

`"pltnn.":{wi:</</:;$[#i:&wi[*x;+\plt`px`pw]&wi[x 1;+\plt`py`ph];*i;0]}
`"pltco.":{[p;x;y]w:p`fh;h:p`fh;X:p`px;Y:p`py;W:p`pw;H:p`ph;C:(X+W%2;Y+H%2);R:(W%2)&(H%2)-h;d:$[`xy~p`t;(X+w;X+W-w;Y+H-h;Y+h);((C-R),C+R)0 2 3 1];a:p`a; ((d 0 1)(a 0 1)'x;(d 2 3)(a 2 3)'y)}
`"pltcl.":{[x;y]p:.`"pltco.";n:`pltnn(x;y); xy:p[plt n;x;y]; `<$[`polar~plt[n;`t];$imag/|xy;"x:",($xy 0), " y:",($xy 1)],_10 32}
`"pltzo.":{[x;y;w;h]n:`pltnn(x;y);p:.`"pltco.";xy0:p[plt n;x;y+h]; xy1:p[plt n;x+w;y]; plt[`a;n]:(xy0,xy1)0 2 1 3; plt[`t;n]:`xy; draw[`plts@`plt;(PW;PH)]}

plot:{n:#x;multi:{x[`pw]%:n;x[`px]:(!n)*x`pw;x};plt::$[(@x)':`L`T;multi@(`plot@)'x;,`plot x];Show[draw[`plts@`plt;(PW;PH)];.`"pltcl.";.`"pltzo."]}

`"plot.":{[d]l:$!d;v:.d; t:$[2~#d;`xy;`polar];
 y:$[t~`xy; $[`L~@y:v 1;y;,y];          $[`L~@y:_*v;y;,y]]
 x:$[t~`xy; $[`L~@x:v 0;x;(,x)@(#y)#0]; $[`L~@x:imag@*v;x;,x]]
 xt:`tics(&/&/x;|/|/x);yt:`tics(&/&/y;|/|/y)
 a:$[t~`xy;(xt 0;*-1#xt;yt 0;*-1#yt);(-a;a;-a;a:*|`tics@0.,|/|/abs@*v)]
 c:c@(#c:11826975 950271 2924588 2631638 12412820 4937356 12744675 8355711 2276796 13614615)/!#x
 style:$[t~`polar;"..";`i~@**y;"||";"--"] / -.| line points bar
 size: $[t~`polar;2;style~"||";(--/((**x),-1#*x))%-1+#*x ;2]
 lines:{`style`size`color`x`y!(style;size;z;x;0.+y)}'[x;y;c]
 pw:PW;ph:PH;`L`T`t`l`a`f`fh`px`py`pw`ph!(lines;"";t;l;a;FONT;FH;0;0;pw;ph)}
 
`"plts.":{[sym];x:.sym;$[`D~@x;`Plot x;,/(`Plot@)'x]}

`"Plot.":{[x];w:x`fh; h:x`fh; X:x`px; Y:x`py; W:x`pw; H:x`ph; a:x`a;T:x`T;grey:13882323
 C:(X+W%2;Y+H%2);R:(W%2)&(H%2)-h
 dst:$[`xy~x`t;(X+w;X+W-w;Y+H-h;Y+h);((C-R),C+R)0 2 3 1];rdst:(X+w;Y+h;W-2*w;H-2*h)
 xs:(a 0 1)(dst 0 1)' /transform axis to canvas
 ys:(a 2 3)(dst 2 3)'
 bars:{[l]$["|"':l`style;(`color;l`color),,/{(`Rect;((-dx%2)+xs x;ys y;dx:-/xs(l`size;0.);(ys a 2)-ys y))}'[l`x;l`y];()]}
 line:{[l]$["-"':l`style;(`linewidth;l`size;`color;l`color;`poly;(xs l`x;ys l`y));()]}
 dots:{[l]$["."':l`style;(`color;l`color),,/{(`Circle;(xs x;ys y;1.5*l`size))}'[l`x;l`y];()]}
 c:(`clip;(X;Y;W;H);`font;(x`f;x`fh);`color;0;`text;((X+W%2;Y+h);1;T))
 xy:{[]c,:(`text;((X+w;Y+H);0;$a 0);`text;((X+W%2;Y+H);1;(x`l)0);`text;((X+W-w;Y+H);2;$a 1))
       c,:(`Text;((X+w;Y+H-h);0;$a 2);`Text;((X+w;Y+H%2);2;(x`l)1);`Text;((X+w;Y+h);2;$a 3))
       c,:(`color;0;`linewidth;2;`rect;rdst)      /todo: clip rdst
       c,:(`linewidth;1;`color;grey)
       c,:(`clip;rdst)
       c,:,/{(`line;0.+(x;dst 2;x;dst 3))}'xs`tics x[`a;0 1]
       c,:,/{(`line;0.+(dst 0;x;dst 1;x))}'ys`tics x[`a;2 3]}
 po:{[]c,:(`text;((C 0;Y+H);1;(x`l)0);`text;(C+.75*R;6;$(x`a)3))
       c,:(`font;(x`f;_h*.8)),,/{(`text;(C+R*(_;imag)@'x;y;z))}'[1@270.+a;0 0 6 6 4 4 2 2;$a:30 60 120 150 210 240 300 330]
       c,:(`color;0),/{(`line;,/+C+(R-w%2;R)*/:(_;imag)@'x)}'1@30.*!12
       /c,:(`clip;C,R) /bug in cairo?
       c,:(`color;grey;`linewidth;1;`line;((-R)+*C;C 1;R+*C;C 1);`line;(*C;(-R)+C 1;*C;R+C 1))
       c,:,/{(`circle;0.+C,x)}'r:(xs@`tics 0.,x[`a;3])-*C
       c,:(`color;0;`linewidth;2;`circle;C,R)}
 $[`xy~x`t;xy[];po[]]
 c,:,/bars'x`L
 c,:,/line'x`L
 c,:,/dots'x`L}
 

xy:`alpha`beta!(0.+!10;(sin 0.2*!10;cos 0.2*!10))
po:`abc!,1a30+.1*?10a
ba:`x`y!(1.+!#i;i:68 56 42 29 0 39 8 61)

/plot xy
/plot(xy;po;ba)
