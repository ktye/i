a:a:12972^<`w2          /all words accepted as inputs
b:b:2315^<`w1           /true words

a:a 0 1 2 3
b:b 0 1 2

score:{ _+/(x':y),x=y } /x:candidate, y:secret

f:{x; `m<,/score/:\:[a;b] }             /run with f 1


/arthur 2022.03.07:
/ here is a python3 program to solve all 2315 wordles.
/ heuristic is countdistinct+inlist. (i put the 12972x2315 scores in file m)
/ 3 helper functions since python3 is not a list language.
/ (i also wrote all this in pandas and numpy but it was much slower)
/ 
/ lap (list map)
/ ind (x index y)
/ grp (x group y)
/ 
/ from functools import partial
/ n=2315
/ x=bytearray(open('../k/m','rb').read())
/ m=[x[i:i+n]for i in range(0,len(x),n)]
/ 
/ lap=lambda f,x:list(map(f,x))
/ ind=lambda x,y:[x[y]for y in y]
/ 
/ def grp(x,y):
/  d={}
/  for x,y in zip(x,y):
/   if x in d:d[x]+=[y]
/   else:d[x]=[y]
/  return d
/ 
/ def p(x):
/  a=lap(partial(ind,y=x),m)
/  a=lap(lambda x:len(set(x))+(m[0][0]in x),a)
/  a=a.index(max(a))
/  x=[x for x in x if x!=a]
/  return grp(ind(m[a],x),x).values()
/ 
/ g=lambda x:len(x)+(len(x)-1 if 3>len(x) else sum(lap(g,p(x))))
/ x=range(n)
/ print(g(x))
/ 
/ takes about 5 seconds to calculate the 7954 moves (3.435 avg  optimal is
/ 3.42)
/ 
/ the k equivalent (1 second) is:
/ 
/ p::m[a;x]!x_:a:*>(!'m@\x)+(!#m)in x
/ g::n+if 3>n:#x then n-1 else sum g'p x
/
