kst:{$[#x;kstn;kst0]x} /roundtrip, single line, cannot handle null
kst0:{(;_34 34;"!0";"0#0.";"0#0a";"0#`";"()";"0#`!0")@@x}
kstn:{v:{" "/:$x};(($;kstc;v;v;v;ksts;kstl;kstd)@@x)x}
kstd:{kst[!x],"!",kst[.x]}
kstc:{c:_34;c,($x),c}
ksts:{"`","`"/:$x}
kstl:{$[1~#x;",",kst@*x;"(",(";"/:kst'x),")"]}
print:kst

hex:  {,/{"0123456789abcdef"@,/(x%16;16/x)}'0+x}  /hex"alpha"
unhex:{,/{+/_16 1*"0123456789abcdef"?x}'(,/(2\#x;2))#x}
bits: {2/128 64 32 16 8 4 2 1\x}

rle:  {(x i;,/1_-':(i:&1,1_~,/=':x),#x)} /rle"aabbb" /("ab";2 3)

tab:    {t:@x; \(_10)/:$[t~6;tochar'x;t~7;tabdict x; 'kx];}         /matrix output
rpad:   {y#x,(y-#x)#" "}                        /rpad["alpha";10]
tabdict:{k:$!x;pad:rpad[;|/#'k]; ((pad'k),'"|"),' 'k'.x}
tochar: {$[1~@x;x; 'kx]}

zip:{[f;l]f./:&l}       /e.g. zip[f3;(x;y;z)] x y z:lists, does not extend
compose:{[f;x]x{y x}/f} /e.g. f:({-x};{_x};..)

odo:{(*/x)#/:&'x#'|*\-1_1,|x}

std:`kst`hex`zip!(kst;hex;zip)

