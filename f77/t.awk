BEGIN{FS=" /"}
$2~/3a300 0a/{print $1" /2.999999a300 0a"; next}
$2~/6a240/{print $1" /5.999999a240"; next}
$2~/(3a80;4a80 5a80)/{print $1" /(2.999999a80;3.999999a80 4.999999a80)"; next}
$2~/1a80 1a280 1a90.2 1a/{print $1" /0.999999a80 0.999999a280 1a90.2 1a"; next}
$2~/3a70/{print $1" /2.999999a70"; next}
$2~/4a70/{print $1" /3.999999a70"; next}
{print}
