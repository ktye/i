k:(`sin!{[R;X]}
   `cos!{[R;X]}
   `exp!{[R;X]}
   `log!{[R;X]}
 `atan2!{[R;X;Y]}
 `hypot!{[R;X;Y]}
  `draw!{[V;x;y;z]}
  `grow!{[r;x]}
`printc!{[V;x;y]}
   `ini!{[r;x]I[0]:134480132;I[4]:67372048;I[12]:1887966018;I[128]:x;p:256;i:8;(i<x)?/(I[4*i]:p;p*:2;i+:1);I[132]:enl mk[1;0];I[136]:enl 0;I[148]:cat[cat[mks 120;mks 121];mks 122];x}
    `bk!{[r;x;y]r:32-*7+y*C x;(r<4)? :4;r}
    `mk!{[r;x;y]t:bk[x;y];b:4*t;m:4*I 128;(0~I b)?/((b>/'m)?(m:grow 1+b%4;I[128]:m;I[b]:1^b_2;m:b;b-:4);b+:4);(128~b)?!0;a:I b;I[b]:I a;j:b-4;(j>/4*t)?/(u:a+1^j_2;I[u]:I j;I[j]:u;j-:4);I[a]:y|x^29;I[a+4]:1;a}
   `mki!{[r;x]r:mk[2;1];I[r+8]:x;r}
   `mkf!{[r;X]r:mk[3;1];F[r+8]:X;r}
   `mkd!{[r;x;y](~5~tp x)?(r:mk[4;1];F[r+8]:`F$0;F[r+16]:`F$1; :add[x;mul[y;r]]);r:nn x;(1~r)?(~1~nn y)?y:enl y;y:lx y;(~r~nn y)?!0;r:l2[x;y];I[r]:2|7^29;r}
   `mkc!{[r;x]r:mk[1;1];C[r+8]:x;r}
   `mks!{[r;x]sc mkc x}
   `mkz!{[r;x;y]r:mkd[x;y];I[r]:2|6^29;r}
    `l2!{[r;x;y]r:mk[6;2];I[r+8]:x;I[r+12]:y;r}
    `l3!{[r;x;y;z]r:mk[6;3];I[r+8]:x;I[r+12]:y;I[r+16]:z;r}
    `nn!{[r;x](x<256)? :1;536870911&I x}
    `tp!{[r;x](x<256)? :0;(I x)_29}
    `fr!{[V;x]xt:tp x;xn:nn x;t:4*bk[xt;xn];I[x]:I t;I[t]:x}
    `dx!{[V;x](x>255)?(xr:I x+4;I[x+4]:xr-1;(1~xr)?(xt:tp x;xn:nn x;xp:8+x;((~xt)+xt>5)?xn/dx I xp+4*i;fr x))}
    `rx!{[V;x]rxn[x;1]}
   `rxn!{[V;x;y](x>255)?(x+:4;I[x]:y+I x)}
    `rl!{[V;x]xn:nn x;xp:8+x;xn/(rx I xp;xp+:4)}
   `rld!{[V;x]rl x;dx x}
    `lx!{[r;x]xt:tp x;xn:nn x;(xt~6)? :x;((xt~7)+xn~1)? :enl x;(0~xt)?!0;r:mk[6;xn];rp:r+8;rxn[x;xn];xn/(I[rp]:atx[x;mki i];rp+:4);dx x;r}
   `til!{[r;x]xt:tp x;xp:8+x;(4~xt)? :zim x;(6~xt)? :ech[x;161];(7~xt)?(r:I xp;rx r;dx x; :r);(~2~xt)?!0;n:I xp;dx x;(n<'0)? :tir[-n];seq[0;n;1]}
   `seq!{[r;x;y;z]r:mk[2;y];rp:8+r;y/(I[rp]:z*i+x;rp+:4);r}
   `tir!{[r;x]r:mk[2;x];rp:4+r+4*x;x/(I[rp]:i;rp-:4);r}
   `upx!{[r;x;y]t:tp x;yt:tp y;(t~yt)? :x;((t~7)+yt~7)?!0;(yt~6)? :lx x;n:nn x;(t<yt)?/(x:up[x;t;n];t+:1);x}
    `up!{[r;x;y;z]r:mk[y+1;z];xp:x+8;rp:r+8;$(1~y;z/(I[rp]:C xp+i;rp+:4);2~y;z/(F[rp]:`F$'I xp;rp+:8;xp+:4);3~y;z/(F[rp]:F xp;F[rp+8]:0.0;xp+:8;rp+:16);!0);dx x;r}
   `atx!{[r;x;y]xt:tp x;xn:nn x;xp:8+x;yt:tp y;yn:nn y;yp:8+y;(~xt)? :cal[x;enl y];(xt~7)? :atd[x;y;yt];(yt>5)? :ecr[x;y;64];(yt~3)?(xt<5)? :phi[x;y];(~yt~2)?!0;r:mk[xt;yn];rp:r+8;w:C xt;yn/(yi:I yp;(xn</yi)?!0;mv[rp;xp+w*yi;w];rp:rp+w;yp+:4);(xt>5)?rl r;(yn~1)?(xt~6)?(rx I r+8;dx r;r:I r+8);dx x;dx y;r}
   `atm!{[r;x;y](~nn y)?(dx x; :y);rx y;f:fst y;t:drop[y;1];(1<nn t)? :atm[atx[x;f];t];t:fst t;nf:nn f;((~f)+~nf~1)?(f?x:atx[x;f]; :ecl[x;t;64]);atx[atx[x;f];t]}
   `atd!{[r;x;y;z]k:I x+8;v:I x+12;(z~5)?(rx k;y:fnd[k;y];z:2);rx v;dx x;atx[v;y]}
   `cal!{[r;x;y]y:lx y;xt:tp x;xn:nn x;xp:8+x;yn:nn y;yp:8+y;xt? :atm[x;y];(yn~1)?(((sadv x)|sadv x-128)? :T[x;`I;`I][fst y];(x<128)?x+:128);(x<128)?((~yn~2)?!0;rld y; :T[x;`I;`I`I][I yp;I yp+4]);(x<256)?((~yn~1)?!0; :T[x;`I;`I][fst y]);(xn~2)?(rld x;a:I xp;$(1~yn; :T[a+128;`I;`I`I][fst y;I xp+4];2~yn;(rld y; :T[a;`I;`I`I`I][I yp;I yp+4;I xp+4]);!0));(xn~3)?(rl x;(1~yn)?y:fst y;r:asi[I x+12;I x+16;y];v:I x+8;dx x; :cal[v;r]);(xn~4)?(a:I x+20;(a>yn)?(a:a-yn;a/y:lcat[y;0]; :prj[x;y;seq[yn;a;1]]); :lcl[x;y]);!0;x}
   `lcl!{[r;x;y]fn:I x+20;(~fn)?(dx y;y:mk[6;0]);(~fn~nn y)?!0;yp:y+8;a:I x+16;ap:a+8;an:nn a;l:mk[2;an];lp:l+8;sp:I 136;an/(d:sp+I ap;I[lp]:I d;v:0;(i<fn)?(v:I yp;rx v;yp+:4);I[d]:v;ap+:4;lp+:4);dx y;r:I x+12;rx r;r:evl r;sp:I 136;lp:l+8;ap:a+8;an/(d:sp+I ap;dx I d;I[d]:I lp;ap+:4;lp+:4);dx l;dx x;r}
   `rev!{[r;x](7~tp x)?(rld x;k:I x+8;v:I x+12; :mkd[rev k;rev v]);n:nn x;(n<2)? :x;atx[x;tir n]}
   `fst!{[r;x]xt:tp x;xn:nn x;(~xn)?(dx x;(xt~0)? :0;(xt~5)? :sc mk[1;0];(xt>5)? :mk[6;0]; :cst[mki xt;mkc 0]);(~xt)? :x;(xt~7)? :fst val x;atx[x;mki 0]}
   `lst!{[r;x](7~tp x)? :lst val x;(~nn x)? :fst x;atx[x;mki (nn x)-1]}
   `cut!{[r;x;y]xt:tp x;xn:nn x;xp:8+x;yt:tp y;yn:nn y;(yt~7)?((xt~2)?((~1~xn)?!0;rld y;k:I y+8;v:I y+12;rx x; :mkd[cut[x;k];cut[x;v]]);rx y; :tkd[exc[til y;x];y]);(~xt~2)?!0;(xn~1)?(r:drop[y;I xp];dx x; :r);r:mk[6;xn];rp:r+8;xn/(a:I xp;b:I xp+4;(i~xn-1)?b:yn;(b<a)?!0;rx y;I[rp]:atx[y;seq[a;b-a;1]];xp+:4;rp+:4);dx x;dx y;r}
   `rsh!{[r;x;y]xt:tp x;xn:nn x;xp:8+x;yt:tp y;(yt~7)?((xt~2)?((~1~xn)?!0;rld y;k:I y+8;v:I y+12;rx x; :mkd[rsh[x;k];rsh[x;v]]); :tkd[x;y]);(~xt~2)?!0;n:prod[xp;xn];r:take[y;n];(xn~1)?((yt~6)?(n~1)?r:enl r;dx x; :r);xn-:1;xe:xp+4*xn;xn/(m:I xe;n:n%m;n:prod[xp;xn-i];r:cut[seq[0;n;m];r];(1~m)?(i>0)?$(6~tp r;r:ech[r;172];r:enl r);xe-:4);(1~I xe)?r:enl r;dx x;r}
  `prod!{[r;x;y]r:1;y/(r:r*I x;x+:4);r}
  `take!{[r;x;y](~nn x)?x:fst x;xn:nn x;o:0;(y<'0)?(o:xn+y;y:-y;(o<'0)? :x);r:seq[o;y;1];(xn<y)?(rp:r+8;y/(I[rp]:i!xn;rp+:4));atx[x;r]}
  `drop!{[r;x;y]xt:tp x;xn:nn x;a:y;(y<'0)?(y:-y;a:0);(y>xn)?(dx x; :mk[xt;0]);x:atx[x;seq[a;xn-y;1]];(xt~6)?(1~xn-y)?x:enl x;x}
   `tkd!{[r;x;y]t:tp x;rld y;k:I y+8;v:I y+12;(~t~5)?!0;rx k;x:fnd[k;x];rx x;v:atx[v;x];(1~nn x)?v:enl v;mkd[atx[k;x];v]}
   `phi!{[r;x;y]n:nn y;r:mk[4;n];rp:r+8;yp:y+8;n/(P:"399d52a246df913f"*\F yp;F[rp]:cos P;F[rp+8]:sin P;rp+:16;yp+:8);dx y;mul[x;r]}
   `use!{[r;x](1~I x+4)? :x;xt:tp x;xn:nn x;xp:8+x;r:mk[xt;xn];rp:r+8;mv[rp;xp;xn*C xt];dx x;r}
    `mv!{[V;x;y;z]z/C[x+i]:C y+i}
   `cat!{[r;x;y]xt:tp x;yt:tp y;(~xt)?(x:enl x;xt:6);(xt~yt)? :ucat[x;y];(xt~6)? :ucat[x;lx y];(yt~6)? :ucat[lx x;y];ucat[lx x;lx y]}
  `ucat!{[r;x;y]xt:tp x;xn:nn x;xp:8+x;yn:nn y;yp:8+y;(xt>5)?(rl x;rl y);(xt~7)?(r:mkd[ucat[I x+8;I y+8];ucat[I x+12;I y+12]];dx x;dx y; :r);r:mk[xt;xn+yn];w:C xt;mv[r+8;xp;w*xn];mv[r+8+w*xn;yp;w*yn];dx x;dx y;r}
  `lcat!{[r;x;y]x:use x;xt:tp x;xn:nn x;xp:8+x;((bk[xt;xn])<bk[xt;xn+1])?(r:mk[xt;xn+1];rld x;mv[r+8;xp;4*xn];x:r;xp:x+8);I[xp+4*xn]:y;I[x]:(xn+1)|6^29;x}
   `enl!{[r;x]r:mk[6;1];I[r+8]:x;r}
   `cnt!{[r;x](7~tp x)?x:til x;r:mki nn x;dx x;r}
   `typ!{[r;x]xt:tp x;r:mk[2;1];I[8+r]:xt;dx x;r}
   `not!{[r;x]t:tp x;(t>5)? :ech[x;126];(~t)?((~x)? :mki 1;dx x; :mki 0);eql[x;mki 0]}
   `wer!{[r;x]xt:tp x;xn:nn x;xp:8+x;(xt~1)? :prs x;(xt~4)? :zan[x;xn;xp];(xt~6)? :flp x;(~xt~2)?!0;rn:0;xn/(rn:rn+I xp;xp+:4);xp:8+x;r:mk[2;rn];rp:r+8;j:0;(j<xn)?/((I xp)/(I[rp]:j;rp+:4);xp+:4;j+:1);dx x;r}
   `mtc!{[r;x;y]r:mk[2;1];I[r+8]:match[x;y];dx x;dx y;r}
 `match!{[r;x;y](x~y)? :1;(~(I x)~I y)? :0;xt:tp x;xn:nn x;xp:8+x;yp:y+8;m:0;$(~xt; :1;1~xt;m:xn;2~xt;m:xn^2;3~xt;m:xn^3;4~xt;m:xn^4;5~xt;m:xn^2;(xn/((~match[I xp;I yp])? :0;xp+:4;yp+:4); :1));m/(~(C xp+i)~C yp+i)? :0;1}
   `fnd!{[r;x;y]xt:tp x;yt:tp y;yn:nn y;yp:8+y;(~xt~yt)?!0;r:mk[2;yn];rp:r+8;w:C yt;yn/(I[rp]:fnx[x;yp];rp+:4;yp:yp+w);dx x;dx y;r}
   `fnx!{[r;x;y]xt:tp x;xn:nn x;xp:8+x;eq:8+xt;w:C xt;xn/((T[eq;`I;`I`I][xp;y])? :i;xp:xp+w);xn}
   `lop!{[r;x;y;z]t:tp y;(~t)? :fxp[x;y;z];(6~t)?(rld y;f:I y+12;y:I y+8; :whl[y;x;f;z]);dx z;0}
   `jon!{[r;x;y]r:lop[x;y;0];r? :r;xt:tp x;xn:nn x;xp:8+x;(~xn)?(r:tp y;dx y;r:mk[r;0];dx x; :r);(1~xn)?(dx y; :fst x);(~xt~6)?(dx y; :x);rl x;r:I xp;rxn[y;xn-2];(xn-1)/(xp+:4;r:cat[cat[r;y];I xp]);dx x;r}
   `spl!{[r;x;y]r:lop[x;y;enl mk[6;0]];r? :r;rx x;yn:nn y;r:fds[x;y];(~nn r)?(dx r; :enl x);r:cut[cat[mki 0;r];x];rn:(nn r)-1;rp:r+8;rn/(rp+:4;I[rp]:drop[I rp;yn]);r}
   `fds!{[r;x;y]xt:tp x;xn:nn x;xp:8+x;yt:tp y;yn:nn y;yp:8+y;((~xt~yt)+xt>5)?!0;(xn<yn)?(dx x;dx y; :mk[2;0]);(~yn)?(dx x;dx y; :drop[seq[0;xn;1];1]);r:mk[2;0];w:C xt;eq:8+xt;j:0;(j<xn)?/(a:0;yn/(k:w*i;a:a+T[eq;`I;`I`I][xp+k;yp+k]);(a~yn)?(r:ucat[r;mki j];j:j+yn-1;xp:xp+w*yn-1);xp:xp+w;j+:1);dx x;dx y;r}
   `exc!{[r;x;y]n:mki nn y;rx x;atx[x;wer eql[n;fnd[y;x]]]}
   `srt!{[r;x]rx x;atx[x;grd x]}
   `gdn!{[r;x]rev grd x}
   `grd!{[r;x]xt:tp x;xn:nn x;xp:8+x;r:seq[0;xn;1];y:seq[0;xn;1];rp:r+8;msrt[y+8;rp;0;xn;xp;xt];dx x;dx y;r}
  `msrt!{[V;x;y;z;x3;x4;x5]((x3-z)>/2)?(c:(x3+z)%2;msrt[y;x;z;c;x4;x5];msrt[y;x;c;x3;x4;x5];mrge[x;y;z;x3;c;x4;x5])}
  `mrge!{[V;x;y;z;x3;x4;x5;x6]k:z;j:x4;w:C x6;i:z;(i<x3)?/(c:k>/x4;(~c)?$(j>/x3;c:0;c:T[x6;`I;`I`I][x5+w*I x+k^2;x5+w*I x+j^2]);$(c;(a:j;j+:1);(a:k;k+:1));I[y+i^2]:I x+a^2;i+:1)}
   `gtc!{[r;x;y](C x)>C y}
   `gti!{[r;x;y](I x)>'I y}
   `gtf!{[r;x;y](F x)>\F y}
   `eqc!{[r;x;y](C x)~C y}
   `eqi!{[r;x;y](I x)~I y}
   `eqf!{[r;x;y]((I x)~I y)*(I x+4)~I y+4}
   `eqz!{[r;x;y](eqf[x;y])*eqf[x+8;y+8]}
   `eqL!{[r;x;y]match[I x;I y]}
   `gtl!{[r;x;y]x:I x;y:I y;xt:tp x;xn:nn x;xp:8+x;yt:tp y;yn:nn y;yp:8+y;(~xt~yt)? :xt>yt;n:xn;(yn<xn)?n:yn;w:C xt;n/(a:xp+i*w;b:yp+i*w;(T[xt;`I;`I`I][a;b])? :1;(T[xt;`I;`I`I][b;a])? :0);xn>yn}
    `sc!{[r;x]k:I 132;n:nn k;x:enl x;r:fnx[k;x+8];$(r<n;dx x;(I[132]:cat[k;x];I[136]:lcat[I 136;0]));r:mki 8+4*r;I[r]:1|5^29;r}
    `cs!{[r;x]r:I (I 132)+I x+8;rx r;dx x;r}
   `eql!{[r;x;y]cmp[x;y;1]}
   `mor!{[r;x;y]cmp[x;y;0]}
   `les!{[r;x;y]cmp[y;x;0]}
   `cmp!{[r;x;y;z]x:upx[x;y];y:upx[y;x];xn:nn x;yn:nn y;(~xn~yn)?((xn~1)?(x:take[x;yn];xn:yn);(yn~1)?(y:take[y;xn];yn:xn));xt:tp x;xp:x+8;yp:y+8;(xt~6)? :ecd[x;y;62-z];f:xt;z?f+:8;w:C xt;r:mk[2;xn];rp:r+8;xn/(I[rp]:T[f;`I;`I`I][xp;yp];xp:xp+w;yp:yp+w;rp+:4);dx x;dx y;r}
   `min!{[r;x;y]mia[x;y;38]}
   `max!{[r;x;y]mia[x;y;124]}
   `mia!{[r;x;y;z]x:upx[x;y];y:upx[y;x];xn:nn x;yn:nn y;(~xn~yn)?((xn~1)?(x:take[x;yn];xn:yn);(yn~1)?(y:take[y;xn];yn:xn));xt:tp x;(xt~6)? :ecd[x;y;z];rx x;rx y;$(z~38;a:les[x;y];a:mor[x;y]);a:wer a;rx a;asi[y;a;atx[x;a]]}
    `nd!{[r;x;y;z;x3]x:upx[x;y];y:upx[y;x];xn:nn x;yn:nn y;(~xn~yn)?((xn~1)?(x:take[x;yn];xn:yn);(yn~1)?(y:take[y;xn];yn:xn));xt:tp x;xp:x+8;yp:y+8;(xt~6)? :ecd[x;y;x3];w:C xt;f:z+xt;r:mk[xt;xn];rp:r+8;xn/(T[f;`V;`I`I`I][xp;yp;rp];xp:xp+w;yp:yp+w;rp:rp+w);dx x;dx y;r}
    `nm!{[r;x;y;z]xt:tp x;xn:nn x;xp:8+x;(xt>5)? :ech[x;z];r:use x;rp:r+8;w:C xt;y:y+xt;xn/(T[y;`V;`I`I][xp;rp];xp:xp+w;rp:rp+w);(xt~4)?(y~19)? :zre r;r}
   `nmf!{[r;x;y]dx x;x:I x+8;y:use cst[mki 3;y];yp:y+8;(nn y)/(F[yp]:T[x;`F;`F][F yp];yp+:8);y}
   `add!{[r;x;y]nd[x;y;143;43]}
   `sub!{[r;x;y]nd[x;y;147;45]}
   `mul!{[r;x;y]nd[x;y;151;42]}
   `diw!{[r;x;y]nd[x;y;155;37]}
   `mod!{[r;x;y]nd[x;y;23;7]}
   `adc!{[V;x;y;z]C[z]:(C x)+C y}
   `adi!{[V;x;y;z]I[z]:(I x)+I y}
   `adf!{[V;x;y;z]F[z]:(F x)+\F y}
   `adz!{[V;x;y;z]adf[x;y;z];adf[x+8;y+8;z+8]}
   `suc!{[V;x;y;z]C[z]:(C x)-C y}
   `sui!{[V;x;y;z]I[z]:(I x)-I y}
   `suf!{[V;x;y;z]F[z]:(F x)-\F y}
   `suz!{[V;x;y;z]suf[x;y;z];suf[x+8;y+8;z+8]}
   `muc!{[V;x;y;z]C[z]:(C x)*C y}
   `mui!{[V;x;y;z]I[z]:(I x)*I y}
   `muf!{[V;x;y;z]F[z]:(F x)*\F y}
   `muz!{[V;x;y;z]F[z]:((F x)*\F y)-\(F y+8)*\F x+8;F[z+8]:((F x)*\F y+8)+\(F x+8)*\F y}
   `dic!{[V;x;y;z]C[z]:(C x)%C y}
   `dii!{[V;x;y;z]I[z]:(I x)%'I y}
   `dif!{[V;x;y;z]F[z]:(F x)%\F y}
   `moi!{[V;x;y;z]x:I x;y:I y;I[z]:(y+`I$x!'y)!'y}
   `diz!{[V;x;y;z]A:F x;B:F x+8;C:F y;D:F y+8;$((+\C)>/\+\D;(R:D%\C;P:C+\R*\D;F[z]:(A+\B*\R)%\P;F[z+8]:(B-\A*\R)%\P);(R:C%\D;P:D+\R*\C;F[z]:(B+\A*\R)%\P;F[z+8]:((B*\R)-\A)%\P))}
   `abx!{[r;x]nm[x;15;171]}
   `neg!{[r;x]nm[x;19;173]}
   `sqr!{[r;x]nm[x;27;165]}
   `abc!{[V;x;y]c:C x;$(is[c;1];C[y]:c-32;C[y]:c)}
   `abi!{[V;x;y]i:I x;$(i<'0;I[y]:-i;I[y]:i)}
   `abf!{[V;x;y]F[y]:+\F x}
   `abz!{[V;x;y]F[y]:hypot[F x;F x+8]}
   `nec!{[V;x;y]c:C x;$(is[c;2];C[y]:c+32;C[y]:c)}
   `nei!{[V;x;y]I[y]:-I x}
   `nef!{[V;x;y]F[y]:-\F x}
   `nez!{[V;x;y]F[y]:-\F x;F[y+8]:-\F x+8}
   `sqc!{[V;x;y]!0}
   `sqi!{[V;x;y]!0}
   `sqf!{[V;x;y]F[y]:%\F x}
   `sqz!{[V;x;y]F[y]:F x;F[y+8]:-\F x+8}
   `lgf!{[r;x]xt:tp x;xn:nn x;xp:8+x;(~xt~3)?!0;x:use x;xp:x+8;xn/(F[xp]:log F xp;xp+:8);x}
   `zre!{[r;x]zri[x;0]}
   `zim!{[r;x]zri[x;8]}
   `zri!{[r;x;y]xn:nn x;xp:8+x;r:mk[3;xn];rp:r+8;xp:xp+y;xn/(F[rp]:F xp;rp+:8;xp+:16);dx x;r}
   `zan!{[r;x;y;z]r:mk[3;y];rp:r+8;y/(F[rp]:ang[F z;F z+8];z+:16;rp+:8);dx x;r}
   `drv!{[r;x;y]r:mk[0;2];I[r+8]:x;I[r+12]:y;r}
   `ecv!{[r;x]drv[40;x]}
   `epv!{[r;x]drv[41;x]}
   `ovv!{[r;x]drv[123;x]}
   `riv!{[r;x]drv[125;x]}
   `scv!{[r;x]drv[91;x]}
   `liv!{[r;x]drv[93;x]}
   `ech!{[r;x;y](tp y)? :bin[y;x];(7~tp x)?(rld x;k:I x+8;v:I x+12; :mkd[k;ech[v;y]]);x:lx x;xn:nn x;xp:8+x;r:mk[6;xn];rp:r+8;rl x;(y<120)?y+:128;xn/(rx y;I[rp]:atx[y;I xp];xp+:4;rp+:4);dx x;dx y;r}
   `ecp!{[r;x;y]rx x;p:fst x;epi[p;x;y]}
   `epi!{[r;x;y;z]n:nn y;(~n)?(dx x;dx z; :y);rxn[y;n];rxn[z;n];r:mk[6;n];rp:r+8;n/(yi:atx[y;mki i];rx yi;I[rp]:cal[z;l2[yi;x]];x:yi;rp+:4);dx yi;dx y;dx z;r}
   `ovr!{[r;x;y]t:tp y;(2~t)? :mod[x;y];ovs[x;y;0;0]}
   `scn!{[r;x;y]t:tp y;t?(t<5)? :diw[x;y];ovs[x;y;enl mk[6;0];0]}
   `ovi!{[r;x;y;z]ovs[y;z;0;x]}
   `sci!{[r;x;y;z]ovs[y;z;enl mk[6;0];x]}
   `scl!{[V;x;y]x?(rx y;xp:x+8;I[xp]:lcat[I xp;y])}
   `ovs!{[r;x;y;z;x3]n:nn x;(~n+x3)?(x:enl x;n:1);rxn[x;n];r:x3;o:1;(~r)?(r:fst x;o:0;n-:1;scl[z;r]);rxn[y;n];n/(r:cal[y;l2[r;atx[x;mki i+1-o]]];scl[z;r]);dx x;dx y;(~z)? :r;dx r;fst z}
   `fxp!{[r;x;y;z]t:x;rx x;1?/(rx x;rx y;r:atx[y;x];((match[r;x])+match[r;t])?(dx x;dx y;dx t;z?r:lcat[fst z;r]; :r);scl[z;x];dx x;x:r);x}
   `ecr!{[r;x;y;z](7~tp y)?(rld y;k:I y+8;v:I y+12; :mkd[k;ecr[x;v;z]]);n:nn y;r:mk[6;n];rp:r+8;rxn[x;n];rxn[y;n];rxn[z;n];n/(I[rp]:cal[z;l2[x;atx[y;mki i]]];rp+:4);dx z;dx x;dx y;r}
   `ecl!{[r;x;y;z](7~tp x)?(rld x;k:I x+8;v:I x+12; :mkd[k;ecl[v;y;z]]);n:nn x;r:mk[6;n];rp:r+8;rxn[x;n];rxn[y;n];rxn[z;n];n/(I[rp]:cal[z;l2[atx[x;mki i];y]];rp+:4);dx z;dx x;dx y;r}
   `whl!{[r;x;y;z;x3]t:tp x;t?(((~t~2)+~1~nn x)?!0;dx x; :nlp[y;z;x3;I x+8]);r:y;scl[x3;r];n:mki 0;1?/(rx x;rx r;t:atx[x;r];(match[t;n])?(dx t;dx n;dx z;dx x;x3?(dx r;r:fst x3); :r);dx t;rx z;r:atx[z;r];scl[x3;r]);x}
   `bin!{[r;x;y]xt:tp x;xn:nn x;xp:8+x;yt:tp y;yn:nn y;yp:8+y;(~xt~yt)?!0;r:mk[2;yn];rp:r+8;w:C xt;yn/(I[rp]:ibin[xp;yp;xn;xt];rp+:4;yp:yp+w);dx x;dx y;r}
  `ibin!{[r;x;y;z;x3]k:0;j:z-1;w:C x3;1?/((k>'j)? :k-1;h:(k+j)_1;$(T[x3;`I;`I`I][x+w*h;y];j:h-1;k:h+1));x}
   `nlp!{[r;x;y;z;x3](x3<0)?!0;r:x;rxn[y;x3];scl[z;x];x3/(r:atx[y;r];scl[z;r]);dx y;z?(dx r;r:fst z);r}
   `ecd!{[r;x;y;z]xn:nn x;yn:nn y;(~xn~yn)?((xn~1)?(x:take[x;yn];xn:yn);(yn~1)?(y:take[y;xn];yn:xn));r:mk[6;xn];rp:r+8;rxn[x;xn];rxn[y;xn];rxn[z;xn];xn/(c:mki i;rx c;I[rp]:cal[z;l2[atx[x;c];atx[y;c]]];rp+:4);dx z;dx x;dx y;r}
   `val!{[r;x]xt:tp x;xn:nn x;$(~xt;((x<256)? :x;rl x;r:mk[6;xn];mv[r+8;x+8;4*xn];(xn~4)?I[r+20]:mki I r+20;dx x);1~xt;(r:prs x;n:(2<nn r)*58~I r+8;r:evl r;n?(dx r;r:0));5~xt;r:lup x;6~xt;r:evl x;7~xt;(r:I x+12;rx r;dx x);!0);r}
   `lup!{[r;x]r:I (I 136)+I x+8;rx r;dx x;r}
   `asn!{[r;x;y](~5~tp x)?!0;p:(I 136)+I x+8;dx I p;I[p]:y;rx y;dx x;y}
   `asd!{[r;x]rld x;v:I x+8;s:I x+12;a:I x+16;u:I x+20;(~v~58)?(rx s;r:lup s;a?(rx a;r:atx[r;a]);u:cal[v;l2[r;u]]);r:u;rx r;a?(rx s;u:asi[lup s;a;u]);dx asn[s;u];r}
   `asi!{[r;x;y;z]xt:tp x;xn:nn x;yt:tp y;yn:nn y;yp:8+y;(xt~7)?(yt<6)?(rld x;k:I x+8;v:I x+12;(yt~5)?(rx k;y:fnd[k;y]); :mkd[k;asi[v;y;z]]);(yt~6)?((xt~7)?(rld x;k:I x+8;v:I x+12;rx y;f:fst y;$(f;(5~tp f)?(rx k;f:fnd[k;f]);f:seq[0;nn k;1]); :mkd[k;asi[v;cat[enl f;drop[y;1]];z]]);((~xt~6)+~yt~6)?!0;r:take[x;xn];(1~xn)?r:enl r;rp:r+8;rx y;a:fst y;y:drop[y;1];(1~nn y)?y:fst y;(~a)?a:seq[0;xn;1];(~2~tp a)?!0;an:nn a;ap:a+8;(an~1)?(dx a;ri:rp+4*I ap;I[ri]:asi[I ri;y;z]; :r);(~yn~2)?!0;(~6~tp z)?z:take[enl z;an];(~an~nn z)?!0;rxn[y;an-1];rl z;zp:z+8;an/(ri:rp+4*I ap;I[ri]:asi[I ri;y;I zp];ap+:4;zp+:4);dx a;dx z; :r);(~yt~2)?!0;zt:tp z;zn:nn z;zp:8+z;(yn>1)?(zn~1)?(~zn~yn)?((~zn~1)?!0;z:take[z;yn];zn:yn;zp:z+8);(xt<6)?((~zt~xt)?!0;r:use x;rp:r+8;w:C xt;yn/(k:I yp;mv[rp+w*k;zp;w];yp+:4;zp:zp+w);dx y;dx z; :r);(xt~6)?(r:take[x;xn];(1~xn)?r:enl r;(1~yn)?(z:enl z;zn:1;zt:6);(~6~zt)?z:lx z;rp:r+8;(~yn~zn)?!0;zp:z+8;rl z;yn/(k:I yp;(~k<xn)?!0;t:rp+4*k;dx I t;I[t]:I zp;yp+:4;zp+:4);dx y;dx z; :r);!0;x}
   `swc!{[r;x]xn:nn x;xp:8+x;i:1;(i<xn)?/(r:I xp+4*i;rx r;r:evl r;((~i!2)|i~xn-1)?(dx x; :r);dx r;i+:1;(~I r+8)?i+:1);dx x;0}
   `ras!{[r;x;y]v:I x+8;(y~3)?(v<256)?((v~58)+v>128)?((v>128)?v-:128;r:I x+12;rxn[r;2];s:fst r;a:drop[r;1];$(nn a;(a:ltr a;an:nn a;(an~1)?a:fst a);(dx a;a:0));u:I x+16;rx u;dx x; :lcat[l3[v;s;a];evl u]);0}
   `ltr!{[r;x]xt:tp x;xn:nn x;xp:8+x;(~xt~6)? :x;rl x;r:mk[6;xn];rp:r+8;xn/(I[rp]:evl I xp;rp+:4;xp+:4);dx x;r}
   `rtl!{[r;x]xt:tp x;xn:nn x;xp:8+x;(~xt~6)? :x;rl x;r:mk[6;xn];rp:r+8+4*xn;xp:xp+4*xn;xn/(rp-:4;xp-:4;I[rp]:evl I xp);dx x;r}
   `evl!{[r;x]xt:tp x;xn:nn x;xp:8+x;(~xt~6)?((xt~5)?(xn~1)? :lup x; :x);(~xn)? :x;(xn~1)? :rtl fst x;v:I xp;(v~36)?(xn>3)? :swc x;r:ras[x;xn];r? :asd r;(v~128)? :lst ltr x;x:rtl x;xn:nn x;xp:x+8;(v~64)?(xn~4)?(rl x;r:asi[I x+12;I x+16;I x+20];dx x; :r);(xn~2)?(rl x;r:atx[I xp;I xp+4];dx x; :r);a:fnl[xp+4;xn-1];a?(rx I x+8; :prj[I x+8;drop[x;1];a]);rx I xp;cal[I xp;drop[x;1]]}
   `prj!{[r;x;y;z]r:mk[0;3];I[r+8]:x;I[r+12]:y;I[r+16]:z;r}
   `fnl!{[r;x;y]r:0;y/((~I x)?((~r)?r:mk[2;0];r:ucat[r;mki i]);x+:4);r}
   `uqg!{[r;x;y]xt:tp x;xn:nn x;xp:8+x;r:mk[xt;0];n:0;w:C xt;xn/(m:fnx[r;xp];(m~n)?(rx x;r:cat[r;atx[x;mki i]];y?y:lcat[y;mk[2;0]];n+:1);y?(yi:y+8+4*m;I[yi]:cat[I yi;mki i]);xp:xp+w);y?r:l2[r;y];dx x;r}
   `unq!{[r;x]uqg[x;0]}
   `grp!{[r;x]uqg[x;mk[6;0]]}
   `flr!{[r;x]xt:tp x;xn:nn x;xp:8+x;(xt>5)? :ech[x;223];(~xt)?(dx x; :mki x);(1~xt)?(dx x; :C xp);(2~xt)?(r:mk[1;xn];rp:r+8;xn/(C[rp+i]:I xp;xp+:4);dx x; :r);(xt~3)?(r:mk[2;xn];rp:r+8;xn/(I[rp]:`I$'F xp;xp+:8;rp+:4);dx x; :r);(xt~4)? :zre x;!0;x}
   `ang!{[R;X;Y]P:57.29577951308232*\atan2[Y;X];(P<\0.0)?P:P+\360.0;P}
   `cst!{[r;x;y]xt:tp x;xn:nn x;yt:tp y;yn:nn y;(xt~5)?(yt~1)?(dx x; :sc y);((~xt~2)+~xn~1)?!0;dx x;x:I x+8;(x<'0)?(x:-x;n:yn%C x;(~yn~n*C x)?!0;r:use y;I[r]:n|x^29; :r);(0~yn)?(dx y;(7~x)? :mkd[mk[5;0];mk[6;0]]; :mk[x;0]);((yt>x)+yt>4)?!0;(8~x)?(n:yn*C yt;r:use y;I[r]:n|1^29; :r);(yt<'x)?/(y:up[y;yt;yn];yt+:1);y}
   `flp!{[r;x]n:nn I x+8;m:nn x;atx[ovr[x;44];ecr[mul[mki n;seq[0;m;1]];seq[0;n;1];43]]}
   `rnd!{[r;x](~1073741825~I x)?!0;dx x;x:I x+8;r:mk[2;x];rp:r+8;x/(I[rp]:rng 0;rp+:4);r}
   `rng!{[r;x]x:I 12;x|:x^13;x|:x_17;x|:x^5;I[12]:x;x}
   `xxx!{[r;x]!0;x}
   `drw!{[r;x;y]xp:8+x;yt:tp y;yn:nn y;w:I xp;(yt~7)?draw[w;I xp+4;y];(yt~2)?draw[w;yn%w;y];dx x;dx y;0}
  `sadv!{[r;x]$(x~39;r:1;x~47;r:1;x~92;r:1;r:0);r}
   `out!{[r;x](~x)? :x;rx x;r:x;(~1~tp r)?r:kst x;printc[r+8;nn r];dx r;x}
   `kst!{[r;x]t:tp x;(~nn x)?(t>1)?(t<6)?(dx x;r:cc[cc[mkc 48;35];48];(t~3)?r:cc[r;46];(t~4)?r:cc[r;97];(t~5)?C[r+10]:96; :r);(7~t)?(rld x;k:I x+8;v:I x+12;k:cc[kst k;33];(0~nn v)?(dx v;v:mki 0); :ucat[k;kst v]);$(6~t;((1~nn x)? :ucat[mkc 44;kst fst x];x:ech[x;235]);x:str x);$(~t;r:x;1~t;r:cc[ucat[mkc 34;x];34];5~t;r:ucat[mkc 96;jon[x;mkc 96]];6~t;r:cc[ucat[mkc 40;jon[x;mkc 59]];41];r:jon[x;mkc 32]);r}
   `str!{[r;x]xt:tp x;xn:nn x;xp:8+x;(xt~1)? :x;(~xt)? :cg[x;xn];((xt>5)+~xn~1)? :ech[x;164];$(2~xt;r:ci I xp;3~xt;r:cf F xp;4~xt;r:cz[F xp;F xp+8];5~xt;(rx x;r:cs x);!0);dx x;r}
    `cc!{[r;x;y]n:nn x;((bk[1;n])<bk[1;n+1])? :ucat[x;mkc y];C[x+8+n]:y;I[x]:1+I x;x}
    `ng!{[r;x;y]y?x:ucat[mkc 45;x];x}
    `cg!{[r;x;y]((~x)+x~128)? :mk[1;0];(x<127)? :mkc x;(x<256)? :cc[mkc x-128;58];(y~2)?(rl x;r:cat[str I x+12;str I x+8]);(y~3)?(rl x;dx I x+16;r:kst I x+12;C[r+8]:91;C[r+7+nn r]:93;r:ucat[str I x+8;r]);(y~4)?(r:I x+8;rx r);dx x;r}
    `ci!{[r;x](~x)? :mkc 48;m:0;(x<'0)?(x:-x;m:1);r:mk[1;0];x?/(c:x!10;r:cc[r;48+c];x%:10);(~nn r)?r:cc[r;48];ng[rev r;m]}
    `cf!{[r;X](~X~\X)? :cc[mkc 48;110];(X~\0.0)? :cc[cc[mkc 48;46];48];m:0;(X<\0.0)?(m:1;X:-\X);(X>\"ffffffffffffef7f")? :ng[cc[mkc 48;119];m];e:0;(X>\1000.0)?/(e+:3;X:X%\1000.0);d:7;(X<\1.0)?(d+:1;(X<\0.1)?(d+:1;(X<\0.01)?(d+:1;(X<\0.001)?(d:7;(X<\1.0)?/(e-:3;X:X*\1000.0)))));n:`I$'X;r:ci n;X:X-\`F$n;d:d-nn r;(d<'1)?d:1;r:cc[r;46];t:0;d/(X:X*\10.0;n:`I$X;r:cc[r;48+n];X:X-\`F$n;t:(1+t)*~n+~i);r:drop[r;-t];e?r:ucat[cc[r;101];ci e];ng[r;m]}
    `cz!{[r;X;Y]A:hypot[X;Y];p:`I$0.5+\ang[X;Y];ucat[cc[cf A;97];ci p]}
   `prs!{[r;x]xt:tp x;xn:nn x;xp:8+x;(~xt~1)?!0;xn:xn+xp;I[8]:xp;(xn>xp)?(47~C xp)?I[8]:com[xp;xn];r:sq xn;$(1~nn r;r:fst r;r:cat[128;r]);dx x;r}
    `sq!{[r;x]r:mk[6;0];q:ex[pt x;x];q?r:lcat[r;q];1?/(v:ws x;p:I 8;(~v)?(v:C p;v:~(v~59)+v~10);v?((p<x)?I[8]:1+p; :r);I[8]:1+p;(~nn r)?r:lcat[r;0];r:lcat[r;ex[pt x;x]]);!0;x}
    `ex!{[r;x;y]((~x)+ws y)? :x;r:C I 8;((is[r;32])+r~10)? :x;r:pt y;(isv r)?(~isv x)? :l3[r;x;ex[pt y;y]];l2[x;ex[r;y]]}
    `pt!{[r;x]r:tok x;(~r)?(p:I 8;(p~x)? :0;l:123~C p;(l+40~C p)?(I[8]:1+p;$(l;(a:0;(91~C 1+p)?(I[8]:2+p;a:sq x;(~nn a)?a:lcat[a;mk[5;0]];a:ovr[a;44]);r:sq x;r:lam[p;I 8;r;a]);(r:sq x;n:nn r;(n~1)?r:fst r;(n>1)?r:enl r))));1?/(p:I 8;b:C p;((p~x)+32~C p-1)? :r;$(is[b;16];r:l2[tok x;r];b~91;(I[8]:1+p;p:sq x;(~nn p)?p:lcat[p;0];r:cat[enl r;p]); :r));!0;r}
   `isv!{[r;x]xt:tp x;xn:nn x;xp:8+x;(~xt)? :1;(xt~6)?(xn~2)?(a:I xp;(a<256)?((is[a;16])|is[a-128;16])? :1);0}
   `lac!{[r;x;y]xt:tp x;xn:nn x;xp:8+x;(xt~6)?((1~xn)?(5~tp I x+8)? :y;xn/(y:lac[I xp;y];xp+:4));(xt~5)?(xn~1)?(p:((I xp)-8)%4;(p>y)?(p<4)? :p);y}
   `loc!{[r;x;y]xt:tp x;xn:nn x;xp:8+x;(~xt~6)? :y;xn/(y:loc[I xp;y];xp+:4);xp:x+8;(xn~3)?(58~I xp)?(r:I xp+4;rx r;s:fst r;n:nn y;(n~fnx[y;s+8])?(rx s;y:cat[y;s]);dx s);y}
   `lam!{[r;x;y;z;x3]$(1~nn z;z:fst z;z:cat[128;z]);(~x3)?(r:I 148;rx r;x3:take[r;lac[z;0]]);v:nn x3;x3:loc[z;x3];n:y-x;t:mk[1;n];mv[t+8;x;n];r:mk[0;4];I[r+8]:t;I[r+12]:z;I[r+16]:x3;I[r+20]:v;r}
    `ws!{[r;x]p:I 8;(47~C p)?(b:C p-1;((b~32)+b~10)?p:com[p;x]);1?/((p~x)?(I[8]:p; :1);b:C p;((b~10)+is[b;64])?(I[8]:p; :0);p+:1;(47~C p)?p:com[p;x]);x}
   `com!{[r;x;y](x<y)?/((10~C x)? :x;x+:1);x}
   `tok!{[r;x](ws x)? :0;p:I 8;b:C p;((is[b;32])+b~10)? :0;5/(r:T[i+136;`I;`I`I`I][b;p;x];r? :r);0}
   `pui!{[r;x;y;z](~is[x;4])? :0;r:0;((is[x;4])*y<z)?/(r*:10;r:r+`I$x-48;y+:1;x:C y);I[8]:y;r}
   `pin!{[r;x;y;z]I[8]:y;u:pui[x;y;z];(~y~I 8)? :mki `I$u ;(x~45)?(y+:1;(y<z)?(x:C y;I[8]:y;u:pui[x;y;z];(~y~I 8)? :mki@-`I$u;I[8]:y-1));0}
   `pfd!{[R;x;y;Z]G:1.0;(Z<\0.0)?G:-\G;1?/(b:C x;$((x<y)*is[b;4];(G:G*\0.1;Z:Z+\G*\`F$b-48);(I[8]:x; :Z));x+:1);Z}
   `pfl!{[r;x;y;z]m:0;(x~45)?(t:C y-1;((t~34)+(t~93)+(t~41)+is[t;7])? :0;m:1);r:pin[x;y;z];y:I 8;((y~z)+~r)? :r;(46~C y)?(r:up[r;2;1];rp:r+8;F[rp]:pfd[y+1;z;F rp]);y:I 8;(y<z)?(101~C y)?(I[8]:y+1;q:pin[C 1+y;1+y;z];(~q)?(I[8]:y; :r);e:I q+8;dx q;N:F r+8;(e<'0)?/(N:N%\10.0;e+:1);(e>0)?/(N:N*\10.0;e-:1);F[r+8]:N);m?(N:F r+8;(N>\0.0)?F[r+8]:-\N);r}
   `num!{[r;x;y;z]r:pfl[x;y;z];(~r)? :r;y:I 8;x:C y;(y<z)?((119~x)+(110~x)+(112~x)+97~x)?((2~tp r)?r:up[r;2;1];y+:1;I[8]:y;(~97~x)?((112~x)?N:3.141592653589793*\F r+8;(110~x)?N:0n;(119~x)?N:0w;F[r+8]:N; :r);r:up[r;3;1];a:pfl[C y;y;z];(~a)?a:mki 0;(2~tp a)?a:up[a;2;1];r:atx[r;a]);r}
   `nms!{[r;x;y;z]r:num[x;y;z];(~r)? :r;1?/(y:I 8;x:C y;((y+2)>z)? :r;(~x~32)? :r;y+:1;q:num[C y;y;z];(~q)?(I[8]:y-1; :r);r:upx[r;q];q:upx[q;r];r:cat[r;q]);r}
   `vrb!{[r;x;y;z](~is[x;24])? :0;(32~C y-1)?((x~92)?(I[8]:1+y; :160);(x~39)?y+:1);r:C y;(z>1+x)?(58~C 1+y)?(y+:1;r+:128);I[8]:1+y;r}
   `chr!{[r;x;y;z](48~x)?(1<z-y)?(120~C 1+y)? :phx[2+y;z];(~x~34)? :0;a:1+y;1?/(y+:1;(y~z)?!0;(34~C y)?(n:y-a;r:mk[1;n];mv[r+8;a;n];I[8]:1+y; :r));r}
   `phx!{[r;x;y]r:mk[1;0];h:1;q:0;1?/(c:C x;((y</x)+~is[c;5])?(I[8]:x; :r);c:c-(48*c<58)+87*c>96;h:~h;h?r:cc[r;c+q^4];q:c;x+:1);x}
   `nam!{[r;x;y;z](~is[x;3])? :0;a:y;1?/(y+:1;((y~z)+~is[C y;7])?(n:y-a;r:mk[1;n];mv[r+8;a;n];I[8]:y; :sc r));x}
   `sym!{[r;x;y;z]((y~z)+~x~96)? :0;y+:1;x:C y;I[8]:y;(y<z)?(r:nam[x;y;z];r? :r;r:chr[x;y;z];r? :sc r);r:mk[5;1];I[r+8]:8;r}
   `sms!{[r;x;y;z]r:sym[x;y;z];(~r)? :r;1?/(y:I 8;q:sym[C y;y;z];(~q)? :enl r;r:cat[r;q]);r}
    `is!{[r;x;y]y&cla x}
   `cla!{[r;x](128<x-32)? :0;C 128+x}
 ,/(`xxx`gtc`gti`gtf`xxx`gti`gtl`mod`xxx`eqc`eqi`eqf`eqz`eqi`eqL`xxx`abc`abi`abf`abz`nec`nei`nef`nez`xxx`moi`xxx`xxx`sqc`sqi`sqf`sqz
    `xxx`mkd`xxx`rsh`cst`diw`min`ecv`ecd`epi`mul`add`cat`sub`cal`ovv`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`asn`xxx`les`eql`mor`fnd
    `atx`xxx`xxx`xxx`xxx`xxx`nmf`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`sci`scv`ecl`exc`cut
    `xxx`xxx`xxx`xxx`drw`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`ovi`max`ecr`mtc`xxx
    `xxx`sin`cos`exp`log`xxx`xxx`xxx`chr`nms`vrb`nam`sms`xxx`xxx`xxx`adc`adi`adf`adz`suc`sui`suf`suz`muc`mui`muf`muz`dic`dii`dif`diz
    `out`til`xxx`cnt`str`sqr`wer`epv`ech`ecp`fst`abx`enl`neg`val`riv`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`lst`xxx`grd`grp`gdn`unq
    `typ`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`scn`liv`spl`srt`flr
    `xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`kst`lgf`xxx`xxx`xxx`prs`xxx`rnd`xxx`xxx`xxx`xxx`xxx`xxx`xxx`xxx`ovr`rev`jon`not`xxx)
  (160;0x204840484848485040604848484848504444444444444444444448604848484848424242424242424242424242424242424242424242424242424240506048484041414141414141414141414141414141414141414141414141414048604800)
   `sin`cos`exp`log`atan2`hypot`draw`grow`printc`ini`mk`mki`mkf`mkd`mkc`mks`mkz`nn`tp`dx`rx`til`atx`cal`fst`drop`mv`cat`ucat`lcat`enl`sc`cs`val`lup`asn`kst`prs)



funcs:{(,/!'x@i)!,/.'x@i:&,/7=@'x}
ftab: {x@*&,/&5=@'x}
data: {x@&,/6=@'k}
exprt:{*|x}
obj:{ 'ox} 
en:{"(",x,")"}
shift:{$["I"~x;">>2";"F"~x;">>3";""]}
rev:{$[6~@x;|x;x]}
return:{$[x;:y;#y;;:y];p:-1+#y;$["return"~6#y p;y;y[p]:"return ",y[p]];y}

c:   {f:funcs x;(_10)/:(,/cdec'[!f;.'.f];cdat x;(ctab x);,/cfun'[!f;.'.f];cmain)}
cdat:{d:data x; "void initdat(){memcpy(C+",($*d),",(char[]){",(","/:$0+d 1),"});}"}
ctab:{"void inittab(){",( ,/{"T[",($x),"]=",($y),";"}'[!#f;f:ftab x] ),"}"}
cfdn:{$[#z 1;""/:(ctyp@*z 2;" ";$y;"(";csig z;")";$[x;cbod z;";"]);""]}
cdec:cfdn[0]
cfun:cfdn[1]
ctyp:{$[`V~x;"void";"Z"<$x;"uint32_t";"float64"]}
csig:{a:1_x 2;n:-1+x 3;","/:(n#a){(ctyp y)," ",$x}'n#a}
cbod:{"{..}",_10}
chead:((_10)/:{"#import<",x,".h>"}'$`stdlib`stdio`stdint`string`math),"char *C;uint32_t *I;double *F;"
cmain:(_34)/:"`"\:"int main(int args, char **argv){C=malloc(1<<16);I=(*uint32_t)C;F=(*double)C;inittab();initdat();ini(16);
}
"

go:   {f:funcs x;o:gohead,(godat x),(gotab x),/gofun'[!f;.'.f];o}
gofun:{$[#y 1;""/:("func ";$x;"(";gosig y;")";goret y;"{";goloc y;gobod y;"}";_10);""]}
gosig:{a:1_x 2;n:-1+x 3;","/:(n#a){($x)," ",gotyp y}'n#a}
goret:{$[`V~r:*x 2;"";"(",($r)," ",(gotyp r),")"]}
gotyp:{$["Z"<$x;"uint32";"float64"]}
goloc:{l:(x 3)_x 2;$[#l;(";"/:{"var ",($x)," ",gotyp x}'l),";";""]}
gotab:{"func initab(){copy(T[0:],[]interface{}{", (","/:$(ftab x)) ,"})};"}
godat:{d:data x; "func initdat(){copy(C[",($*d),":],[]byte{",(","/:$0+d 1),"})};"}
gobod:{b:x 1;";"/:return[`V~*x 2;()gotre/b]}
goasn:{$[x 1;gosto x;(x 2),($[(#go2)~n:go2?,v:x 3;v;go2 1+n]),"=",x 0]}
gocnd:{(5_,/{"else if ",(gobin x 0),"{",(";"/:rev x 1),"}"}'((2\#x),2)#x),"else{",(";"/:rev@*|x),"}"  }
gofoc:{"for i:=uint32(0);i<",x,";i++"}
gofob:{x,"{",(";"/:rev y),"}"}
goatx:{$[|/"CIF"~\:x;golod[x;y];|/go1~\:x;gomon[go1;x;y];"for"~3#x;gofob[x;y];x,en y]}
gocal:{$["?{"~x;gowhl y;|/(go2," ",goc)~\:x; godya[x;y 1;y 0]; "T"~x;goicl y; x,en@ ","/:rev y]}
goicl:{a:","/:(d:`I`F`V!("uint32";"float64";""))[,/`$'x 0];r:d@`$x 1;"T[",(x 2),"].(func(",a,")",r,")"}
gowhl:{"for ",(gobin x[1])," {",(";"/:rev@*x),"}"}
gobin:{$["ib("~3#x;2_x;"n32(ib("~7#x;"!((",7_x;"n32"~3#x;"1==",x;"bi(",x,")"]}
gonot:{"n32(",x,")"}
golod:{x:x,"[",(en y),(shift x),"]";$["C"~*x;"uint32(",x,")";x]}
gosto:{(x 2),"[",(en x 1),(shift x 2),"]=",$["C"~x 2;"byte(",(x 0),")";x 0]}
goiff:{"if ",(gobin x),"{",(";"/:rev y),"}"}
goclz:{"uint32(bits.LeadingZeros32(",x,"))"}
gomon:{s:$[(#x)~n:x?,y; y; x 1+n];$["("~*|s;s,z,")";s,z]}
godya:{[v;x;y]$[v~":"; :x,"=",y;v~"$"; :gocvt[x;y;""];v~"$("; :gocvt[x;y;"'"]; 0]; ib:""; l:$[(#goc)~n:goc?,v;go2;*(goc;ib:"ib")]; s:$[(#l)~n:l?,v;v;l 1+n]; ib:gosnb[s;ib]; ib,"(",gosgn[s;x],gosgn[s;s],gosgn[s;y],")"}
gocvt:{$["I"~x;"uint32";"float64"],"(",gosgn[z;y],")"}
gocon:{$["0n"~x;"math.Float64frombits(18444492273895866368)";"0w"~x;"math.Inf(1)";x~"3.141592";"math.Pi";x~"57.29577";"57.29577951308232";x~"ffffffffffffef7f";"1.7976931348623157e+308";x~"399d52a246df913f";"0.017453292519943295";x]}
gosgn:{o:$["'"~*|x;-1_x;x];$[x~y;o;$[o~x;y;"int32(",y,")"]]}
gosnb:{$[~""~y;y;"'"~*|x;"uint32";""]}
go1:" "\:"+[ math.Abs( -[ - - - %[ math.Sqrt( ! trap( ~ ! _\ math.Floor("
go2:" "\:"% / ! % !( %' ^ << | ^ _ >> $( $' +[ + + + -[ - - - *[ * * * %( /' %[ /"
goc:" "\:"~ == ~[ == < < <( <' <{ <= <{( <=' <[ < >( >' >{ >= >{( >=' >{[ >= >[ > > >"
gotre:{p:-1+#x;m:16/y;n:16\y;$[~y;
 0~m;x,:,gocon@$%y
 1~m;x[p]:$[164~n;gocnd[|x p];128>n;gofoc x p;"*"~n1:_n-128;goclz x p;"~"~n1;gonot x p;en gomon[go1;n1;x p]]
 2~m;x:(-2_x),,$[63~n;goiff;64~n;goatx;46~n;gocal;godya(_n)][x p;x p-1]
 4~m;x,:,$_n
 5~m;x,:,$%y
 6~m;x:((-n)_x),,(-n)#x
 7~m;x:(-4_x),,goasn -4#x
 9~m;x[p]:"return ",x p
15~m;;$y];x}

gohead:(_34)/:"`"\:"package main
import (`unsafe`;`fmt`;`bufio`;`os`;`math`;`math/bits`;`io/ioutil`;`bytes`;`strings`)
var C []byte;var I []uint32;var F []float64;var T [256]interface{}
func init(){initab();C=make([]byte,1<<16);msl();initdat()}
func trap(x uint32){panic(`trap`)}
func n32(x uint32)uint32{if x==0{return 1}else{return 0}}
func bi(x uint32)bool{if 0==x{return false}else{return true}}
func ib(x bool)uint32{if x{return 1}else{return 0}}
func grow(x uint32)uint32{if x>31{panic(`oom`)};c:=make([]byte,1<<x);copy(c,C);C=c;msl();return x}
var sin,cos,exp,atan2,hypot,log=math.Sin,math.Cos,math.Exp,math.Atan2,math.Hypot,math.Log
func draw(x, y, z uint32){fmt.Printf(`draw %x %x %x\n`, x, y, z)}
func printc(x, y uint32){fmt.Printf(`%s\n`, string(C[x:x+y]))}
type sl struct {p uintptr;l,c int}
func msl() {
 cp := *(*sl)(unsafe.Pointer(&C));ip := *(*sl)(unsafe.Pointer(&I));fp := *(*sl)(unsafe.Pointer(&F))
 ip.l, ip.c, ip.p = cp.l/4, cp.c/4, cp.p; fp.l, fp.c, fp.p = ip.l/2, ip.c/2, ip.p
 C = *(*[]byte)(unsafe.Pointer(&cp));I = *(*[]uint32)(unsafe.Pointer(&ip));F = *(*[]float64)(unsafe.Pointer(&fp))}
func mkcs(b []byte)uint32{x:=mk(1,uint32(len(b)));copy(C[x+8:],b);return x}
func fatal(e error){if e!=nil{panic(e)}}
func load(f string){b,e:=ioutil.ReadFile(f);fatal(e);if n:=bytes.Index(b,[]byte(`\n\\`));n!=-1{b=b[:n+1]};dx(out(val(mkcs(b))))}
func main(){ini(16);for _, a := range os.Args[1:]{if strings.HasSuffix(a,`.k`){load(a)};if a==`-e`{os.Exit(0)}}
 s:=bufio.NewScanner(os.Stdin);fmt.Printf(`k.k(go)\n `);
 for s.Scan(){t:=s.Text();if t==`\\\\`{os.Exit(0)};dx(out(val(mkcs([]byte(t)))));os.Stdout.Write([]byte{32})}};
"
 
 