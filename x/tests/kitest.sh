k ../ki.k -e 'x:kixx@<`"blank.i"' > out1
k blank.k   -e 'x:`<"\n"/:(`lxy 40 100)@+`T`P`I`S!(T;P;I;S)' > out2

vimdiff out1 out2
