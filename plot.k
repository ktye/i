/plot `x`y!(xF;yF)     line with titles x and y
/plot `x`y!(xF;yL)     multiple lines, individual axes(xL;yL)
/plot `t!,xZ           polar plot
/plot 3 4 5!10 1 22    bar plot xI!yI


PW:600;PH:400

plot:{plt::`plot x;Plot`plt;s:`c`k`z`d!"plt.",/:"ckzd"
 ui::"<canvas width=",($PW)," height=",($PH)," id='",(s`c),"' data-click='",(s`k),"' data-zoom='",(s`z),"'></canvas><div id='",(s`d),"'></div>"}

`"plot.":{[d];l:$!d;v:.d; t:$[2~#l;`xy;`polar];
 y:$[t~`xy; $[`L~@y:v 1;y;,y];              _v]
 x:$[t~`xy; $[`L~@x:v 0;x;(,x)@(#y)#0]; imag v]
 xt:`tics(&/&/x;|/|/x);yt:`tics(&/&/y;|/|/y)
 a:(xt 0;*-1#xt;yt 0;*-1#yt)
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
 c,:`fillText!(($a 0;w;H);(s;2\W-w*#s:(x`l)0;PH);(s;W-w*1+#s:$a 1;H))
 c,:`fillText!((("[",(", "/:$a 2 3),"]");w;h);(s;2\W-w*#s:(x`l)1;h))
 c,:`rect`strokeStyle`lineWidth`stroke`clip!(rdst;"black";2;();())
 c,:d@,/line'x`L
 xt:`xt \`tics x[`a;0 1]
 yt:`yt \`tics x[`a;2 3]
 back:{[a;dst;x;y]((dst 0 1)(a 0 1)'x;(dst 2 3)(a 2 3)'y)}[a;dst]  /transform canvas to axis
 click:{[b;s;x;y];(`0)[s;"\n"/:("x:";"y:"),'$b[x;y]]}[back;`$($sym),".d"]
 zoom:{[b;s;x]a:b[x 0;x 1],b[x 2;x 3];.[s;(`a;0 2 1 3);@[(^a 0 2),^a 1 3;0 2 1 3]];Plot s}[back;sym]
 (`0)[`$($sym),".c";c]     /assign to dynamic symbol
 (`0)[`$($sym),".k";click]
 (`0)[`$($sym),".d";""]
 (`0)[`$($sym),".z";zoom]
}


`"tics.":{[minmax]nice:{[x;r]f:x%10^ex:_log[10;x];(1 2 5 10.@1+(1 2 5.;1.5 3 7)[r]'f)*10^ex}
 e:nice[-/|`minmax \minmax;0];s:nice[e%4.;1];n:_1.5+e%s;s*(_(*minmax)%s)+!n}

x:(0 9)(0 4p)'!10
y:(sin x;cos x)
p:plot `xtitle`ytitle!(x;y)

 


/n := int(math.Ceil(extent/spacing + 0.5))


/func NiceLimits(min, max float64) (niceMin, niceMax, spacing float64) {
/	maxTicks := 5
/
/	extent := NiceNumber(max-min, false)
/	spacing = NiceNumber(extent/float64(maxTicks-1), true)
/	niceMin = math.Floor(min/spacing) * spacing
/	niceMax = math.Ceil(max/spacing) * spacing
/
/	return niceMin, niceMax, spacing
/}
/func NiceNumber(extent float64, round bool) float64 {
/	exponent := math.Floor(math.Log10(extent))
/	fraction := extent / math.Pow10(int(exponent))
/	var niceFraction float64
/	if round {
/		if fraction < 1.5 {
/			niceFraction = 1
/		} else if fraction < 3 {
/			niceFraction = 2
/		} else if fraction < 7 {
/			niceFraction = 5
/		} else {
/			niceFraction = 10
/		}
/	} else {
/		if fraction <= 1 {
/			niceFraction = 1
/		} else if fraction <= 2 {
/			niceFraction = 2
/		} else if fraction <= 5 {
/			niceFraction = 5
/		} else {
/			niceFraction = 10
/		}
/	}
/	return niceFraction * math.Pow10(int(exponent))
/}