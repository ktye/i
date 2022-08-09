months:`jan`feb`mar`apr`may`jun`jul`aug`sep`oct`nov`dec
t:+`month`amount!(30?months;?30)
m:`jan
         link[`;`t]
dropdown:show[`;months]              /show,link,expr return the id
filter:  expr[`fil;"t@&m=t`month"]

/when the dropdown is changed we call f
js[0#`;"document.getElementById('",($dropdown),"').onchange=function(e){K.Kx('f',K.Ks(e.target.value));document.update()}"][]
f:{m::x}
