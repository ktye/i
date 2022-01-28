### draw library

[draw.c](draw.c) is a library that draws 2d vector graphics as a software renderer.  
It uses `stb_image_write.h`, `stb_truetype.h` from [nothings](http://nothings.org/stb)
and [nanosvg](https://github.com/memononen/nanosvg) as a rasterizer.

Besides the header files and the stdlib, no libaries are needed.

Showing an image in a window and reacting to events is done in [another library](../ray).

## image representation

In k a raster `image` is represented a 2-element list with the height and image data `(h;d)`.

```
h int
d ints length w*h
```

Each element of `d` represents the color value of a pixel in row major order starting at top left.  
Colors are 24 bit integers stored in int32, e.g:

```
red:255; green:255*256; blue:255*256*256
```


## functions

|call|arg|return|description|
|---|---|---|---|
|`png`|`image`|C|encode image as png bytes|
|`loadfont`|`(s name;C ttfdata)`|0|decode ttf data|
|`draw`|L,wh|`image`|rasterize vector graphics|

`png` and `loadfont` is available only in the c verions, not for js.
The web version uses system fonts, the c version needs to load fonts from ttf data.

### draw

```
draw[d;image]
draw[d;wh]
```

`draw` takes 2 arguments, vector graphics calls and the background `image`.
If the seconds argument is `I (w;h)` instead of an `image`, a new all-white image is used as a background.  

The first argument `d` is a general list with symbol-argument pairs:  

|symbol|argument|description|
|---|---|---|
|`color`|`i`|for both filling and stroking|
|`font`|`C`|"20px fontname"|
|`linewith`|`w`|for stroking|
|`rect`|`x y w h`|stroke rectangle|
|`Rect`|`x y w h`|fill rectangle|
|`circle`|`x y r`|stroke circle|
|`Circle`|`x y r`|fill circle|
|`poly`|`X Y`|stroke poly line|
|`Poly`|`X Y`|fill polygon|
|`text`|`(x;y;text)`|draw text|
|`Text`|`(x;y;text)`|draw rotated text|

numeric arguments may be float or int.

Before using fonts, use `loadfont` to register a font under a given name.

example:

```
d:(`rect;0 0 300 200
   `color;255
   `Rect;20 10 40 30
   `color;255*256;`Circle;100 100 20
   `color;0;`linewidth;2;`circle;100 100 50
   `color;255*256*256;`poly;(0 10 20 30 40;190 100 180 110 190)
   `color;255;`Poly;(200 250 280;10 100 30)
   `color;0;`text;(200;120;"a b c")
   `Text;(180;120;"E F G"))
   
show draw[d;300 200]  /show is defined in ../ray
```

