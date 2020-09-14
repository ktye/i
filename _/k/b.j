\listbox.k
\m.k
tmsolve:{[A;BT](qr A)qrsolve/BT}

ucal:{[s;t;r]n:#r;i:(u:?r)?r;X:tmsolve[t,'0a+i=(n#,!#u);&s];K:(#t 0)#'X}
/ s0:1a0 2a30
/ s:(s0;s0+2a0 0a0;s0+0a0 2a0)
/ t:(0a0 0a0;1a0 0a0;0a0 1a0)
/ r:`a`a`a

pcal:{[s;u]i:!2\#s;&tmsolve[(s i)-s i+1;&u i]}
pcnl:{[s;u;f;n]i:!2\#s;pcal[s g:(n&#i)#>f s i;u g]}

uidx:{[u;a]solve[1@0.+a;u]}
/ u:4.8a48 4.9a47 1.6a277 1.7a279
/ a:(0 0;0 0;0 180;0 180)



Tags:`List`idx
List:`u`a!(4.8a48 4.9a47 1.6a277 1.7a279;(0 0;0 0;0 180;0 180))
idx:{l:walk path;uidx[l`u;l`a]}


Tags:`List`unb
List:`s`t`r`c!((3a20 4a50;12.8a5 4.1a51;3.1a19 12.9a14;3a300 2a10);(0a0 0a0;10a0 0a0;0a0 10a0;0a0 0a0);`a`a`a`a;1 1 1 0)
unb:{w:&List`c;s:List[`s;w];t:List[`t;w];r:List[`r;w];K:qr[ucal[s;t;r]];l:walk path;K qrsolve/l`s}

tags: Tags
