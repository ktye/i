b:256*g:256*r:255
/`"red.png"  <png[(100*100)#r;100]
/`"green.png"<png[(100*100)#g;100]
/`"blue.png" <png[(100*100)#b;100]
/ `"nano.png"<png[*1_svg@<`"img/nano.svg";640]
/`"circle.png"<png[*1_svg@<`"img/circle.svg";100]
`"rect.png"<png[*1_svg@<`"img/rect.svg";100]

p: ((`color;r)
    (`fillRect;10;10;60;40)
    (`color;b)
    (`lineWidth;4)
    (`drawRect;20;20;40;60)
    (`color;g)
    (`fillElli;20;70;30;20))
p,:((`drawCirc;50;50;20)
    (`color;0)
    (`drawPoly;0 30 60 90.;20 80 10 90.))
    
m:draw[p;100 100]
i:m 1;
w:(#i)%*m;
`"draw.png"<png[i;w]
