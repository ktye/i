t:+`a`b`c`d!(1 2;3 4.;`abc`def;("alpha";"beta"))
c:sqlite `t!,t
s:sqlq[c;"select sum(b) as M from t"]
T:sqlite c

 \"query ",$$[7.~*s`M;`ok;`fail]
 \"roundtrip ",$$[T~`t!,t;`ok;`fail]
