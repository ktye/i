
t:{(`color;255;    `Circle;x,4;`color;0;`text;(x;y;"Abcg",$y))}
T:{(`color;255*256;`Circle;x,4;`color;0;`Text;(x;y;"Abcg",$y))}

d:(`font;("Courier New";40)) 
d,: t[10 50;0],  t[100 100;1], t[200 150;2]
d,: t[10 200;7], t[100 250;8], t[200 300;3]
d,: t[10 350;6], t[100 400;5], t[200 450;4]

d,: T[250 150;0], T[250 300;8], T[250 400;4]

show draw[d;300 600]
