add:I:II{x+y}
/
/sum:I:I{x/r+:i;r}
/           loop r    i    +  →r   i    1    +  →i;i x    <  repeat end r    / w:  x/r+:i;r
/sum:I:I:II: 0340 2001 2002 6a 2101 2002 4101 6a 2202 2000 49 0d00   0b  2001 / k:  +/!x
/ 0340  loop void
/  2001  local.get r  /r:r+i
/  2002  local.get i
/  6a    i32.add
/  2101  local.set r
/  
/  2002  local.get i  /i:i+1;i
/  4101  i32.const 1
/  6a    i32.add
/  2202  local.tee i
/  
/  2000  local.get x  / (i)<x
/  49    i32.lt_u
/  0d00  br_if        / if break loop
/ 0b  end
/ 
/ 2001  local.get r   / return r
/ 0b  end
/ 
/ (module
/  (func $f (param $n i32) (result i32) 
/   (local $r i32)
/   (local $i i32)
/   (loop
/    (set_local $r (i32.add (get_local $r) (get_local $i)))
/    (set_local $i (i32.add (get_local $i) (i32.const 1)))
/    (br_if 0 (i32.lt_u (get_local $i)(get_local $n)))
/   )
/   (get_local $r)
/  )
/  (export "f" (func 0))
/ )
/                               . ..
/ 0000000: 0061 736d 0100 0000 0106 0160 017f 017f  .asm.......`....
/ 0000010: 0302 0100 0503 0100 0107 0701 0373 756d  .............sum
/ 0000020: 0000 0a1e 011c 0102 7f03 4020 0120 026a  ..........@ . .j
/ 0000030: 2101 2002 4101 6a22 0220 0049 0d00 0b20  !. .A.j". .I...
/ 0000040: 010b
