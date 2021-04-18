# pre processor
function join(v, a, b,  i)
{
	r = ""
	for (i=a; i<=b; i++) {
		if (i!=a) {
			r = r " "
		}
		r = r v[i]
	}
	return r
}
function replaceLine(s, p)
{
	gsub(/%/, s, p)
	return p	
}
function replace(s,  i)
{
	n = split(s, a)
	for (i=1; i<=n; i++) {
		v = a[i]
		w = v "*"
		if (w in m) {
			s = join(a, 1+i, n)	
			a[i] = replaceLine(s, m[w])
			n = i
			break
		}
		w = v "'"
		if (w in m) {
			a[i] = ""
			while (i<n) {
				i++
				a[i] = replaceLine(a[i], m[w])
			}
			break
		}
		if (v in m) {
			a[i] = m[v] 
		}
		if (substr(v, 1, 2) == ":.") {
			a[i] = "local.tee $" substr(v, 3)
		} else if (substr(v, 1, 1) == ".") {
			a[i] = "local.get $" substr(v, 2)
		} else if (substr(v, 1, 1) == ":") {
			a[i] = "local.set $" substr(v, 2)
		}
	}
	return join(a, 1, n)
}

BEGIN{ print "(module" }

/^def /{
	sym = $2
	repl = ""
	for (i=3; i <= NF; i++)
		repl = repl " " $(i)
	m[sym] = repl
	next
}

{
	x = $0
	do {
		y = x
		x = replace(x)
	} while	(x != y)
	print y
}

END{ print ")" }
