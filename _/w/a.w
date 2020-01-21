/           x    y    +
add:I:II::  2000 2001 6a         /add two integers

/          block loop  x    1    -  tee  0    =  ifbr     r    1    +  br0(continue) end  r    / x/r+:i;r
sum:I:I:I: 0240  03400 2000 4101 6b 2200 4100 46 0d01 6a  2001 4101 6a 0c00          0b0b 2001 / sum loop (+/!x)