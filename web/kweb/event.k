months:`jan`feb`mar`apr`may`jun`jul`aug`sep`oct`nov`dec
t:+`month`amount!(30?months;?30)
m:`jan
         show[`;`t]
dropdown:show[`;months]              /show returns the id
filter:  show[`fil;{t@&m=t`month}]

/when the dropdown is changed we call f
js[0#`;"document.getElementById('",($dropdown),"').onchange=function(e){K.Kx('f',K.Ks(e.target.value));document.update()}"][]
f:{m::x}
