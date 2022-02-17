/plot `x`y!(xF;yF)     line with titles x and y
/plot `x`y!(xF;yL)     multiple lines, individual axes(xL;yL)
/plot `t!,xZ           polar plot
/plot 3 4 5!10 1 22    bar plot xI!yI

PW:800;PH:600;FH:20;FONT:"monospace"

`"pltco":{[p;x;y]w:p`fh;h:p`fh;X:p`px;Y:p`py;W:p`pw;H:p`ph; d:(X+w;X+W-w;Y+H-h;Y+h);a:p`a; ((d 0 1)(a 0 1)'x;(d 2 3)(a 2 3)'y)}
`"pltcl":{[x;y]p:.`"pltco";xy:p[plt;x;y]; `<"x:",($xy 0), " y:",($xy 1),_10 32}
`"pltzo":{[x;y;w;h]p:.`"pltco";xy0:p[plt;x;y+h]; xy1:p[plt;x+w;y]; plt[`a]:(xy0,xy1)0 2 1 3; plt[`t]:$[`polar~t:plt`t;`xy;t]; draw[`Plot@`plt;(PW;PH)]}

plot:{plt::`plot x;Show[draw[`Plot@`plt;(PW;PH)];.`"pltcl";.`"pltzo"]}

`"plot.":{[d];l:$!d;v:.d; t:$[2~#d;`xy;`polar];
 y:$[t~`xy; $[`L~@y:v 1;y;,y];          $[`L~@y:_*v;y;,y]]
 x:$[t~`xy; $[`L~@x:v 0;x;(,x)@(#y)#0]; $[`L~@x:imag@*v;x;,x]]
 xt:`tics(&/&/x;|/|/x);yt:`tics(&/&/y;|/|/y)
 a:$[t~`xy;(xt 0;*-1#xt;yt 0;*-1#yt);(-a;a;-a;a:*|`tics@0.,|/|/abs@*v)]
 c:c@(#c:11826975 950271 2924588 2631638 12412820 4937356 12744675 8355711 2276796 13614615)/!#x
 style:$[t~`polar;"..";`i~@**y;"||";"--"] / -.| line points bar
 size: $[t~`polar;2;style~"||";(--/((**x),-1#*x))%-1+#*x ;2]
 lines:{`style`size`color`x`y!(style;size;z;x;0.+y)}'[x;y;c]
 `L`t`l`a`f`fw`fh`px`py`pw`ph!(lines;t;l;a;FONT;FW;FH;50;30;$[t~`xy;PW-50;PH-30];PH-30)}

`"Plot.":{[sym];x:.sym;w:x`fh; h:x`fh; X:x`px; Y:x`py; W:x`pw; H:x`ph; a:x`a;T:"title";grey:13882323
 dst:(X+w;X+W-w;Y+H-h;Y+h);rdst:(X+w;Y+h;W-2*w;H-2*h)
 xs:(a 0 1)(dst 0 1)' /transform axis to canvas
 ys:(a 2 3)(dst 2 3)'
 d:{(x i)!x 1+i:2*!2\#x}
 bars:{[l]$["|"':l`style;(`color;l`color),,/{(`Rect;((-dx%2)+xs x;ys y;dx:-/xs(l`size;0.);(ys a 2)-ys y))}'[l`x;l`y];()]}
 line:{[l]$["-"':l`style;(`linewidth;l`size;`color;l`color;`poly;(xs l`x;ys l`y));()]}
 dots:{[l]$["."':l`style;(`color;l`color),,/{(`Circle;(xs x;ys y;1.5*l`size))}'[l`x;l`y];()]}
 c:(`clip;(X;Y;W;H);`font;(x`f;x`fh);`text;((X+W%2;Y+h);1;T))
 xy:{[]c,:(`text;((X+w;Y+H);0;$a 0);`text;((X+W%2;Y+H);1;(x`l)0);`text;((X+W-w;Y+H);2;$a 1))
       c,:(`Text;((X+w;Y+H-h);0;$a 2);`Text;((X+w;Y+H%2);2;(x`l)1);`Text;((X+w;Y+h);2;$a 3))
       c,:(`color;0;`linewidth;2;`rect;rdst)      /todo: clip rdst
       c,:(`linewidth;1;`color;grey)
       c,:(`clip;rdst)
       c,:,/{(`line;0.+(x;dst 2;x;dst 3))}'xs`tics x[`a;0 1]
       c,:,/{(`line;0.+(dst 0;x;dst 1;x))}'ys`tics x[`a;2 3]}
 po:{[]C:(X+W%2;Y+H%2);R:(H%2)-h
       c,:(`text;((C 0;Y+H);1;(x`l)0);`text;(C+.75*R;6;$(x`a)3))
       c,:,/{(`text;(C+R*(_;imag)@'x;y;z))}'[1@270.+a;0 0 6 6 4 4 2 2;$a:30 60 120 150 210 240 300 330]
       c,:(`color;0),/{(`line;,/+C+(R-w%2;R)*/:(_;imag)@'x)}'1@30.*!12
       /c,:(`clip;C,R) /bug in cairo?
       c,:(`color;grey;`linewidth;1;`line;(X+w;C 1;X+W-h;C 1);`line;(*C;Y+h;*C;Y+H-h))
       c,:,/{(`circle;0.+C,x)}'r:(xs@`tics 0.,x[`a;3])-*C
       c,:(`color;0;`linewidth;2;`circle;C,R)}
       
 $[`xy~x`t;xy[];po[]]
 c,:,/bars'x`L
 c,:,/line'x`L
 c,:,/dots'x`L}
 

xy:`alpha`beta!(0.+!10;(sin 0.2*!10;cos 0.2*!10))
po:`abc!,1a30+.1*?10a
ba:`x`y!(1.+!#i;i:68 56 42 29 0 39 8 61)

