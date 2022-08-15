a:1.;b:2.;r:.1
A:B:0.
x:.1p*!30
f:{[a;b;r;x](r*?#x)+(a*sin x)+(b* cos x)}
p:{[]`x`y!(x;(f[a;b;r;x];f[A;B;0;x]))}
