Ki(K:i){i|32\2} Kf(K:f){r:Ktn[1;3];Fr0:x;jr|32\-3} Kz(K:z){r:Ktn[1;4];Zr0:x;jr|32\-4}                   /  +   -   *   %   !   &   |   <   >   =   ~   ,   ^   #   _   $   =   @   .  
Mal(p:i) Mfr(:p) Kfr(:j){t:i.32/j;p:px;c:Ip[-2]-1;!c? if(!t)Ip[-1]/Kfr Jpi; Mfr x; :Ix[-2]:c}           / abs neg sqr sqt til flp rev asc dsc grp not enl srt cnt flr str unq fst val
Ktn(K:ii){r:Mal@y\x?[3;0;2;3;4];Ir[-2]:1;Ir[-1]:y;ri|32\x}                                              / add sub mul div mod min max les mor eql neq cat cut tak drp fmt fnd atx cal
Ufy{(!nx)? x: t:tx0; !t? x: nx/if(t~txi){t:0;break};t?[n#c.xi;n#i.xi;n#f.xi;n#z.xi]}                    /  '   '   '   '   /   /   /   /   \   \   \  \     $[c;cnd;else]     if(c)e;
Exl{r:Ktn[0;xn]tx?[x;xn/Jri:K.Cxi;xn/Jri:K.Ixi;xn/Jri:K.Fxi;xn/Jri:K.Zxi];r}                            / ech pri bot bin ovr fix ecr jon scn fxs ecl spl   ?[i;swc;def]   while(c)e;

Ech{jy<0? x y :ny#x yi} Ecr{jz<0? x[y;z] :nz#x[y;zi]} Ecl{jy<0? x[y;z] :ny#x[yi;z]}                     / i/f: <=> ~ ~neq %div i/\shift |&^ _min $max n/ndo  n#vec cast:i. idx: x0 xi
Bot{jy<0? jz<0? x[y;z] :nz#x[y;zi] jz<0? ny#x[yi;z] :if(nx~ny)Q; nx#x[yi;zi]}                           / n(length) t(type) r(refcount) p(pointer) s(element-size) 
Pri{if(jy<0)Q;x[yn#yi;yn#y i-1]} Ovr{s:@x;(ny-1)/s:x[s;y 1+i];s} Scn{s:@x; ny#i? x[s;yi] :s}
Jon{if(jy<0)Q;r:@y;(ny-1)/{r,:x;r,:y 1+i}; r} Spl{if(!jx<0)Q; if(jy<0)Q; r:(0,!y=x)^y;if(nr)r0:1_r0; r}
Fix{r:y;while(r~y:f y)r:y; r} Fxs{r:y;s:,r;while(r~f y)r:y;s,:r; r}
Sum{jx<0? x :tx?[Ovr[Add;x];xn#Sum Upt[2;x];{s:i.0;nx/s+:Ixi;Ki s};{s:f.0;nx/s+:Fxi;Kf s};{s:z.0;nx/s+:Zxi;Kz s};Q]}

Fst{x0} Enl{Ufy 1#x} Til{tx=2? Wer x :tx~-2? Q :ix#i} Cnt{Ki(jx<0? 1 tx>4? tx-4 :nx)}
Grp{a:Unq x;l2(a;ta? na#!x=ai :na#!x(~~)\ai)} Rev{jx<0? x: x@nx#nx-i-1}
Drp{if(tx~-2)Q;if(jy<0)Q;ix<0? (Ki max[0;ix+ny])#y :!nx>ix? Ktn[ty;0] :n:ny-ix;y n#ix+i}
Tak{if(tx~-2)Q;jy<0? y@!x :ix<0? (Ki nx+max[nx;-ix-1])_y :ix>ny? k:i0;y@x#k<x? +k: 0 :y@!x}
Cut{if(tx=-2){n:ix<0? ny%-ix :ix; x:(Ki ny%n)*!Ki n};if(tx~2)Q;n:nx;x,:K.nx;n#x Kxi+!Ki x[1+i]-xi}
Wer{if(tx~2)Q;m:s:i0;nx/{s+:xi;m:max[m;Ixi]};r:Ktn[2;s];k:i.-1 m<2? nx/xi? r[+k]:i :nx/{a:Ixi;a/Ir[+k]:i}}
Unq{jx<0? x: nx<2? x: x:<x;n:i0;(xn-1)/if(xi~xi+1)+n;r:Ktn[tx;n];k:i.-1;(xn-1)/if(xi~xi+1)Ir[+k]:xi;r} /or x@!(nx#i)=x?x 
Cat{if(jx<0)x:,x;if(jy<0)y:,y;r:Ktn[tx=ty? tx :0;(i.Cnt x)+i.Cnt y];if(!rt){nx/+Kxi;ny/+Kyi};Mcp[pr;px;sx\nx];Mcp[pr+sx\nx;py;sy\nx]}
Fnd{jx<0? In[x;y] :tx=-ty? {tx?[xn/if(Cxi=iy)break;xn/if(Ixi=iy)break;xn/if(Fxi=Fx0)break;xn/if(Zxi=Zx0)break;Q];Ki(i)} :yt>-1? Ecr[Fnd;x;y]}
Ati(K:Ki){tx?[x y;K.Cx y;K.Ix y;K.Fx y;K.Zx y]} Atd{jy<0? x :!ny? Ktn[tx;0] :1=ny? Atx[x;y0] :Dmd[x y0;1_y]}
Atx{tx>4? x.,y :(!ty)? ny#x yi; ty=-2? Ati[x;iy] :ty~2? Q :tx?[ny#x Iyi;ny#Cx Iyi;ny#Ix Iyi;ny#Fx Iyi;ny#Zx Iyi]}
Cal{if(jy<0)Q; !tx>4? Atd[x;y] :a:tx-4;ny>a? Q :ny<a? Prj[x;y;!#y] : ny?[K[ix];K[x0;ix];K[x0;x1;ix];K[x0;x1;x2;ix];K[x0;x1;x2;x3;ix];Q]}
Dmd(4){1=ny? Amd[x;y0;z;a] :u:@x;v:1_x;Amd[x;u;j.0;Dmd[Atx[x;u];y;z;a]]}
Amd(4){if(jx<0)Q;if((2~ty)&2~-ty)Q;x:Use x;if(Kz)a:z.(x y;a);if((jy<0)&ta~-tx)x:Exl x;jy<0? tx?[Kfr jx iy;Jx[iy]:ja; Cx[iy]:ia; Ix[iy]:ia; Fx[iy]:Fa0; Zx[iy]:Za0]
 :if(Ka<0)a:K.ny#a; if(na~ny)Q; if(tx~ta)x:Exl x;a:Exl a; tx?[ny/Kfr jxi;Jx[i]:Ja yi; ny/Cx[yi]:Ca yi; ny/Ix[yi]:Ia yi; ny/Fx[yi]:Fa yi; ny/Fz[yi]:Za yi]}

Abs{jx<0? -1-tx?[x;Ki(ix|-ix);Kf(+Fx0);Kz(+Zx0);Q] :tx?[nx#Abs xi;x;nx#Ixi|-Ixi;nx#+Fxi;nx#+Zxi]} Sqt{jx<0? -1-tx?[Q;%Kf f.ix;%Fx0;Q] :tx?[nx#Sqt xi;Q;nx#%f.ix;nx#%Fxi;Q]}
Neq{tx~ty? 0 :jx<0? Ki(-1-tx?[ix~iy;ix~iy;Fx0~Fy0;Zx0~Zy0;i0]) :nx~ny? 0 :tx?[nx/if(Cxi~Cyi)break;nx/if(Ixi~Iyi)break;nx/if(Fxi~Fyi)break;nx/if(Zxi~Zyi)break;i:i.-1];Ki(i~nx)} Not{0=x}
Flp{tx? x :x(m*!#x)+/!m:|/#'x} Neg{0-x} Sqr{x*x} Con(i:KK){jx<0? jy<0? tx&ty :ty$-tx :jy<0? tx$-ty :nx~ny? Q !tx? 0 :!ty? 0 :tx$ty} Imo(i:ii){!y? x :x:y!x;x<0? x+y :x}
Upt(K:iK){(!t)|(x=ty)|x=-ty? y :y:ty+3?[Kz z.Fy0;Kf f.ix;Ki ix;Q;ny#i.Cxi;ny#f.Ixi;ny#z.Fxi;Q];x~ty? Upt[x;y] :y}

Mod{t:Con[x;y];x:Upt[t;x];y:Upt[t;y];!t? x!'y :if((t~2)&t~-2)Q; t<0? Imo(iy;ix) :jx<0? ny#Imo[yi;ix] :jx<0? nx#Imo[iy;xi] :nx#Imo[yi;xi]}
Min{t:Con[x;y];x:Upt[t;x];y:Upt[t;y];!t? x&'y :t<0? -1-t?[Ki ix_iy;Ki ix_iy;Kf Fx0_Fy0;Q] :jx<0? t-1?[ny#ix_Cyi;ny#ix_Iyi;ny#Fx0_Fyi;Q] :jy<0? y&x :t-1?[ny#i.Cxi_Cyi;ny#Ixi_Iyi;;ny#Fxi_Fyi;Q]}
Max{t:Con[x;y];x:Upt[t;x];y:Upt[t;y];!t? x|'y :t<0? -1-t?[Ki ix$iy;Ki ix$iy;Kf Fx0$Fy0;Q] :jx<0? t-1?[ny#ix$Cyi;ny#ix$Iyi;ny#Fx0$Fyi;Q] :jy<0? y|x :t-1?[ny#i.Cxi$Cyi;ny#Ixi$Iyi;;ny#Fxi$Fyi;Q]}
Les{t:Con[x;y];x:Upt[t;x];y:Upt[t;y];!t? x<'y :t<0? Ki(-1-t?[ix<iy;ix<iy;Fx0<Fy0;Q]) :jx<0? t-1?[ny#ix<Cyi;ny#ix<Iyi;ny#Fx0<Fyi;Q] :jy<0? y>x :t-1?[ny#Cxi<Cyi;ny#Ixi<Iyi;;ny#Fxi<Fyi]}
Mor{t:Con[x;y];x:Upt[t;x];y:Upt[t;y];!t? x>'y :t<0? Ki(-1-t?[ix>iy;ix>iy;Fx0>Fy0;Q]) :jx<0? t-1?[ny#ix>Cyi;ny#ix>Iyi;ny#Fx0>Fyi;Q] :jy<0? y<x :t-1?[ny#Cxi<Cyi;ny#Ixi<Iyi;;ny#Fxi<Fyi]}
Eql{t:Con[x;y];x:Upt[t;x];y:Upt[t;y];!t? x='y :t<0? Ki(-1-t?[ix=iy;ix=iy;Fx0=Fy0;Q]) :jx<0? t-1?[ny#ix=Cyi;ny#ix=Iyi;ny#Fx0=Fyi;Q] :jy<0? y=x :t-1?[ny#Cxi=Cyi;ny#Ixi=Iyi;;ny#Fxi=Fyi]}
Add{t:Con[x;y];x:Upt[t;x];y:Upt[t;y];!t? x+'y :t<0? -1-t?[Ki ix+iy;Ki ix+iy;Kf Fx0+Fy0;Kz Fx0+Fy0] :jx<0? t-1?[ny#ix+Cyi;ny#ix+Iyi;ny#Fx0+Fyi;ny#Zx0+Zy0] :jy<0? y+x  :t-1?[ny#i.Cxi+Cyi;ny#Ixi+Iyi;ny#Fxi+Fyi;ny#Zx0+Zx1]}
Sub{t:Con[x;y];x:Upt[t;x];y:Upt[t;y];!t? x-'y :t<0? -1-t?[Ki ix-iy;Ki ix-iy;Kf Fx0-Fy0;Kz Fx0-Fy0] :jx<0? t-1?[ny#ix-Cyi;ny#ix-Iyi;ny#Fx0-Fyi;ny#Zx0-Zy0] :jy<0? y+-x :t-1?[ny#i.Cxi-Cyi;ny#Ixi-Iyi;ny#Fxi-Fyi;ny#Zx0-Zx1]}
Mul{t:Con[x;y];x:Upt[t;x];y:Upt[t;y];!t? x*'y :t<0? -1-t?[Ki ix*iy;Ki ix*iy;Kf Fx0*Fy0;Kz Fx0*Fy0] :jx<0? t-1?[ny#ix*Cyi;ny#ix*Iyi;ny#Fx0*Fyi;ny#Zx0*Zy0] :jy<0? y*x  :t-1?[ny#i.Cxi*Cyi;ny#Ixi*Iyi;ny#Fxi*Fyi;ny#Zx0*Zx1]}
Div{t:Con[x;y];x:Upt[t;x];y:Upt[t;y];!t? x%'y :t<0? -1-t?[Ki ix%iy;Ki ix%iy;Kf Fx0%Fy0;Kz Fx0%Fy0] :jx<0? t-1?[ny#ix%Cyi;ny#ix%Iyi;ny#Fx0%Fyi;ny#Zx0%Zy0] :jy<0? (1%y)*x :t-1?[ny#i.Cxi%Cyi;ny#Ixi%Iyi;ny#Fxi%Fyi;ny#Zx0%Zx1]}
Bin{y} Enc{y} Dec{y} Str{x} Fmt{y} Srt{x}
