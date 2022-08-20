rem[4]@-5 -4 -3 /3 0 1
rem[-4]@-5 -4 -3 -2 -1 0 1 2 /-1 0 -3 -2 -1 0 -3 -2
`split" a bc  def  " /(,"a";"bc";"def")
`split("ab c";" d e f") /(("ab";,"c");(,"d";,"e";,"f"))
dot[(3^!9);!3] /5 14 23
x:solve[a:3^1 -2 3 8 4 3 9 2 1.;b:3 2 1.];1e-13>|/b+dot[+a;-x] /1
A:(3 5 -8 12.;-2 3 3 0.;7 -8 2 1.);b:dot[(+A);x:1 2 3];0.00001>|/abs x-solve[A;b] /1
A:(3 5 -8 12.;-2 3 3 0.;7 -8 2 1.);b:dot[(+A);x:1 2 3];0.00001>|/abs x-solve[(qr A);b] /1
A:+0a0+(1 -2a90 3;5a90 3 2;2 3 1;4 -1 1);0.0001>|/abs r-solve[A;dot[+A;r:1a30 2a30 3a30]] /1
ej[`k;+`k`a!(`x`x`y`z;1 2 3 4);+`k`b!(`x`y;`p`q)] /+`k`a`b!(`x`x`y;1 2 3;`p`p`q)
x:10000?20;0.1>abs(1;9.5;20*%1%12.)-(20=#?x;avg x;std x) /1 1 1
x:?10000;0.01>abs(0.5;%1%12.)-(avg x;std x) /1 1
x:?-10000;0.02>abs(0 1)-(avg x;std x) /1 1
0.1>_avg@?1000a /1
`pack 1 /0x690400000001000000
`unpack 0x690400000002000000 /2
x~`unpack `pack x:1.2 3.4 /1
x~`unpack `pack x:(1;2 3;`beta) /1
x~`unpack `pack x:`a`b!2 3 /1
x~`unpack `pack x:+`a`b!(1 2;3 4.) /1
