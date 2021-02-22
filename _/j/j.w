j:I:I{(~I 0)? :ii x;s:I 8;p:I s0;t:I s1;(1~I s4)?((41~x)?s4::0;rz);(40~x)?(s4::1;rz);                                  /p
 n:I s2;(47<x)?(58>x)?(x-:48;(~x|n)?(t:pc 1;rz);n*:10;n+:x;s2::n;rz);n?(t:pc 1|n<<1;s2::0);                            /a
 y:I s3;(96<x)?(123>x)?(y*:32;y+:x-96;s3::y;rz);y?(t:pc 2|y<<2;s3::0);                                                 /r
 (33>x)?((10~x)?((~t~p)?!;ex p;s0::mk 0;s1::I s+8; :1);rz);                                                            /s x cute
 (91~x)?(t:pc mk 0;s1::I(pl t);rz);(93~x)?(t:pa t;(~t)?!;s1::t;rz);t:pc 4|(x-33)<<3;0}                                 /e
 
ex.V:I{(~xl)?!;r:0;p:x+8;l:x;nx?l:pl x;
 (p<=l)?/(c:I p;t:p~l;$[t*2~c&3;(pu lu c;tc);t*93~c;(tc);t*127~c;(e:po;h:po;$[ip;(dx e;pu h;tc);(dx h;pu e;tc)]);
 2~c&3;ex lu c;~4~c&7;pu rx c;740~c;(h:po;r:sw r;pu h;r:sw r);724~c;(r:sw r;h:po;r:sw r;pu h);(V.c>>3)(x)];p+:4);
 dx x;dx r}  tc:{dx x;x:po;p:x+4;l:pl x}                                                                               /tail call

ii.I:I{0::x;p:128;i:7;(i<x)?/((4*i)::p;p*:2;i+:1);4::mk sz;(4+I 4)::0;s:mk 5;(8+s)::mk 0;(12+s)::I s+8;8::s;12::mk 0;x}sz:{126} 
bk.I:I{r:32-*7+4*x;(r<4)? :4;r}rx.I:I{xl?(x::1+I x);x}nn.I:I{I 4+x}nx:{(nn x)}xl:{(~x&7)}rz:{ :0}                      /bucket type, ref
dx.V:I{x?xl?(r:(I x)-1;x::r;(~r)?(n:I x+4;p:x+8;n/(dx I p;p+:4);fr x))}fr.V:I{p:4*bk I 4+x;x::I p;p::x}                /unref, free
mk.I:I{t:bk x;i:4*t;m:4*I 0;(~I i)?/((i>='m)?!;i+:4);a:I i;i::I a;k:i-4;(k>=4*t)?/(u:a+1<<k>>2;u::I k;k::u;k-:4);a::1;(a+4)::x;a}
sw.I:I{(~x)?x:mk 0;s:I 4;4::x;s}s0:{(s+8)}s1:{(s+12)}s2:{(s+16)}s3:{(s+20)}s4:{(s+24)}

lc.I:II{n:nx;(1~I x)?((bk 1+n)~bk n)?((8+x+4*n)::y;(4+x)::1+n; :x);r:mk 1+n;(cp(x;r;n))::y;dx x;r}cp.I:III{x+:8;y+:8;z/(y::rx I x;x+:4;y+:4);y}
pc.I:I{s:I 8;p:I 8+s;t:I 12+s;q:pa t;r:lc(t;x);(12+s)::r;(t~p)?((8+s)::r; :r);(pl q)::r;r}
pa.I:I{p:I 8+I 8;1?/((~nn p)? :p;l:I(pl p);((l~x)+(~l)+p~x)? :p;p:l);p}
pl.I:I{4+x+4*nx}fi.I:I{(~nx)?!;r:rx I x+8;dx x;r}us.I:I{(1~I x)? :x;n:nx;r:mk n;n:cp(x;r;n);dx x;r}
ip:{(ipop x)}ipop.I:I{x:po;(~x&1)?!;I?x%'2}lp:{(lpop x)}lpop.I:I{x:po;(7&x)?!;x}
pi.V:I{pu 1+2*x}ln.I:I{p:(I 12)fn x;(~p)?!;rx I p}pu.V:I{s:I 4;n:nn s;(n~sz)?!;(4+s)::1+n;(pl s)::x}px:{pu x}
fn.I:II{n:nx>>1;p:x+8;n/((y~I p)? :4+p;p+:8);0}lu.I:I{p:(I 12)fn x;(~p)?!;rx I p}
po:{(pop x)}pop.I:I{x:I 4;n:nx;(~n)?!;r:I x+4*1+n;(4+x)::n-1;r}
   
add.V:I{pi ip+ip}sub.V:I{pi(-ip)+ip}mul.V:I{pi ip*ip}div.V:I{swp x;pi ip%ip}                                            /+-*/     
mod.V:I{swp x;pi I?ip\'ip}lti.V:I{pi ip>'ip}eql.V:I{pi ip~ip}gti.V:I{pi ip<'ip}                                         /%<=> 
dup.V:I{x:po;px;px}drp.V:I{dx po}swp.V:I{x:po;y:po;px;pu y}rol.V:I{x:po;y:po;z:po;px;pu z;pu y}                          /"_/~|                                  
cnt.V:I{x:po;r:-1;xl?r:1+2*nx;px;pu r}atx.V:I{i:ip;l:lp;((i<0)+i>='nn l)?!;pu rx I 8+l+4*i;dx l}                        /#@
amd.V:I{v:po;i:ip;a:us lp;n:nn a;$[i~n;a:a lc v;(i<'0)+i>'n;!;(ap:8+a+4*i;x:rx I ap;ap::v)];pu a}                       /$
cat.V:I{y:po;x:po;(7&x)?x:(mk 0)lc x;$[7&y;x:x lc y;(yp:y+8;(nn y)/(x:x lc rx I yp;yp+:4);dx y)];px}                    /,
asn.V:I{y:fi lp;(~2~y&3)?!;v:po;(v&7)?(v:(mk 0)lc v);s:I 12;p:s fn y;(~p)?(s:s lc y;s:s lc 1;p:pl s);dx I p;p::v;12::s} /: 
ife.V:I{e:po;t:po;$[~ip;(dx t;ex e);(dx e;ex t)]}exe.V:I{ex lp}stk.V:I{!}                                               /?.!
0:{stk;dup;cnt;amd;mod;xxx;xxx;xxx;xxx;mul;add;cat;sub;exe;div}25:{asn;xxx;lti;eql;gti;ife;atx}62:{drp}91:{rol;xxx;swp}
