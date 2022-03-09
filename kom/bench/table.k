/ k out table.k -e

real:{0*x}
user:{1*x}
sys: {0*x}
s:()
.out
s:`$s
t:+s!.'s
`<$"\n"/:`l@`compiled`interpreted`ratio!t,%/t
