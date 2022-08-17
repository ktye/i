
echo -n "const plotdict_=\""
cat << 'EOF' | sed -e 's/\\/\\\\/g' -e 's/"/\\"/g' -e 's/$/\\n\\/g'
{[d]l:$!d;v:.d; t:$[2~#d;`xy;`polar];
 y:$[t~`xy; $[`L~@y:v 1;y;,y];          $[`L~@y:_*v;y;,y]]
 x:$[t~`xy; $[`L~@x:v 0;x;(,x)@(#y)#0]; $[`L~@x:imag@*v;x;,x]]
 xt:`tics(&/&/x;|/|/x);yt:`tics(&/&/y;|/|/y)
 a:$[t~`xy;(xt 0;*-1#xt;yt 0;*-1#yt);(-a;a;-a;a:*|`tics@0.,|/|/abs@*v)]
 c:c@(#c:11826975 950271 2924588 2631638 12412820 4937356 12744675 8355711 2276796 13614615)/!#x
 style:$[t~`polar;"..";`i~@**y;"||";"--"] / -.| line points bar
 size: $[t~`polar;2;style~"||";(--/((**x),-1#*x))%-1+#*x ;2]
 lines:{`style`size`color`x`y!(style;size;z;x;0.+y)}'[x;y;c]
 `L`T`t`l`a!(lines;"";t;l;a)}
EOF


echo -n "const plotwh_=\""
cat << 'EOF' | sed -e 's/\\/\\\\/g' -e 's/"/\\"/g' -e 's/$/\\n\\/g'
{[x;fs;W;H]; ; a:x`a;T:x`T;grey:13882323
 C:(W%2;H%2);R:(W%2)&(H%2)-fs
 dst:$[`xy~x`t;(fs;W-fs;H-fs;fs);((C-R),C+R)0 2 3 1];rdst:(fs;fs;W-2*fs;H-2*fs)
 xs:(a 0 1)(dst 0 1)' /transform axis to canvas
 ys:(a 2 3)(dst 2 3)'
 bars:{[l]$["|"':l`style;(`color;l`color),,/{(`Rect;((-dx%2)+xs x;ys y;dx:-/xs(l`size;0.);(ys a 2)-ys y))}'[l`x;l`y];()]}
 line:{[l]$["-"':l`style;(`linewidth;l`size;`color;l`color;`poly;(xs l`x;ys l`y));()]}
 dots:{[l]$["."':l`style;(`color;l`color),,/{(`Circle;(xs x;ys y;1.5*l`size))}'[l`x;l`y];()]}
 c:(`clip;(0;0;W;H);`font;("monospace";fs);`color;0;`text;((W%2;fs);1;T))
 xy:{[]c,:(`text;((fs;H);0;$a 0);`text;((W%2;H);1;(x`l)0);`text;((W-fs;H);2;$a 1))
       c,:(`Text;((fs;H-fs);0;$a 2);`Text;((fs;H%2);2;(x`l)1);`Text;((fs;fs);2;$a 3))
       c,:(`color;0;`linewidth;2;`rect;rdst)      /todo: clip rdst
       c,:(`linewidth;1;`color;grey)
       c,:(`clip;rdst)
       c,:,/{(`line;0.+(x;dst 2;x;dst 3))}'xs`tics x[`a;0 1]
       c,:,/{(`line;0.+(dst 0;x;dst 1;x))}'ys`tics x[`a;2 3]}
 po:{[]c,:(`text;((C 0;H);1;(x`l)0);`text;(C+.75*R;6;$(x`a)3))
       c,:(`font;("monospace";_fs*.8)),,/{(`text;(C+R*(_;imag)@'x;y;z))}'[1@270.+a;0 0 6 6 4 4 2 2;$a:30 60 120 150 210 240 300 330]
       c,:(`color;0),/{(`line;,/+C+(R-fs%2;R)*/:(_;imag)@'x)}'1@30.*!12
       /c,:(`clip;C,R) /bug in cairo?
       c,:(`color;grey;`linewidth;1;`line;((-R)+*C;C 1;R+*C;C 1);`line;(*C;(-R)+C 1;*C;R+C 1))
       c,:,/{(`circle;0.+C,x)}'r:(xs@`tics 0.,x[`a;3])-*C
       c,:(`color;0;`linewidth;2;`circle;C,R)}
 $[`xy~x`t;xy[];po[]]
 c,:,/bars'x`L
 c,:,/line'x`L
 c,:,/dots'x`L}
EOF
