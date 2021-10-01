/plot `x`y!(xF;yF)     line with titles x and y
/plot `x`y!(xF;yL)     multiple lines, individual axes(xL;yL)
/plot `t!,xZ           polar plot
/plot 3 4 5!10 1 22    bar plot xI!yI

PW:600;PH:400

plot:{[d];l:$!d;v:.d; t:$[2~#l;`xy;`polar];
 y:$[t~`xy; $[`L~@y:v 1;y;,y];              _v]
 x:$[t~`xy; $[`L~@x:v 0;x;(,x)@(#y)#0]; imag v]
 a:(&/&/x;|/|/x;&/&/y;|/|/y)
 
 c:c@(#c:","\:"#1f77b4,#ff7f0e,#2ca02c,#d62728,#9467bd,#8c564b,#e377c2,#7f7f7f,#bcbd22,#17becf")/!#x
 style:$[`t~`xy;"--";".."] / -.| line points bar
 size: $[`t~`xy;2;5]
 lines:{`style`size`color`x`y!(style;size;z;x;y)}'[x;y;c]
 plt::`lines`t`l`a`f`fw`fh`pw`ph!(lines;t;l;a;"20px monospace";FW;FH;$[`t~`xy;PW;PH];PH)
 /Plot`plt
}

x:(0 9)(0 4p)'!10
y:(sin x;cos x)
p:plot `xtitle`ytitle!(x;y)


Plot:{[sym];x:.sym; w:x`fw; h:x`fh; W:x`pw; H:x`ph; a:x`a
 dst:(w;W-w;H-h;h);rdst:(w;h;W-2*w;H-2*h)
 xs:(a 0 1)(dst 0 1)'
 ys:(a 2 3)(dst 2 3)'
 d:{(x i)!x 1+i:2*!2\#x}
 / line:{[x;y;w;c](`lineWidth;w;`strokeStyle;c;`beginPath;()),(,/(`moveTo,(-1+#x)#`lineTo),','+(xs x;ys y)),(`stroke;())}
 line:{[l](`lineWidth;l`size;`strokeStyle;l`color;`beginPath;()),(,/(`moveTo,(-1+#l`x)#`lineTo),','+(xs l`x;ys l`y)),(`stroke;())}
 c:`font!,x`f
 c,:`fillText!(($a 0;w;H);(s;2\W-w*#s:(x`l)0;PH);(s;W-w*1+#s:$a 1;H))
 c,:`fillText!((("[",(", "/:$a 2 3),"]");w;h);(s;2\W-w*#s:(x`l)1;h))
 c,:`rect`strokeStyle`lineWidth`stroke`clip!(rdst;"black";2;();())
 c,:d@,/line'[x`x;x`y;2;(#x`x)#plotcolors]
 (`0)[symc:`($sym),".c";c] /assign to dynamic symbol
 ui::"<canvas width=",($PW)," height=",($PH)," id='",($symc)"'/>"
}

axis:{axes::x;Plot(.`"plt.")}

x:(0 9)(0 4p)'!10
y:(sin x;cos x)
plot `xtitle`ytitle!(x;y)


/ PW:600;PH:400
/ plotcolors:","\:"#1f77b4,#ff7f0e,#2ca02c,#d62728,#9467bd,#8c564b,#e377c2,#7f7f7f,#bcbd22,#17becf"
/ plot:{[d];v:.d
/  plty::$[`L~@y;y;,y];pltx::$[`L~@x;x;(,x)@(#plty)#0];pltl::!d
/  axes::(&/&/pltx;|/|/pltx;&/&/plty;|/|/plty)
/  replot[]
/ }
/ replot:{ /FW FH PW PH axes pltx plty pltl -> plt ui
/  dst:(FW;PW-FW;PH-FH;FH);rdst:(FW;FH;PW-2*FW;PH-2*FH)
/  xs:(axes 0 1)(dst 0 1)'
/  ys:(axes 2 3)(dst 2 3)'
/  d:{(x i)!x 1+i:2*!2\#x}
/  line:{[x;y;w;c](`lineWidth;w;`strokeStyle;c;`beginPath;()),(,/(`moveTo,(-1+#x)#`lineTo),','+(xs x;ys y)),(`stroke;())}
/  plt::`font!,("20px monospace")
/  plt,:`fillText!(($axes 0;FW;PH);(s;2\PW-FW*#s:pltl 0;PH);(s;PW-FW*1+#s:$axes 1;PH))
/  plt,:`fillText!((("[",(", "/:$axes 2 3),"]");FW;FH);(s;2\PW-FW*#s:pltl 1;FH))
/  plt,:`rect`strokeStyle`lineWidth`stroke`clip!(rdst;"black";2;();())
/  plt,:d@,/line'[pltx;plty;2;(#pltx)#plotcolors]
/  ui::"<canvas width=",($PW)," height=",($PH)," id='plt'/>"
/ }
/ x:(0 9)(0 4p)'!10
/ y:(sin x;cos x)
/ plot `xtitle`ytitle!(x;y)

