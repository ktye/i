a:1.;b:2.;r:.5;
x:.4*!50
R:?-#x

f:{[a;b;r;x](r*R)+(a*sin x)+(b* cos x)}
fit:{[]ab::solve[(sin x;cos x);f[a;b;r;x]]}
fit[]

p:{[]`x`y!(x;(f[a;b;r;x];f[ab 0;ab 1;0;x]))}
