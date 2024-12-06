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
 
(func $storeB (param $rp i32) (param $v v128) (result i32)
 local.get $rp
 local.get $v
 v128.const i32x4 0 0 0 0
 i8x16.shuffle 0 16 16 16 1 16 16 16 2 16 16 16 3 16 16 16
 v128.store
 
 local.get $rp
 i32.const 16
 i32.add
 local.tee $rp
 local.get $v
 v128.const i32x4 0 0 0 0
 i8x16.shuffle 4 16 16 16 5 16 16 16 6 16 16 16 7 16 16 16 
 v128.store
 
 local.get $rp
 i32.const 16
 i32.add
 local.tee $rp
 local.get $v
 v128.const i32x4 0 0 0 0
 i8x16.shuffle 8 16 16 16 9 16 16 16 10 16 16 16 11 16 16 16 
 v128.store
  
 local.get $rp
 i32.const 16
 i32.add
 local.tee $rp
 local.get $v
 v128.const i32x4 0 0 0 0
 i8x16.shuffle 16 16 16 12 16 16 16 13 16 16 16 14 16 16 16 15
 v128.store
 
 local.get $rp
 i32.const 16
 i32.add
)

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
  
  local.get $xp
  i32.const 16
  i32.add
  local.set $xp
  local.get $yp
  i32.const 16
  i32.add
  local.set $yp
  
  call $storeB
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
 
(func $ltcC (param $xp i32) (param $yp i32) (param $rp i32) (param $ep i32) (local $x v128)
 local.get $xp
 v128.load8_splat
 local.set $x
 loop
  local.get $rp
  local.get $x
  local.get $yp
  v128.load
  i8x16.lt_s
  v128.const i8x16 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1
  v128.and
  local.get $yp
  i32.const 16
  i32.add
  local.set $yp
  call $storeB
  local.tee $rp
  local.get $ep
  i32.lt_s
  br_if 0
 end)
(func $eqcC (param $xp i32) (param $yp i32) (param $rp i32) (param $ep i32) (local $x v128)
 local.get $xp
 v128.load8_splat
 local.set $x
 loop
  local.get $rp
  local.get $x
  local.get $yp
  v128.load
  i8x16.eq
  v128.const i8x16 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1
  v128.and
  local.get $yp
  i32.const 16
  i32.add
  local.set $yp
  call $storeB
  local.tee $rp
  local.get $ep
  i32.lt_s
  br_if 0
 end)
(func $gtcC (param $xp i32) (param $yp i32) (param $rp i32) (param $ep i32) (local $x v128)
 local.get $xp
 v128.load8_splat
 local.set $x
 loop
  local.get $rp
  local.get $x
  local.get $yp
  v128.load
  i8x16.gt_s
  v128.const i8x16 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1
  v128.and
  local.get $yp
  i32.const 16
  i32.add
  local.set $yp
  call $storeB
  local.tee $rp
  local.get $ep
  i32.lt_s
  br_if 0
 end)
(func $ltiI (param $xp i32) (param $yp i32) (param $rp i32) (param $ep i32) (local $x v128)
 local.get $xp
 v128.load32_splat
 local.set $x
 loop
  local.get $rp
  local.get $x
  local.get $yp
  v128.load
  i32x4.lt_s
  v128.const i32x4 1 1 1 1
  v128.and
  v128.store
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
 
 
(func $eqiI (param $xp i32) (param $yp i32) (param $rp i32) (param $ep i32) (local $x v128)
 local.get $xp
 v128.load32_splat
 local.set $x
 loop
  local.get $rp
  local.get $x
  local.get $yp
  v128.load
  i32x4.eq
  v128.const i32x4 1 1 1 1
  v128.and
  v128.store
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
(func $gtiI (param $xp i32) (param $yp i32) (param $rp i32) (param $ep i32) (local $x v128)
 local.get $xp
 v128.load32_splat
 local.set $x
 loop
  local.get $rp
  local.get $x
  local.get $yp
  v128.load
  i32x4.gt_s
  v128.const i32x4 1 1 1 1
  v128.and
  v128.store
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
(func $addI (param $x i32) (param $y i32) (param $r i32) (param $e i32)
 loop
  local.get $r
  local.get $x
  v128.load
  local.get $y
  v128.load
  i32x4.add
  v128.store
  local.get $x
  i32.const 16
  i32.add
  local.set $x
  local.get $y
  i32.const 16
  i32.add
  local.set $y
  local.get $r
  i32.const 16
  i32.add
  local.tee $r
  local.get $e
  i32.lt_s
  br_if 0
 end)
(func $addiI (param $x i32) (param $y i32) (param $r i32) (param $e i32) (local $v v128)
 local.get $x
 v128.load32_splat
 local.set $v
 loop
  local.get $r
  local.get $v
  local.get $y
  v128.load
  i32x4.add
  v128.store
  local.get $y
  i32.const 16
  i32.add
  local.set $y
  local.get $r
  i32.const 16
  i32.add
  local.tee $r
  local.get $e
  i32.lt_s
  br_if 0
 end)
(func $subI (param $x i32) (param $y i32) (param $r i32) (param $e i32)
 loop
  local.get $r
  local.get $x
  v128.load
  local.get $y
  v128.load
  i32x4.sub
  v128.store
  local.get $x
  i32.const 16
  i32.add
  local.set $x
  local.get $y
  i32.const 16
  i32.add
  local.set $y
  local.get $r
  i32.const 16
  i32.add
  local.tee $r
  local.get $e
  i32.lt_s
  br_if 0
 end)
(func $subiI (param $x i32) (param $y i32) (param $r i32) (param $e i32) (local $v v128)
 local.get $x
 v128.load32_splat
 local.set $v
 loop
  local.get $r
  local.get $v
  local.get $y
  v128.load
  i32x4.sub
  v128.store
  local.get $y
  i32.const 16
  i32.add
  local.set $y
  local.get $r
  i32.const 16
  i32.add
  local.tee $r
  local.get $e
  i32.lt_s
  br_if 0
 end)
(func $mulI (param $x i32) (param $y i32) (param $r i32) (param $e i32)
 loop
  local.get $r
  local.get $x
  v128.load
  local.get $y
  v128.load
  i32x4.mul
  v128.store
  local.get $x
  i32.const 16
  i32.add
  local.set $x
  local.get $y
  i32.const 16
  i32.add
  local.set $y
  local.get $r
  i32.const 16
  i32.add
  local.tee $r
  local.get $e
  i32.lt_s
  br_if 0
 end)
(func $muliI (param $x i32) (param $y i32) (param $r i32) (param $e i32) (local $v v128)
 local.get $x
 v128.load32_splat
 local.set $v
 loop
  local.get $r
  local.get $v
  local.get $y
  v128.load
  i32x4.mul
  v128.store
  local.get $y
  i32.const 16
  i32.add
  local.set $y
  local.get $r
  i32.const 16
  i32.add
  local.tee $r
  local.get $e
  i32.lt_s
  br_if 0
 end)
(func $minI (param $x i32) (param $y i32) (param $r i32) (param $e i32)
 loop
  local.get $r
  local.get $x
  v128.load
  local.get $y
  v128.load
  i32x4.min_s
  v128.store
  local.get $x
  i32.const 16
  i32.add
  local.set $x
  local.get $y
  i32.const 16
  i32.add
  local.set $y
  local.get $r
  i32.const 16
  i32.add
  local.tee $r
  local.get $e
  i32.lt_s
  br_if 0
 end)
(func $miniI (param $x i32) (param $y i32) (param $r i32) (param $e i32) (local $v v128)
 local.get $x
 v128.load32_splat
 local.set $v
 loop
  local.get $r
  local.get $v
  local.get $y
  v128.load
  i32x4.min_s
  v128.store
  local.get $y
  i32.const 16
  i32.add
  local.set $y
  local.get $r
  i32.const 16
  i32.add
  local.tee $r
  local.get $e
  i32.lt_s
  br_if 0
 end)
(func $maxI (param $x i32) (param $y i32) (param $r i32) (param $e i32)
 loop
  local.get $r
  local.get $x
  v128.load
  local.get $y
  v128.load
  i32x4.max_s
  v128.store
  local.get $x
  i32.const 16
  i32.add
  local.set $x
  local.get $y
  i32.const 16
  i32.add
  local.set $y
  local.get $r
  i32.const 16
  i32.add
  local.tee $r
  local.get $e
  i32.lt_s
  br_if 0
 end)
(func $maxiI (param $x i32) (param $y i32) (param $r i32) (param $e i32) (local $v v128)
 local.get $x
 v128.load32_splat
 local.set $v
 loop
  local.get $r
  local.get $v
  local.get $y
  v128.load
  i32x4.max_s
  v128.store
  local.get $y
  i32.const 16
  i32.add
  local.set $y
  local.get $r
  i32.const 16
  i32.add
  local.tee $r
  local.get $e
  i32.lt_s
  br_if 0
 end)