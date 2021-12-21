#img
| Command | Description |
| --- | --- |
| `png x` | convert `(i height;I pixels)` to `C` (png bytes)  |
| `svg x` | rasterize C (svg-string) to `(i height;I pixels)` |
| `rgb x` | convert `(r;g;b) 0-255` to I |
| `draw x y`| rasterize draw calls x to image of size y (w;h) |

#draw

Draw calls are a flat list of commands (symbols) with arguments.
Numeric values may be int or float.

Much of the api follows svg.

| Command | Arguments | Description |
| `M`|   2  | create new path and move to x y |
| `L`|   2  | line to x y |
| `l`|   2  | rel line to +x +y |
| `H`|   1  | hor line to x |
| `h`|   1  | rel hor line to +x |
| `V`|   1  | ver line to y |
| `v`|   1  | rel ver line to +y |
| `C`|   6  | cube to x1 y1 x2 y2 x y |
| `c`|   6  | rel cube to +x1 +y1 +x2 +y2 +x +y |
| `S`|   4  | short cube to x2 y2 x y |
| `s`|   4  | rel short cube to +x2 +y2 +x +y |
| `Q`|   4  | quad to x1 y1 x y |
| `q`|   4  | rel quad to +x1 +y1 +x +y |
| `T`|   2  | short quad to x y |
| `t`|   2  | rel short quad to +x +y |
| `A`|   7  | arc to rx ry xrot large sweep x y |
| `a`|   7  | rel arc to rx ry xrot large sweep +x +y |
| `Z`|   0  |  close path |
| `rect`| 4 | add rectangle x y w h |
| `circ`| 3 | add circle cx cy r |
| `elli`| 4 | add ellipse cx cy rx ry |
| `poly`| 2F| add polygon path for points in X and Y |
| `fill`|  0| fill path |
| `stroke`|0| stroke path |
| `flush`| 0| flush pending stroke and fill |
| `save`|  0| save transformation |
| `rest`|  0| restore transformation |
| `co`|    1| set color value (int) |
| `lw`|    1| set line width |
| `tr`|    2| translate x y |
| `ro`|    1| rotate deg |
| `sc`|    2| scale sx sy |

draw is based on `github.com/memononen/nanosvg`

and `nothings.org/stb/stb_truetype.h` for fonts

png uses `nothings.org/stb/stb_image_write.h`
