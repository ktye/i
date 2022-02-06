### draw library

[draw.c](draw.c) is a library that draws 2d vector graphics using [cairo](https://www.cairographics.org).

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
|`png[image]`|`image`|C|encode image as png bytes|
|`draw[d;wh]`|L, I|`image`|rasterize vector graphics|
|`draw[d;image]`|L, L|`image`|draw over background image|

`png` is available only in the c verions, not for js.

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
|`font`|`C i`|name size|
|`linewith`|`w`|for stroking|
|`rect`|`x y w h`|stroke rectangle|
|`Rect`|`x y w h`|fill rectangle|
|`circle`|`x y r`|stroke circle|
|`Circle`|`x y r`|fill circle|
|`poly`|`X Y`|stroke poly line|
|`Poly`|`X Y`|fill polygon|
|`text`|`(x;y;text)`|draw text|
|`Text`|`(x;y;text)`|draw rotated text|

Most numeric arguments may be float or int.


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

