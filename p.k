
S:" \t"                    /space
L:"({[;\n"                 /left
R:")}]"                    /right
A:"'/\\"                   /adverb
V:"+-*%&|<>=^!~,#_$?@."    /verb

C:(S,L;R;V;A;".",_,/65 97+\:!26;"0123456789"),"\"`_-:e\\/\n"
C:@[;;]/[-1+&256;0+C;!#C]@0+

t:1_'(&"\n"=t)^t:1_("
 ;)+'a0q`_-:e\/n
;;)+'a0q`_-+a'/;
);)+'a0q`_++a'';
+;)+'a0q`_-+a'';
';)+'a0q`_-:a'';
a;)+'bbq`b++b'';
0;)+'11q`_+:e'';
qrrrrrrtrrrrrsrr
`;)+'b0q`b++b'';
_;)+'b0q`_-+b'';
-;)+'a1q`_-+a'';
/cccccccccccccc;
ccccccccccccccc;
:;)+'a0q`_-+a'';
b;)+'bbq`b++b'';
1;)+'11q`_++e'';
e;)+'11q`_1+1'';
rrrrrrrtrrrrrsrr
srrrrrrrrrrrrrrr
t;)+'a0q`_++a'';")

t[10]:"/cccccccccccccc;"

T:t[;0]?/:1_'t

w:{x@&{$[(*x)':S;0b;"/"=*x;1=#x;1b]}'x:(&11>{T[x;y]}\C x)^x:" ",x}

D:`plus`minus`times`dividedby`and`or`less`more`equal`power`mod`match`join`take`drop`form`find`at`dot
M:`flip`negate`first`reciprocal`where`reverse`upgrade`downgrade`group`shape`enumerate`not`enlist`count`floor`format`unique`atom`value
DMA:$D,M,`each`over`scan`eachprior`eachright`eachleft
dma:(,'V),(V,'":"),(,'A),A,'":"

q:{(*|x)':V}
p:{x:w x
   c:*'x
   d:q'x
   e:c=":"
   b:1_(c':"[;\n",R,A),1b
   i:&d&(-1_1b,d|e|c': L,A)&~b 
   x[i]:x[i],'":"
   x:x(&b&e&~-1_1,c':L)_!#x
   -1_,/{$[0N~i:dma?x;x;DMA i]," "}'x}
