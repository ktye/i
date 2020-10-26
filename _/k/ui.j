ui:{$[1~@x;key x;#o;mouse x;*(o::(*/x)#0;C::20\w::*x;R::20\h::x 1)]};o:!w:h:R:C:0
key:{term x} // $[x~_10;(Y::Y+20;X::-10);X::X+10];o[istr[X;Y;w;x]]::white;o};X:-10;Y:0
mouse:{o[iarc[x;10]]::red;o}


X:0;
term:{X+:10;o::@[@[o;block X;0];block X+10;white];o::$[10=x;scroll o;@[o;istr[X;(h-20);w;x];white]]}
scroll:{x:@[x;block X+10;0];X::0; @[(n_x),(n:20*w)#0;block[10];white]}

block:{[x]; (x+(#o)-20*w)+(10/i)+w*10\i:!10*20} /((10*x)+(#o)-20*h)
iarc:{[xy;r]x:(w/i:!#o)+r-xy 0;y:(w\i)+r-xy 1;&r>((x*x)+y*y)-r*r}
irgb:{[rgb]+/1 256 65536*rgb};white:irgb 255 255 255;red:irgb 255 0 0
istr:{[x;y;w;c]i:glyphs 0+c; (10*!#s)+(x+w*y)+(10/i)+w*10\i}
\font.k
