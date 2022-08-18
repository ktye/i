
/draw commands a list of cmd,arg pairs
l:(`color;      255*256
   `linewidth;  3
   `rect;       10 10 200 100   /stroke(lowercase)
   `Rect;       20 20 50 30     /fill  (uppercase)
   `line;       0 0 200 150
   `color;      255
   `circle;     50 50 80
   `Circle;     50 50 20
   `linewidth;  2
   `poly;       (10 20 30 40;0 10 0 10)
   `Poly;       (10+50*!3;150 100 150)
   `font;       ("monospace";20)
   `text;       (200 150;0;"horizontal") /align:0..8
   `Text;       (200 140;1;"vertical"))
   
/also `clip;x y r(circular) or `clip;x y w h(rectangular)

show[`type`width`height!(`draw;400;300); l]
   