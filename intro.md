a brief introduction about the pecularities of ktye/k

# get/compile
get any version from the [table](https://github.com/ktye/i#build) e.g.
```
$wget https://github.com/ktye/i/releases/download/latest/k.c
$gcc -ok k.c -lm
$k
ktye/k
 █
```

# command line arguments
- `k a.k b.k c.k` executes all files in order, as if they were catenated and drops into the repl
- `k a.k -e '+/x'` loads a.k, evaluates the expression and exits. *e* stands for both: *eval* and *end* (like in awk)
- `k a.k -e` prevents the repl

executing a file parses and executes everything in one go. it is not line oriented as opposed to other apl/k.
If you want output from a line, print or debug:

# print, debug
debug with a backslash: `x+ \y`; it is also dyadic: to include a label ``x+`Y \y:3`` prints `` `Y:3``.
that means you sometimes have to include an @ if you want to force the monadic form.

in the repl output is converted to (clipped) 2d form by default.
if you prepand a space, it uses k syntax in one long line.

both forms are also available as `` `l x`` and `` `k x`` and return chars.
they are defined in [z.k](z.k) which is built in.

for file i/o there are no numeric verbs and there is only 1 form:
- read: `` <`file`` returns chars, e.g. ``x:<`file``
- write: `` `file<"chars only\n" ``
- to stdout: `` `<"..."``


# special forms, adverbs and overloads
the only keyword is `while`. there is block `[x;y;z]` similar to `*(z;y;x)` for use in cond `$[a;b;c]`.

adverbs have verb overloads if the left arg is not a verb: `' / \ ': /: \:` are `bin mod div in join split`.
- split and join allow vectors (e.g. words) on the left and are not restricted to chars
- div: `x\y` is the same as `y%x`. division always remains within their domain, e.g. ints. to get floats uptype must be explicit: `x%2.`
- encode decode need a double slash: `x//y` and `x\\y` (i always implement them late when everything is filled)


# dots, names
a dot is not part of a symbol: `a . b` can be written as `a.b`. to index a list/dict use `` d`a``.

there are no namespaces, only flat global variables and locals.

there are also no undefined variables/errors. when k creates a new symbol, it also creates an associated spot for a global variable which is zero (null verb).
