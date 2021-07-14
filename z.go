package main

import . "github.com/ktye/wg/module"

func zk() {
	Data(600, "/32 40 48   56 64   72  80   88   96    104   112\n `k`l`while`in`find`abs`imag`conj`angle`solve`dot\n\n`\".dotmv\":{{+/x*y}\\:[x;y]}\n`\".dot\":{f:.`\".dotmv\";f[x;y]}\n\n`\".slv\":{q:$[`i~@*|x;x;(.`\".qr\")@x];s:.`\".qslv\";$[`L~@y;s/:[q;y];s[q;y]]}\n\n`\".qslv\":{H:x 0;r:x 1;n:x 2;m:x 3;j:0;K:!m\n while[j<n;y[K]-:(+/(conj H[j;K])*y K)*H[j;K];K:1_K;j+:1]\n i:n-1;J:!n;y[i]%:r@i\n while[i;j:i_J;i-:1;y[i]:(y[i]-+/H[j;i]*y@j)%r@i]\n n#y}\n \n`\".qr\":{K:!m:#*x;I:!n:#x;j:0;r:n#0a;turn:$[`Z~@*x;{(-x)angle angle y};{x*1. -1@y>0}]\n while[j<n;I:1_I\n  r[j]:turn[s:0. abs/j_x j;xx:x[j;j]]\n  x[j;j]-:r[j]\n  x[j;K]%:%s*(s+abs xx)\n  x[I;K]-:{+/x*y}/:[(conj x[j;K]);x[I;K]]*\\:x[j;K]\n  K:1_K;j+:1];(x;r;n;m)}\n")
	zn := int32(670) // should end before 2k
	x := mk(Ct, zn)
	Memorycopy(int32(x), 600, zn)
	dx(Val(x))
}
