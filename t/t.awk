{
	x = $0
	sub(/^ /, "", x)
	getline < "t"
	y = $0
	if (match(y, /^\//)) next
	gsub(/.* \//, "")
	if (match($0, /^"[^"]*"$/)) x = "\"" x "\""
	if (x != $0) {
		print "t:" NR ": " y " ? " x
		exit 1
	}
	if (x == "`ok") {
		print "t ok"
		exit 0
	}
}
