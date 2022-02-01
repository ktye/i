/ d:(`rect;0 0 300 200
/    `color;255
/    `Rect;20 10 40 30
/    `color;255*256;`Circle;100 100 20
/    `color;0;`linewidth;2;`circle;100 100 50
/    `color;255*256*256;`poly;(0 10 20 30 40;190 100 180 110 190)
/    `color;255;`Poly;(200 250 280;10 100 30)
/    `color;0;`text;(200;120;"a b c")
/    `Text;(180;120;"E F G"))
/ show draw[d;300 200]

b:256*g:256*r:255

loadfont[`monospace;<`"font.ttf"]
m:draw[(`linewidth;5;`color;0;`rect;1 1 398 298);400 300]
m:draw[(`linewidth;5;`color;r;`circle;200 100 80);m]
m:draw[(`linewidth;5;`color;g;`circle;130 200 80);m]
m:draw[(`linewidth;5;`color;b;`circle;270 200 80);m]
m:draw[(`color;255;`font;"20px monospace";`text;(0;20;"ABCDEFG"));m]
show m

/d:(`color;255
/   `Rect;20 10 40 30
/   `linewidth;2
/   `color;255*255;`rect;10 10 100 50
/   `color;255*256*256
/   `Circle;100 100 20
/   `circle;100 100 50
/   `color;255;`line;5 60 10 5
/   `color;255*256*256;`poly;(0 10 20 30 40;190 100 180 110 190)
/   `color;255;`Poly;(100 150 180;10 100 30))
/show draw[d;m]
