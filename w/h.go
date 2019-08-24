package main

const h = `<html><head><script>

</script>
<meta charset="utf-8">
<style>
	body,textarea { font-family: monospace; margin:0pt; }
	textarea {
		background-color: black;
		color: white;
		border: none;
		resize: none;
	}
	.col { float: left; width:50%; height:100%; }
	.row:after { content: ""; display: table; clear: both }
</style>
</head><body>
<div id="dropbox">
<div class="row"><textarea id="term" class="col"></textarea><canvas id="draw" class="col"></canvas></div></div>
<script>
function e(s) {
	var req = new XMLHttpRequest()
	req.onreadystatechange = function() { if (this.readyState == (this.DONE || 4)) { O(req.response+" ");term.scrollTo(0, term.scrollHeight)  } }
	req.open("POST", "")
	req.send(s)
}

var term = document.getElementById("term")
var hold = false

term.value = window.location.hash.substr(1)
if (term.value) {
	term.value += "\n" + e(term.value) + "\n "
} else {
	term.value = "ESC(toggle hold) ENTER(exec selection or current line) \\c(clear console)\n "
}
term.onkeydown = function (evt) {
	if (evt.which === 27) {
		evt.preventDefault()
		hold = !hold	
		term.style.border = "none"
		if (hold) 
			term.style.border = "2px solid blue"
	} else if (evt.which === 13) {
		if (hold) {
			return
		}
		evt.preventDefault()
		var a = term.selectionStart
		var b = term.selectionEnd
		var s = term.value.substring(a, b)
		if (b == a) {
			if (term.value[a] == "\n")
				a -= 1
			a = term.value.lastIndexOf("\n", a)
			if (a == -1)
				a = 0
			b = term.value.indexOf("\n", b)
			if (b == -1)
				b = term.selectionEnd
			s = term.value.substring(a, b)
		}
		if (term.selectionEnd != term.value.length)
			O(s)
		O("\n")
		s = s.trim()
		if (s === "\\c") {
			term.value = " "
		} else if (s == "\\m") {
			term.value += "\n" + xxd() + " "
		} else if (s.length && s[0] == '/') {
			ls(s.substring(1))
		} else {
			e(s)
			return
		}
		P()
	}
}
function O(s) { term.value += s }
function P() { term.value += "\n "; term.scrollTo(0, term.scrollHeight) }

document.getElementById("dropbox").ondragover = function(ev) {
	ev.preventDefault()
}
document.getElementById("dropbox").ondrop = function(ev) {
	ev.preventDefault()
	if (ev.dataTransfer.items) {
		for (var i = 0; i< ev.dataTransfer.items.length; i++) {
			if (ev.dataTransfer.items[i].kind == 'file') {
				var file = ev.dataTransfer.items[i].getAsFile()
				addfile(file)
			}
		}
	} else
		for (var i = 0; i<ev.dataTransfer.files.length; i++)
			addfile(ev.dataTransfer.files[i])
}
function addfile(f) {
	var r = new FileReader()
	r.onload = function() {
		sendfile("/"+f.name, r.result)
	}
	r.readAsArrayBuffer(f)
}
function sendfile(name, buf) {
	var req = new XMLHttpRequest()
	req.onreadystatechange = function() { if (this.readyState == (this.DONE || 4)) { O(req.response+" ");term.scrollTo(0, term.scrollHeight)  } }
	req.open("POST", "")
	req.setRequestHeader("file", name)
	req.send(buf)
}
/* TODO: some kind of readdir or read "." within k?
function ls(name) { // list files (empty name), or show
	if (name.length == 0) {
		for (var name in files)
			O("/"+name+"\n")
		return
	}
	var f = files[name]
	if (f == undefined) {
		O("?")
		return
	}
	var r = new FileReader()
	r.onload = function(f) {
		return function(e) {
			O(e.target.result)
			P()
		}
	}(f)
	r.readAsText(f) // readAsArrayBuffer...
}
*/
</script></body></html>`
