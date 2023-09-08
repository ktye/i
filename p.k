
S:" \t"                    /space
L:"({[;\n"                 /left
R:")}]"                    /right
A:"'/\\"                   /adverb
V:"+-*%&|<>=^!~,#_$?@."    /verb

C:(S,L;R;V;A;".",_,/65 97+\!26;"0123456789"),"\"`-:e\\/\n"
C:@[;;]/[-1+&256;0+C;!#C]@0+

t:1_'(&"\n"=t)^t:1_("
 ;)+'a0q`-:e\/n
;;)+'a0q`-+a'/;
);)+'a0q`++a'';
+;)+'a0q`-+a'';
';)+'a0q`-:a'';
a;)+'bbq`++b'';
0;)+'11q`+:e'';
qrrrrrrtrrrrsrr
`;)+'b0q`++b'';
-;)+'a1q`-+a'';
/ccccccccccccc;
cccccccccccccc;
:;)+'a0q`-+a'';
b;)+'bbq`++b'';
1;)+'11q`++e'';
e;)+'11q`1+1'';
rrrrrrrtrrrrsrr
srrrrrrrrrrrrrr
t;)+'a0q`++a'';")

T:t[;0]?/1_'t

w:{x@&{$[(*x)?S;0;"/"=*x;1=#x;1]}'x:(&10>{T[x;y]}\C x)^x:" ",x}

D:`plus`minus`times`dividedby`and`or`less`more`equal`power`mod`match`join`take`drop`form`find`at`dot
M:`flip`negate`first`reciprocal`where`reverse`upgrade`downgrade`group`shape`enumerate`not`enlist`count`floor`format`unique`atom`value
DMA:$D,M,`each`over`scan`eachprior`eachright`eachleft
dma:(,'V),(V,'":"),(,'A),A,'":"

q:{(*|x)?V}
p:{x:w x
   c:*'x
   d:q'x
   e:c=":"
   b:1_(c?\"[;\n",R,A),1
   i:&d&(-1_1,d|e|c?\L,A)&~b 
   x[i]:x[i],'":"
   x:x(&b&e&~-1_1,c?\L)_!#x
   -1_,/{$[0N~i:dma?x;x;DMA i]," "}'x}

 \p"1_'(&x=y)_y:x,y /comment"
