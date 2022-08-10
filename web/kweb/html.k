html:js[`x;"document.body.innerHTML+=x"]
tag:{"<",x,">",(,/y),"</",(x:$x),">"}
h1:tag`h1
li:tag`li
ul:(tag`ul)@li'
td:tag`td
tr:(tag`tr)@td'
table:(tag`table)@tr'

html h1"This is khtml"
html "let's itemize"
html ul("one";"two";"three")
html h1"a table"
html table($'4^?12)
