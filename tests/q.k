ts:{(./!x)!,,/((./.x)*\:(./.y))}
sp:{+/((,/.x)*(,/.y))}
op:{(.x)*\/(.y)}
dim:1025#0
dim[2]:1
dim[4]:2
dim[8]:3
dim[16]:4
dim[32]:5
dim[64]:6
dim[128]:7
dim[256]:8
dim[512]:9
dim[1024]:10
nzpos:{&~0=./.x}
nzcoef:{(./.x)[nzpos x]}
sz:{dim[#,/.x]#2}
basis:{,'/'$(+sz[x]\\nzpos x)}
ket:{"|",/(basis x)[y],">"}
bra:{"<",/(basis x)[y],"|"}
ppk:{ket[x]'(!#basis x)}
ppb:{bra[x]'(!#basis x)}
qpp:{$[(./!x)=`bra;ppb x;ppk x]}
pst:{(nzcoef x),'(+,qpp x)}
mmu:{x(+/*)/:y}
apply:{(!y)!mmu[x;.y]}
mkr:{,/(,/'')(+')x*\:\:y}
inv:{conj'+x}
mpowI:{$[y>0;mmu/y#,x;mmu/(abs y)#,inv x]}
id:{=/:/(!x;!x)}
p2:{x*x}
pN:{t:y#x;$[y=0;1a;*/t]}
mkrN:{mkr/y#,x}
om:{angle[1;360%x]}
cntrl:{[x;y;z;m]t:x#,I;v:x#,I;t[y-1]:P0;v[y-1]:P1;v[z-1]:m;(mkr/t)+mkr/v} / x=total, y=ctrl, z=tgt, uCX = cntrl[2;1;2;X]
                                                                          / numbered from the top
prob:{(p2'abs'nzcoef x),'(,'ppk x)}
nprob:{(p2'abs'nzcoef x)} / returns the numeric probability
pprob:{(,'((`k@'(100*(p2'abs'nzcoef x))),\:"%")),'(,'ppk x)}
SHF:{(,((x-1)#0),1),(id (x-1)),'0}
CLK:{(pN[om x]'(!x))*id x}
zround:{angle'[abs x;_0.5+angle x]}
SQ2:%2
ISQ2:1%%2
vDB:3.%16.
K0:`ket!,1. 0.
K1:`ket!,0. 1.
B0:`bra!,1. 0.
B1:`bra!,0. 1.
P0:(1. 0.;0. 0.)
P1:(0. 0.;0. 1.)
I:(1. 0.;0. 1.)
X:(0. 1.;1. 0.)
Z:(1. 0.;0. -1.)
H:ISQ2*(1 1;1 -1)
Y:(0 1a270;1a90 0)
S:(1 0;0 1a90)
P:{x:180%x;((1;0);(0;angle[1;x]))}
T:(1 0;0 1a45)
U:{x:1p%x;y:180%y;z:180%z;xx:(((cos x%2);-(sin x%2)*1%angle[1;z]);((sin x%2)*1%angle[1;y];(cos x%2)*angle[1;y+z]));zround xx}
Rx:{x:1p%x;xx:(((cos x%2);-(sin x%2)*angle[1;90]);(-(sin x%2)*angle[1;90];(cos x%2)));zround xx}
Ry:{x:1p%x;xx:((cos x%2;-sin x%2);(sin x%2;cos x%2));zround xx}
Rz:{x:180%x;xx:((1%angle[1;x%2];0);(0;angle[1;x%2]));zround xx}
Rn:{P x}
CP:{x:180%x;t:angle[1;x];(1 0 0 0;0 1 0 0;0 0 1 0;(0;0;0;t))}
SX:(0.5*SQ2)*(1a45 1a315;1a315 1a45)
uCX:(1. 0 0 0;0 1. 0 0;0 0 0 1.;0 0 1. 0)
dCX:(1. 0 0 0;0 0 0 1.;0 0 1. 0;0 1. 0 0)
M:ISQ2*(1 1a90 0 0;0 0 1a90 1;0 0 1a90 -1;1 1a270 0 0)
ECR:ISQ2*(0 0 1 1a90;0 0 1a90 1;1 1a270 0 0;1a270 1 0 0)
SWAP:(1. 0 0 0;0 0 1. 0;0 1. 0 0;0 0 0 1.)
PG:(1 0 0 0 0 0 0 0;0 1 0 0 0 0 0 0;0 0 1 0 0 0 0 0;0 0 0 1 0 0 0 0;0 0 0 0 0 0 0 1;0 0 0 0 0 0 1 0;0 0 0 0 1 0 0 0;0 0 0 0 0 1 0 0)
CCX:(1 0 0 0 0 0 0 0;0 1 0 0 0 0 0 0;0 0 1 0 0 0 0 0;0 0 0 1 0 0 0 0;0 0 0 0 1 0 0 0;0 0 0 0 0 1 0 0;0 0 0 0 0 0 0 1;0 0 0 0 0 0 1 0)
QFT:{xx:(1%%x)*(pN[om x]''((!x)*\:!x));zround xx}
/ circuits (example)
QFT3:{HII:mkr[H;I];HII:mkr[HII;I];kk:apply[HII;x];
R221:cntrl[3;2;1;Rn 2];kk:apply[R221;kk];
R431:cntrl[3;3;1;Rn 4];kk:apply[R431;kk];
IHI:mkr[I;H];IHI:mkr[IHI;I];kk:apply[IHI;kk];
R232:cntrl[3;3;2;Rn 2];kk:apply[R232;kk];
IIH:mkr[I;I];IIH:mkr[IIH;H];kk:apply[IIH;kk];
ISW:mkr[I;SWAP];SWI:mkr[SWAP;I];kk:apply[ISW;kk];
kk:apply[SWI;kk];apply[ISW;kk]}
