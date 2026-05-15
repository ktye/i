BEGIN{FS=" /"}
$2~/3a300 0a/{next}
$2~/6a240/{next}
$2~/(3a80;4a80 5a80)/{next}
$2~/1a80 1a280 1a90.2 1a/{next}
$2~/3a70/{next}
$2~/4a70/{next}
{print}
