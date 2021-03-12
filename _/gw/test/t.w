neg:I:I{}add:I:II{}
f2:I:II{x+y}		/1 2	3
f3:I:II{x-y}		/1 2	4294967295
f4:I:I{a:x;1+a}  	/3	4
f5:I:II{x f2 y}		/1 2	3
f6:I:II{x add y}	/1 2	3
f7:I:I{neg 0-x}		/3	3
f8:I:II{(I.2*x)(x;y)}	/5 2	7
f9:I:II{r:0;x/(r+:y);r}	/3 5	15
f10:I:I{I x}		/0	1234
f11:I:I{I x}		/4	5678
f12:F:I{F x}		/8	1.2e-125
f13:I:III{x? :y;z}      /1 2 3	2
f14:I:III{x? :y;z}      /0 2 3	3
f15:I:III{$[0; :1;x; :y; :z];7}	/0 1 2	2
f16:I:III{$[0; :1;x; :y; :z];7}	/1 2 3	2
f17:I:I{~x}		/1	0
f18:I:II{x~y}		/2 2	1
f19:I:II{x<y}		/2 3	1
f20:I:II{x<'y-3}	/1 2	0
f21:I:II{x>'y-3}	/1 2	1
f22:I:II{x>y}		/1 2	0

0!{d20400002e160000d6609c29013f0026}/1234 5678 1.2e-125
10:{f2;f3;f4}