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
<div class="row"><textarea id="term" class="col"></textarea><image id="dpy" class="col"></image></div></div>
<script>
function e(s) {
	var img = document.getElementById("dpy")
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
	req.open("POST", "")
	req.setRequestHeader("width", img.width)
	req.setRequestHeader("height", img.width)
	req.send(s)
}
var hold = false
var term = document.getElementById("term")
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
			e(s)
			return
		}
		P()
	}
}
function O(s) { term.value += s }
function P() { term.value += "\n "; term.scrollTo(0, term.scrollHeight) }

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
e("") // first call always responds with an image
</script></body></html>`
