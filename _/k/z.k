kst:{$[#x;kstn;kst0]x} /roundtrip, single line
kst0:{(;_34 34;"!0";"0#0.";"0#0a";"0#`";"()";"0#`!0")@@x}
kstn:{v:{" "/:$x};(($;kstc;v;v;v;ksts;kstl;kstd)@@x)x}
kstd:{kst[!x],"!",kst[.x]}
kstc:{c:_34;c,($x),c}
ksts:{"`","`"/:$x}
kstl:{$[1~#x;",",kst@*x;"(",(";"/:kst'x),")"]}
print:kst

hex:{,/{"0123456789abcdef"@,/(x%16;16/x)}'0+x}  /hex"alpha"


std:`kst`hex!(kst;hex);