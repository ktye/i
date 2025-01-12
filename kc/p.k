
S:" \t"                    /space
L:"({[;\n"                 /left
R:")}]"                    /right
A:"'/\\"                   /adverb
V:"+-*%&|<>=^!~,#_$?@."    /verb

C:(S,L;R;V;A;".",_,/65 97+\!26;"0123456789"),"\"`-:e\\/\n"
C:@[;;]/[&256;0+C;!#C]

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

em:{"{",x,"}"}
ar:em@","/$'
`<"int C[256]=",(ar C),";"

t:","/ar'T

`<"\nint T[",($#T),"][",($#*T),"]={",t,"};\n"
