months:`jan`feb`mar`apr`may`jun`jul`aug`sep`oct`nov`dec
t:+`month`amount!(15?months;?15)
m:`jan

show[`;`t]
show[`onchange!{m::`$x};months]
show[`;{t@&m=t`month}]

c:0
show[`type`onclick!(`button;{c+:1});"click me"]
show[`readOnly!1;`c]
