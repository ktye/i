#patch tail calls
#substitute call and call indirect with return call
#if they are the last instruction in a function.
# (next line starts with ^\(
#or if the next line is return
#dont subsitute if the function does not have a result
#
# e.g. use:
# wg -tags 'small,simd4' -nomain . |awk -f tco.awk| wat2wasm --enable-tail-call - -o /c/k/ktye.github.io/k.wasm
#

BEGIN{c=0;r=0}
#{ print ">>>" $0 }
/^\(func/{r=0}
/^\(func.*result/{r=1}
/^ *return/{
	if(c&&r){print " return_" x; x=""}
	else{
		print x
		x=$0
	}
	c=0
	next
}
/^\(/{
	if(c&&r) print " return_" x
	else     print x
	c=0
	x=$0
	next
}
/^ *call/{
	c=1
	print x
	sub(/^ */,"")
	x=$0
	next
}
{
	print x
	c=0
	x=$0
}
END{print x}
