# kitest asn
# kitest cal..
p="$1"
k ../ki.k -e 'x:kixx@<`"'$p'.i"' > out1
k ${p}.k   -e 'x:`<"\n"/:(`lxy 40 100)@+`T`P`I`S!(T;P;I;S)' > out2

diff out1 out2
