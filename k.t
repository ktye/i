1 /1
1.1 -0w /1.1 -0w
1. 0n 0w -1.1 /1. 0n 0w -1.1
1. 0n 0w -1.1 -0w -0n /1. 0n 0w -1.1 -0w 0n
1. /1.
.1 /0.1
1.1 /1.1
1a /1a
1a 2a20 0a /1a 2a20 0a
1a90 /1a90
1.+2. /3.
1 2+3 /4 5
1+2 3 /3 4
1 2+3 4 /4 6
-(1 2 3) /-1 -2 -3
-(1;2 3) /(-1;-2 -3)
- 0101b /0 -1 0 -1
1a /1a
."1+2" /3
`2 /+
`1+(+) /-
+[2;] /2+
+[;2] /+[;2]
.(1;2;`66) /3
1 /1
1b /1b
101b /101b
`x /`x
``x /``x
`"a" /`a
"" /""
"3" /"3"
1+2 /3
1-2 /-1
1 2 3+4 /5 6 7
1 2 3+4a /5a 6a 7a
1-2 3 /-1 -2
1 2+(3;4 5) /(4;6 7)
3+`a`b!1 2 /`a`b!4 5
(`a`b!7 8)+`a`b!1 2 /`a`b!8 10
1 2,3.5,1. 2.,1a,1a 2a45,"cd",1b /(1;2;3.5;1.;2.;1a;1a;2a45;"c";"d";1b)
"aBc"<"abc" /010b
~1 0 2 /010b
~```a`b /1100b
~" ab" /000b
1a=1 2 /10b
-|1+2 /-3
-`a`b!1 2 /`a`b!-1 -2
0. 1.!0 /0. 1.!0 0
10 20%2. /5. 10.
((4;8 8);1 2 3 4)%(2 3;2) /((2;2 2);0 1 1 2)
101b&110b /100b
101b|110b /111b
#`a`b`c!!3 /3
#+`a`b!(1 2 3;4 5 6) /3
3#!5 /0 1 2
-3#!5 /2 3 4
6#!3 /0 1 2 0N 0N 0N
-6#!3 /0N 0N 0N 0 1 2
"ab"#"abc" /"ab"
2#`a`b`c!1 2 3 /`a`b!1 2
3_!5 /3 4
-3_!5 /0 1
6_!3 /!0
-6_!3 /!0
1_("ab";"cd") /,"cd"
"ab"_"abc" /,"c"
2 5^"alphabeta" /("pha";"beta")
3^!8 /(0 1;2 3;4 5)
"ABC"^"abcCdeAgh" /("abc";"de";"gh")
`a!1 /(,`a)!,1
+(1;2 3;4 5 6.) /((1;2;4.);(1;3;5.);(1;0N;6.))
+("abc";"def") /("ad";"be";"cf")
+`a`b!(1 2;3 4) /+`a`b!(1 2;3 4)
{+/x*y}\:[(1 2 3;4 5 6);7 8 9] /50 122
"b"\:"abc" /(,"a";,"c")
"x"\:"abxdexfg" /("ab";"de";"fg")
"xd"\:"abxdexfg" /("ab";"exfg")
"xy"\:"abcdx" /,"abcdx"
"x"/:("ab";,"c";"ef") /"abxcxef"
"xy"/:("ab";,"c";"ef") /"abxycxyef"
_(97;-2.3 2.3;2a90) /("a";-3 2;0.)
=`b`a`a`c`b /`b`a`c!(0 4;1 2;,3)
?4 3 3 2 4 /4 3 2
?"alpha" /"alph"
1 2 2 1?2 /1
0001001b?1b /3
1 2 2 1?1 2 3 /0 1 4
("abc";"de")?"de" /1
("abc";"de")?("de";"gh") /1 2
((1;"a");(2;"b"))?,(2;"b") /,1
"abc"?(("bc";"a");"cc") /((1 2;0);2 2)
(`a`b!1 2)?2 /`b
3 in 0 1 2 /0b
3 in !5 /1b
"a"in"abc" /1b
"ad"in"abc" /10b
in 000b /0b
in 010b /1b
"ab" find "aaabcabca" /2 5
"aa" find "aaaaaaa" /0 2 4
"ab" find "xab" /,1
abs(1;2a30) /(1;2.)
(1%%2)~imag 1a45 /1b
imag -conj 1 2 imag(3 4.;5 6) /(3. 4.;5. 6.)
(3;4 5)angle 80 /(3a80;4a80 5a80)
angle(3;4 5)angle 80 /(80.;80. 80.)
3=3 /1b
$1b /"1b"
3 4 5=4 /010b
1 2 3 /1 2 3
+/1 2 3 /6
3+/1 2 3 /9
+/0.+!1000 /499500.
,/(1 2 3;,4) /1 2 3 4
|/0#0. /-0w
&/1.1 2.2 /1.1
-2.&/1.1 2.2 /-2.
+\10010b /1 1 1 2 2
4+\5 2 1 /9 11 12
-\5 2 3 1 /5 3 0 -1
3-\5 2 3 1 /-2 -4 -7 -8
*\8 2 1 -3 /8 16 16 -48
|\00b /00b
|\0001101b /0001111b
&\111010101b /111000000b
-'1. 2. /-1. -2.
f:{-x};f'1 2;f'1 2;f'1 2 /-1 -2
v:{$[1<#x;j@*x;x]};j:{v x};j(1;2 3) /1
-'1 2a /1a180 2a180
-'!5 /0 -1 -2 -3 -4
-':3 2 4 0 /0 -1 2 -4
2-':3 2 4 0 /1 -1 2 -4
=':3 2 2 3 /1010b
2~':3 2 2 3 /0010b
1(*-)/:(1;2 3) /0 -1
-/:[1;2 3] /-1 -2
(3+2) /5
(1+2)*3 /9
() /()
*() /""
+/() /""
2#() /("";"")
(1+2;2) /3 2
(;1) /(;1)
(1;;) /(1;;)
1 2@3 /0N
1 2@1b /2
1 2@10b /2 1
"abc"@!5 /"abc  "
3 4 5@(1;(2;0 1)) /(4;(5;3 4))
(2 1)0 /2
1 2 3[1] /2
1 2 3[2 1][0] /3
.[1 2;0] /1
.[(1 2;3 4);1 0] /3
(3^!9).(,1;0 1) /,3 4
(1 2;3 4)[1;0] /3
(1 2;3 4)[1;1 0] /4 3
(1 2;3 4)[1 0;1] /4 2
(1 2;3 4)[;1] /2 4
(`a`b!1 2)`a /1
(`a`b!1 2)`b`a /2 1
@[!5;2 3;9] /0 1 9 9 4
@[1 2;0;3] /3 2
@[1 2;1;+;1] /1 3
@[1 2 3.;1 2;4.] /1. 4. 4.
@[1 2 3a;1 2;4a] /1a 4a 4a
.[2^!6;1 2;9] /(0 1 2;3 4 9)
.[2^!6;1 2;9 9] /(0 1 2;(3;4;9 9))
.[2^!6;1 2;+;1] /(0 1 2;3 4 6)
.[2^!6;(1;1 2);+;1] /(0 1 2;3 5 6)
.[2^!6;1 2;+;1 2] /(0 1 2;(3;4;6 7))
.[3^!9;(1 2;0 2);(0 1;2 3)] /(0 1 2;0 4 1;2 7 3)
.[3^!9;(1 2;0 2);-1] /(0 1 2;-1 4 -1;-1 7 -1)
.[3^!9;(;1);9] /(0 9 2;3 9 5;6 9 8)
.[3^!9;(1;);9] /(0 1 2;9 9 9;6 7 8)
.[3^!9;(0 1;);9] /(9 9 9;9 9 9;6 7 8)
.[+`a!,!5;(`a;1);5] /+(,`a)!,0 5 2 3 4
x:2^!6;x[1;2]:0;x /(0 1 2;3 4 0)
x:!5;x[2]:2. 3.;x /(0;1;2. 3.;3;4)
x:!5;x[2 3]:2. 3.;x /(0;1;2.;3.;4)
d:`a`b!1 2;d[`a]:3;d /`a`b!3 2
@[2 3!4 5;3;+;1] /2 3!4 6
d:`a`b!(1 2!3 4;5);d[`a;2 1]:3 4;d /`a`b!(1 2!4 3;5)
+- /+-
(*-)1 2 /-1
1+ /1+
(1+)2 /3
1+- /1+-
1+/- /1+/-
-1+/- /-1+/-
x:1 /
x:1;x /1
x*x:2 /4
x::1;x /1
x+:x:1;x /2
(1+x;x:3) /4 3
x[1]+:*x:1 2;x /1 3
x:!5;x[2 3]:9;x /0 1 9 9 4
x:!5;x[2 3]+:1 2;x /0 1 3 5 4
f:+;f[3;4] /7
f:1+/;f 3 4 5 /13
(-+)[3;5] /-8
1;2 /2
1;;;2 /2
;;2 /2
{x+y} /{x+y}
.{x+y} /((`y;.;`x;.;+);`x`y;"{x+y}";2)
(.{a:3*x;a:5;x:2})1 /`x`a
(.{a:3*x})3 /1
{x+y}[3;4] /7
{z*x+y}[3;4;5] /35
{{x+y}[x;y]}[1;2] /3
$23 /"23"
$(+':) /"+':"
`i$"31" /31
`F$"1 2 3" /1. 2. 3.
`I$/:("1";"") /(,1;!0)
`$"abc" /`abc
`b`B`i`i`s`S`f$'("0b";"101b";"3";"3b";"`abc";"``a`b";"3.14") /(0b;101b;3;;`abc;``a`b;3.14)
`I`Z$\:"3 4 5" /(3 4 5;3a 4a 5a)
`k@23 /"23"
`l 11 22 33 /,"11 22 33"
`l@`a`be!122 3 /("`a |122";"`be|3")
1>2 /0b
1<2. /1b
1<2a /1b
3 4<`a`b!5 2 /`a`b!10b
`x=`y`x`z /010b
"alphabetagamma"="m" /00000000000110b
+/"alpha"="a" /2
+/"abc"="rbx" /1
1101001b~1101001b /1b
*|!100000 /99999
*`a`b!2 1 /2
!`a`b!1 2 /`a`b
!+`a`b!(1 2;3 4) /`a`b
.`a`b!1 2 /1 2
^2 1 3 4 /1 2 3 4
^1 5 8 2 -5 /-5 1 2 5 8
^10010010b /00000111b
^"abracadabra" /"aaaaabbcdrr"
^("alpha";"abc";"";"ab") /("";"ab";"abc";"alpha")
^-2.1 1.1 0w 0n -0w /0n -0w -2.1 1.1 0w
0n<0w /1b
0n>0w /0b
{x@<x}@-2.1 1.1 0w 0n -0w 0.5 /0n -0w -2.1 0.5 1.1 0w
{x@>x}@-2.1 1.1 0w 0n -0w 0.5 /0w 1.1 0.5 -2.1 -0w 0n
<("alpha";"abc";"";"ab") /2 3 1 0
^-!8 /-7 -6 -5 -4 -3 -2 -1 0
>3 4 /1 0
>4 3 /0 1
<1 8 1 2 5 9 /0 2 3 4 1 5
>1 8 1 2 5 9 /5 1 4 3 0 2
<0 3 2 5 6 5 1 4 9 2 8 0 2 1 3 7 1 5 8 1 8 3 3 1 5 5 4 2 6 1 /0 11 6 13 16 19 23 29 2 9 12 27 1 14 21 22 7 26 3 5 17 24 25 4 28 15 10 18 20 8
>0 3 2 5 6 5 1 4 9 2 8 0 2 1 3 7 1 5 8 1 8 3 3 1 5 5 4 2 6 1 /8 10 18 20 15 4 28 3 5 17 24 25 7 26 1 14 21 22 2 9 12 27 6 13 16 19 23 29 0 11
$[1;2;3+4] /2
$[1>2;2;3+4] /7
while[1-1;2+2] /
n:2;x:3;2*while[n>0;n-:1;x+:1] /10
(3^!9)dot!3 /5 14 23
A:(3 5 -8 12.;-2 3 3 0.;7 -8 2 1.);b:(+A)dot x:1 2 3;0.00001>|/abs x-A solve b /1b
A:+0a0+(1 -2a90 3;5a90 3 2;2 3 1;4 -1 1);0.0001>|/abs r-A solve(+A)dot r:1a30 2a30 3a30 /1b