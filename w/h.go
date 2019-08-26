package main

const h = `<html><head><script>

</script>
<meta charset="utf-8">
<style>
	body,textarea { font-family: monospace; margin:0pt; }
	textarea {background-color:black;color:white;border:none;resize: none;}
	.col { float: left; width:50%; height:100%; }
	.row:after { content: ""; display: table; clear: both }
</style>
</head><body>
<div id="dropbox">
<div class="row">
	<textarea id="term" class="col"></textarea> /* left */
	<textarea id="edit" class="col"></textarea> /* right */
	<image id="dpy" class="col"></image>        /* right(alt) */
</div></div>
<script>
var term = document.getElementById("term")
var edit = document.getElementById("edit")
var dpy  = document.getElementById("dpy")
function e(k, n, f) {
	var req = new XMLHttpRequest()
	req.onreadystatechange = function() { 
		if (this.readyState == (this.DONE || 4)) { 
			if (req.getResponseHeader('Content-Type') == "image/png") {
				img.src = req.response;O(" ")
			} else {
				O(req.response+" ");term.scrollTo(0, term.scrollHeight)  
			}
		} 
	}
	var a = edit.selectionStart
	var b = edit.selectionEnd
	var s = edit.value.substring(a, b)
	req.open("POST", "")
	req.setRequestHeader("n", n) // file name
	req.setRequestHeader("s", s) // selected text   (only .e)
	req.setRequestHeader("a", a) // selection start (only .e)
	req.setRequestHeader("b", b) // selection end   (only .e)
	req.setRequestHeader("w", dpy.width) 
	req.setRequestHeader("h", dpy.height)
	req.setRequestHeader("k", k) // term value(current line)
	req.send(f)
}
var hold = false
term.value = " "
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
		} else if (s.length && s[0] == '/') {
			// ls(s.substring(1))
			O("TODO ls")
		} else {
			e(s, "", edit.value)
			return
		}
		P()
	}
}
function O(s) { term.value += s }
function P() { term.value += "\n "; term.scrollTo(0, term.scrollHeight) }
func show(e, b) { b?e.style.display="block":e.style.display="none" }
func edit(b) { if(b){show(edit,true);show(dpy,false)}else{show(edit,false);show(display,true)}}

var dropbox = document.getElementById("dropbox")
dropbox.ondragover = function(ev) { ev.preventDefault() }
dropbox.ondrop = function(ev) {
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
		e("", f.name, r.result)
	}
	r.readAsArrayBuffer(f)
}
edit(true)
e("", ".e", "") // receives image response
</script></body></html>`
