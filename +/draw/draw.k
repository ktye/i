d:(`color;123;`Rect;0 0 100 50)
m:draw[d;100 50]

 \"draw ",$[m~(50;(100*50)#123);"ok";"fail"]


d:(`rect;0 0 300 200
   `color;255
   `Rect;20 10 40 30
   `color;255*256;`Circle;100 100 20
   `color;0;`linewidth;2;`circle;100 100 50
   `color;255*256*256;`poly;(0 10 20 30 40;190 100 180 110 190)
   `color;255;`Poly;(200 250 280;10 100 30)
   `color;0;`text;(200 120;0;"a b c")
   `Text;(180 120;0;"E F G"))
d,:(`clip;200 150 20;`Rect;180 130 100 50)

m:draw[d;300 200]

/ show m

