st:()
ps:{ev'`t x;`0}
pu:{st,:x}
po:{r: *|st;st::-1_st;r}
xy:{r:-2#st;st::-2_st;r}
sy:{$[r:bi x;r[];x]}
ev:{t:@x;$[t~`s;sy x;pu x]}
dy:{pu@x/xy[]}
ad:{dy(+)}
ro:{(x#y),x_y}
bi:{
 stack:{ \st}
 add:{dy(+)}
 sub:{dy(-)}
 eq: {dy(~)}
 ne: {dy(~~)}
 dup:{pu[*|st]}
 pop:{po[]}
 copy:{[]n:-po[];st,:n#st}
 roll:{[]i:po[];n:-po[];st::(n_st),ro[-i;n#st]}
 exch:{pu'xy[]} }

