/plot `x`y!(xF;yF)     line with titles x and y
/plot `x`y!(xF;yL)     multiple lines, individual axes(xL;yL)
/plot `t!,xZ           polar plot
/plot 3 4 5!10 1 22    bar plot xI!yI

PW:600;PH:400
plotcolors:","\:"#1f77b4,#ff7f0e,#2ca02c,#d62728,#9467bd,#8c564b,#e377c2,#7f7f7f,#bcbd22,#17becf"
plotstyle

plot:{[d];l:!d;v:.d; y:$[`L~@y:v 1;y;,y]; x:$[`L;~@x:v 0;x;(,x)@(#y)#0]
 a:(&/&/x;|/|/x;&/&/y;|/|/y)
 Plot `x`y`t`l`a`f`fw`fh`pw`ph!(x;y;`xy;l;a;"20px monospace";FW;FH;PW;PH)
}


Plot:{[x]; w:x`fw; h:x`fh; W:x`pw; H:x`ph; a:x`a
 dst:(w;W-w;H-h;h);rdst:(w;h;W-2*w;H-2*h)
 xs:(a 0 1)(dst 0 1)'
 ys:(a 2 3)(dst 2 3)'
 d:{(x i)!x 1+i:2*!2\#x}
 line:{[x;y;w;c](`lineWidth;w;`strokeStyle;c;`beginPath;()),(,/(`moveTo,(-1+#x)#`lineTo),','+(xs x;ys y)),(`stroke;())}
 plt::`font!,x`f
 plt,:`fillText!(($a 0;w;H);(s;2\W-w*#s:(x`l)0;PH);(s;W-w*1+#s:$a 1;H))
 plt,:`fillText!((("[",(", "/:$a 2 3),"]");w;h);(s;2\W-w*#s:(x`l)1;h))
 plt,:`rect`strokeStyle`lineWidth`stroke`clip!(rdst;"black";2;();())
 plt,:d@,/line'[x`x;x`y;2;(#x`x)#plotcolors]
 ui::"<canvas width=",($PW)," height=",($PH)," id='plt'/>"
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

