/plot `x`y!(xF;yF)     line with titles x and y
/plot `x`y!(xF;yL)     multiple lines, individual axes(xL;yL)
/plot `t!,xZ           polar plot
/plot 3 4 5!10 1 22    bar plot xI!yI


PW:600;PH:400

plot:{plt::`plot x;Plot`plt;s:`c`k`z`d!"plt.",/:"ckzd"
 ui::"<canvas width=",($PW)," height=",($PH)," id='",(s`c),"' data-click='",(s`k),"' data-zoom='",(s`z),"'></canvas><div id='",(s`d),"'></div>"}

`"plot.":{[d];l:$!d;v:.d; t:$[2~#d;`xy;`polar];
 y:$[t~`xy; $[`L~@y:v 1;y;,y];          $[`L~@y:_*v;y;,y]]
 x:$[t~`xy; $[`L~@x:v 0;x;(,x)@(#y)#0]; $[`L~@x:imag@*v;x;,x]]
 xt:`tics(&/&/x;|/|/x);yt:`tics(&/&/y;|/|/y)
 a:$[t~`xy;(xt 0;*-1#xt;yt 0;*-1#yt);(-a;a;-a;a:*|`tics@0.,|/|/abs@*v)]
 c:c@(#c:","\:"#1f77b4,#ff7f0e,#2ca02c,#d62728,#9467bd,#8c564b,#e377c2,#7f7f7f,#bcbd22,#17becf")/!#x
 style:$[t~`xy;"--";".."] / -.| line points bar
 size: $[t~`xy;2;5]
 lines:{`style`size`color`x`y!(style;size;z;x;y)}'[x;y;c]
 `L`t`l`a`f`fw`fh`pw`ph!(lines;t;l;a;"20px monospace";FW;FH;$[t~`xy;PW;PH];PH)}


Plot:{[sym];x:.sym; w:x`fw; h:x`fh; W:x`pw; H:x`ph; a:x`a
 dst:(w;W-w;H-h;h);rdst:(w;h;W-2*w;H-2*h)
 xs:(a 0 1)(dst 0 1)' /transform axis to canvas
 ys:(a 2 3)(dst 2 3)'
 d:{(x i)!x 1+i:2*!2\#x}
 line:{[l](`lineWidth;l`size;`strokeStyle;l`color;`beginPath;()),(,/(`moveTo,(-1+#l`x)#`lineTo),','+(xs l`x;ys l`y)),(`stroke;())}
 c:`font`textBaseline!(x`f;"bottom")
 xy:{[]c,:`fillText!(($a 0;w;H);(s;2\W-w*#s:(x`l)0;PH);(s;W-w*1+#s:$a 1;H))
       c,:`fillText!((("[",(", "/:$a 2 3),"]");w;h);(s;2\W-w*#s:(x`l)1;h))
       c,:`rect`strokeStyle`lineWidth`stroke`clip!(rdst;"black";2;();())
       c,:`strokeStyle`lineWidth!("lightgrey";1)
       c,:d@,/{(`moveTo;(x;dst 2);`lineTo;(x;dst 3);`stroke;())}'xs`tics x[`a;0 1]
       c,:d@,/{(`moveTo;(dst 0;x);`lineTo;(dst 1;x);`stroke;())}'ys`tics x[`a;2 3]}
 po:{[]c,:`fillText!((s;2\W-w*#s:(x`l)0;h);($(x`a)3;W%2;H))
       c,:`strokeStyle`lineWidth`beginPath`arc`stroke`clip!("black";2;();(W%2;H%2;(H%2)-h;0;2p);();())
       c,:`strokeStyle`lineWidth`moveTo`lineTo`stroke`moveTo`lineTo`stroke!("lightgrey";1;(h;H%2);(W-h;H%2);();(W%2;h);(W%2;H-h);())
       c,:d@,/{(`beginPath;();`arc;(W%2;H%2;x;0;2p);`stroke;())}'(xs`tics 0.,x[`a;3])-W%2
 };$[`xy~x`t;xy[];po[]]
 
 c,:d@,/line'x`L
 back:{[a;dst;x;y]((dst 0 1)(a 0 1)'x;(dst 2 3)(a 2 3)'y)}[a;dst]  /transform canvas to axis
 click:{[b;s;t;x;y];(`0)[s;$[t~`xy;"\n"/:("x:";"y:"),'$b[x;y]; $imag/|b[x;y]]]}[back;`$($sym),".d";x`t]
 zoom:{[b;s;x]a:b[x 0;x 1],b[x 2;x 3];.[s;(`a;0 2 1 3);@[(^a 0 2),^a 1 3;0 2 1 3]];.[s;`t;`xy];Plot s}[back;sym]
 (`0)[`$($sym),".c";c]     /assign to dynamic symbol
 (`0)[`$($sym),".k";click]
 (`0)[`$($sym),".d";""]
 (`0)[`$($sym),".z";zoom]
}



`"tics.":{[minmax]nice:{[x;r]f:x%10^ex:_log[10;x];(1 2 5 10.@*&(f<1.5 3 7 0w;~f>1 2 5 0w)[r])*10^ex}
 e:nice[-/|minmax;0];s:nice[e%4.;1];n:_1.5+e%s;s*(_(*minmax)%s)+!n}

x:(0 9)(0 4p)'!10
y:(sin x;cos x)
/p:plot `xtitle`ytitle!(x;y)
p:plot `title!,(2a30+?10a;3a270+?10a)



