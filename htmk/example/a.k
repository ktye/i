url["/tablenames"]{x;,/tag["<span class='big link'>";"</span>"]'$tables}
url["/buttons"   ]{x;,/{tag["<span class='big link' hx-get='exec?q=",x,"'>";"</span>";x]}'$buttons}
url["/exec?q=s"  ]{x;$q}
url["/table"     ]{x;t:L;(,/tag["<option>";"</option>"]'fmtt(100&#t)#t)}

fmtt:{ html'"|"/'+{x@\(!|/#'x)}'(,'$!x),'$.x }
tag:{x,z,y}

tables:`L`D
buttons:`abc`def`ghi

L:+`a`b`ceta!(!10;?10;?10a)