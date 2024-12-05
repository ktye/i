(func $seqI (param $xp i32) (param $e i32) (local $v v128)
 v128.const i32x4 0 1 2 3
 local.set $v
 loop
  local.get $xp
  local.get $v
  v128.store
  local.get $v
  i32.const 4
  i32x4.splat
  i32x4.add
  local.set $v
  local.get $xp
  i32.const 16
  i32.add
  local.tee $xp
  local.get $e
  i32.lt_u
  br_if 0
 end)
(func $negI (param $xp i32) (param $e i32)
 loop
  local.get $xp
  local.get $xp
  v128.load
  i32x4.neg
  v128.store
  local.get $xp
  i32.const 16
  i32.add
  local.tee $xp
  local.get $e
  i32.lt_s
  br_if 0
 end)
(func $negF (param $xp i32) (param $e i32)
 loop
  local.get $xp
  local.get $xp
  v128.load
  f64x2.neg
  v128.store
  local.get $xp
  i32.const 16
  i32.add
  local.tee $xp
  local.get $e
  i32.lt_s
  br_if 0
 end)
 (func $absI (param $xp i32) (param $e i32)
 loop
  local.get $xp
  local.get $xp
  v128.load
  i32x4.abs
  v128.store
  local.get $xp
  i32.const 16
  i32.add
  local.tee $xp
  local.get $e
  i32.lt_s
  br_if 0
 end)
(func $absF (param $xp i32) (param $e i32)
 loop
  local.get $xp
  local.get $xp
  v128.load
  f64x2.abs
  v128.store
  local.get $xp
  i32.const 16
  i32.add
  local.tee $xp
  local.get $e
  i32.lt_s
  br_if 0
 end)
(func $sqrF (param $xp i32) (param $e i32)
 loop
  local.get $xp
  local.get $xp
  v128.load
  f64x2.sqrt
  v128.store
  local.get $xp
  i32.const 16
  i32.add
  local.tee $xp
  local.get $e
  i32.lt_s
  br_if 0
 end)
 
(func $ltC (param $xp i32) (param $yp i32) (param $rp i32) (param $ep i32) (local $z v128)
 loop
  local.get $rp
  local.get $xp
  v128.load
  local.get $yp
  v128.load
  i8x16.lt_s
  v128.const i8x16 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1
  v128.and
  local.tee $z
  
  v128.const i32x4 0 0 0 0
  i8x16.shuffle 0 16 16 16 1 16 16 16 2 16 16 16 3 16 16 16 
  v128.store
  
  local.get $rp
  i32.const 16
  i32.add
  local.tee $rp
  local.get $z
  v128.const i32x4 0 0 0 0
  i8x16.shuffle 4 16 16 16 5 16 16 16 6 16 16 16 7 16 16 16 
  v128.store
  
  local.get $rp
  i32.const 16
  i32.add
  local.tee $rp
  local.get $z
  v128.const i32x4 0 0 0 0
  i8x16.shuffle 8 16 16 16 9 16 16 16 10 16 16 16 11 16 16 16 
  v128.store
  
  local.get $rp
  i32.const 16
  i32.add
  local.tee $rp
  local.get $z
  v128.const i32x4 0 0 0 0
  i8x16.shuffle 16 16 16 12 16 16 16 13 16 16 16 14 16 16 16 15
  v128.store
  
  local.get $xp
  i32.const 16
  i32.add
  local.set $xp
  local.get $yp
  i32.const 16
  i32.add
  local.set $yp
  local.get $rp
  i32.const 16
  i32.add
  local.tee $rp
  local.get $ep
  i32.lt_s
  br_if 0
 end)
(func $eqC (param $xp i32) (param $yp i32) (param $rp i32) (param $ep i32) (local $z v128)
 loop
  local.get $rp
  local.get $xp
  v128.load
  local.get $yp
  v128.load
  i8x16.eq
  v128.const i8x16 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1
  v128.and
  local.tee $z
  
  v128.const i32x4 0 0 0 0
  i8x16.shuffle 0 16 16 16 1 16 16 16 2 16 16 16 3 16 16 16 
  v128.store
  
  local.get $rp
  i32.const 16
  i32.add
  local.tee $rp
  local.get $z
  v128.const i32x4 0 0 0 0
  i8x16.shuffle 4 16 16 16 5 16 16 16 6 16 16 16 7 16 16 16 
  v128.store
  
  local.get $rp
  i32.const 16
  i32.add
  local.tee $rp
  local.get $z
  v128.const i32x4 0 0 0 0
  i8x16.shuffle 8 16 16 16 9 16 16 16 10 16 16 16 11 16 16 16 
  v128.store
  
  local.get $rp
  i32.const 16
  i32.add
  local.tee $rp
  local.get $z
  v128.const i32x4 0 0 0 0
  i8x16.shuffle 16 16 16 12 16 16 16 13 16 16 16 14 16 16 16 15
  v128.store
  
  local.get $xp
  i32.const 16
  i32.add
  local.set $xp
  local.get $yp
  i32.const 16
  i32.add
  local.set $yp
  local.get $rp
  i32.const 16
  i32.add
  local.tee $rp
  local.get $ep
  i32.lt_s
  br_if 0
 end)
(func $gtC (param $xp i32) (param $yp i32) (param $rp i32) (param $ep i32) (local $z v128)
 loop
  local.get $rp
  local.get $xp
  v128.load
  local.get $yp
  v128.load
  i8x16.gt_s
  v128.const i8x16 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1
  v128.and
  local.tee $z
  
  v128.const i32x4 0 0 0 0
  i8x16.shuffle 0 16 16 16 1 16 16 16 2 16 16 16 3 16 16 16 
  v128.store
  
  local.get $rp
  i32.const 16
  i32.add
  local.tee $rp
  local.get $z
  v128.const i32x4 0 0 0 0
  i8x16.shuffle 4 16 16 16 5 16 16 16 6 16 16 16 7 16 16 16 
  v128.store
  
  local.get $rp
  i32.const 16
  i32.add
  local.tee $rp
  local.get $z
  v128.const i32x4 0 0 0 0
  i8x16.shuffle 8 16 16 16 9 16 16 16 10 16 16 16 11 16 16 16 
  v128.store
  
  local.get $rp
  i32.const 16
  i32.add
  local.tee $rp
  local.get $z
  v128.const i32x4 0 0 0 0
  i8x16.shuffle 16 16 16 12 16 16 16 13 16 16 16 14 16 16 16 15
  v128.store
  
  local.get $xp
  i32.const 16
  i32.add
  local.set $xp
  local.get $yp
  i32.const 16
  i32.add
  local.set $yp
  local.get $rp
  i32.const 16
  i32.add
  local.tee $rp
  local.get $ep
  i32.lt_s
  br_if 0
 end)
 
(func $ltI (param $xp i32) (param $yp i32) (param $rp i32) (param $ep i32)
 loop
  local.get $rp
  local.get $xp
  v128.load
  local.get $yp
  v128.load
  i32x4.lt_s
  v128.const i32x4 1 1 1 1
  v128.and
  v128.store
  local.get $xp
  i32.const 16
  i32.add
  local.set $xp
  local.get $yp
  i32.const 16
  i32.add
  local.set $yp
  local.get $rp
  i32.const 16
  i32.add
  local.tee $rp
  local.get $ep
  i32.lt_s
  br_if 0
 end)
(func $gtI (param $xp i32) (param $yp i32) (param $rp i32) (param $ep i32)
 loop
  local.get $rp
  local.get $xp
  v128.load
  local.get $yp
  v128.load
  i32x4.gt_s
  v128.const i32x4 1 1 1 1
  v128.and
  v128.store
  local.get $xp
  i32.const 16
  i32.add
  local.set $xp
  local.get $yp
  i32.const 16
  i32.add
  local.set $yp
  local.get $rp
  i32.const 16
  i32.add
  local.tee $rp
  local.get $ep
  i32.lt_s
  br_if 0
 end)
(func $eqI (param $xp i32) (param $yp i32) (param $rp i32) (param $ep i32)
 loop
  local.get $rp
  local.get $xp
  v128.load
  local.get $yp
  v128.load
  i32x4.eq
  v128.const i32x4 1 1 1 1
  v128.and
  v128.store
  local.get $xp
  i32.const 16
  i32.add
  local.set $xp
  local.get $yp
  i32.const 16
  i32.add
  local.set $yp
  local.get $rp
  i32.const 16
  i32.add
  local.tee $rp
  local.get $ep
  i32.lt_s
  br_if 0
 end)
