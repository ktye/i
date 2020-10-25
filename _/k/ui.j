ui:{$[1~@x;key x;#o;mouse x;*(o::(*/x)#0;w::*x;h::x 1)]};o:!w:h:0
key:{$[x~_10;(Y::Y+20;X::-10);X::X+10];o[istr[X;Y;w;x]]::white;o};X:-10;Y:0
mouse:{o[iarc[x;10]]::red;o}

iarc:{[xy;r]x:(w/i:!#o)+r-xy 0;y:(w\i)+r-xy 1;&r>((x*x)+y*y)-r*r}
irgb:{[rgb]+/1 256 65536*rgb};white:irgb 255 255 255;red:irgb 255 0 0
istr:{[x;y;w;c]i:glyphs 0+c; (10*!#s)+(x+w*y)+(10/i)+w*10\i}
\font.k
