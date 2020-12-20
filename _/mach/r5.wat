;; risc-v virtual machine in webassembly
;; - instantiate module
;; - grow exported memory as needed
;; - copy risc-v program to memory at 384 (first 32*4 bytes are i32 registers, then 32*8 f64 registers)
;; - call "jmp" with pc argument (linear memory addr)
(module
 (type (func (param i32 i32) (result i32)))
 (type (func (param f64 f64) (result f64)))
 (type (func (param i32 i32) (result f64)))
 (type (func (param i32 i32)))
 (type (func (param i32 f64)))
 (func $jmp (param $pc i32) 
  (local $i  i32) (local $op i32) (local $rd i32) (local $r1 i32) (local $r2 i32) 
  (local $f3 i32) (local $f7 i32) (local $x  i32) (local $y  i32) (local $z  i32) (local $im i32) (local $is i32) (local $ib i32)
  (local.set $i (local.get $pc))
  
  (local.set $op (i32.and (local.get $i) (i32.const 127)))
  (local.set $rd (i32.and (i32.shr_u (local.get $i) (i32.const  7)) (i32.const  31)))
  (local.set $r1 (i32.and (i32.shr_u (local.get $i) (i32.const 15)) (i32.const  31)))
  (local.set $r2 (i32.and (i32.shr_u (local.get $i) (i32.const 20)) (i32.const  31)))
  (local.set $f3 (i32.and (i32.shr_u (local.get $i) (i32.const 12)) (i32.const   7)))
  (local.set $f7 (i32.and (i32.shr_u (local.get $i) (i32.const 25)) (i32.const 127)))
  (local.set $im (i32.shr_s (local.get $i) (i32.const 20)))
  (local.set $is (i32.or (local.get $rd) (i32.shr_s (local.get $i) (i32.const 25))))
  (local.set $ib (i32.or (i32.and   (i32.const 30) (i32.shr_u (local.get $i)  (i32.const 7)))
                 (i32.or (i32.shl   (i32.and (i32.const 128)  (local.get $i)) (i32.const 4)) 
                 (i32.or (i32.shr_u (i32.and (i32.const 0x7e000000) (local.get $i)) (i32.const 20)) 
                         (i32.shr_s (i32.and (i32.const 0x80000000) (local.get $i)) (i32.const 19))))))
 
  
  (if       (i32.eq (local.get $op) (i32.const   3)) (then (i32.store (local.get $rd) (call_indirect (type 0) (local.get $f3) (local.get $r1) (local.get $im)))) ;; ld 0..7
  (else (if (i32.eq (local.get $op) (i32.const   7)) (then (f64.store offset=128 (local.get $rd) (call $fld (local.get $r1) (local.get $im))))                   ;; fld
  (else (if (i32.eq (local.get $op) (i32.const  19)) (then (i32.store (local.get $rd) (call_indirect (type 0) (i32.add (local.get $f3) (i32.const 8)) (i32.load (local.get $r1)) (local.get $im)))) ;; addi 8..15
  (else (if (i32.eq (local.get $op) (i32.const  35)) (then (call_indirect (type 3) (i32.add (local.get $f3) (i32.const 15)) (i32.add (local.get $r1) (local.get $is)) (local.get $r2)))             ;; sb  15..18
  (else (if (i32.eq (local.get $op) (i32.const  39)) (then (call $fsd (i32.add (local.get $r1) (local.get $is)) (f64.load offset=128 (local.get $rd))))          ;; fsd
  (else (if (i32.eq (local.get $op) (i32.const  51)) (then (i32.store (local.get $rd) (call_indirect (type 0) (i32.add (local.get $f3)
                                                                            (if (result i32) (i32.eqz (local.get $f7)) (then (i32.const 8))                      ;; add 8..15
	                                                                    (else (if (result i32) (i32.eq (local.get $f7) (i32.const 1)) (then (i32.const 19))  ;; mul 19..26
	                                                                    (else (i32.const 27))))))                                                            ;; sub 27..34 
                                                                                         (i32.load (local.get $r1)) (i32.load (local.get $r2)))))
  (else (if (i32.eq (local.get $op) (i32.const  99)) (then (if (call_indirect (type 0) (i32.add (local.get $f3) (i32.const 35)) (local.get $r1) (local.get $r2))
              (then (local.set $pc (i32.sub (i32.add (local.get $pc) (local.get $ib)) (i32.const 4))))))                                                         ;; beq 35..42
         
  (else (unreachable)))))))))))))))
  
  (local.set $pc (i32.add (local.get $pc) (i32.const 4)))
 )
 (func $xxx  (type 0) unreachable)
 (func $lb   (type 0) (i32.load8_s offset=1152 (i32.add (local.get 0) (local.get 1))))
 (func $lw   (type 0) (i32.load    offset=1152 (i32.add (local.get 0) (local.get 1))))
 (func $lbu  (type 0) (i32.load8_u offset=1152 (i32.add (local.get 0) (local.get 1))))
 (func $sb   (type 3) (i32.store8  offset=1152 (local.get 0) (local.get 1)))
 (func $sw   (type 3) (i32.store   offset=1152 (local.get 0) (local.get 1)))
 (func $fsd  (type 4) (f64.store   offset=1152 (local.get 0) (local.get 1)))
 (func $add  (type 0) (i32.add   (local.get 0) (local.get 1)))
 (func $sub  (type 0) (i32.sub   (local.get 0) (local.get 1)))
 (func $fadd (type 1) (f64.add   (local.get 0) (local.get 1)))
 (func $shl  (type 0) (i32.shl   (local.get 0) (local.get 1)))
 (func $shr  (type 0) (i32.shr_u (local.get 0) (local.get 1)))
 (func $sra  (type 0) (i32.shr_s (local.get 0) (local.get 1)))
 (func $xor  (type 0) (i32.xor   (local.get 0) (local.get 1)))
 (func $or   (type 0) (i32.or    (local.get 0) (local.get 1)))
 (func $and  (type 0) (i32.and   (local.get 0) (local.get 1)))
 (func $mul  (type 0) (i32.mul   (local.get 0) (local.get 1)))
 (func $div  (type 0) (i32.div_s (local.get 0) (local.get 1)))
 (func $divu (type 0) (i32.div_u (local.get 0) (local.get 1)))
 (func $rem  (type 0) (i32.rem_s (local.get 0) (local.get 1)))
 (func $remu (type 0) (i32.rem_u (local.get 0) (local.get 1)))
 (func $beq  (type 0) (i32.eq    (local.get 0) (local.get 1)))
 (func $bne  (type 0) (i32.ne    (local.get 0) (local.get 1)))
 (func $blt  (type 0) (i32.lt_s  (local.get 0) (local.get 1)))
 (func $bltu (type 0) (i32.lt_u  (local.get 0) (local.get 1)))
 (func $bge  (type 0) (i32.ge_s  (local.get 0) (local.get 1)))
 (func $bgeu (type 0) (i32.ge_u  (local.get 0) (local.get 1)))
 (func $fld  (type 2) (f64.load offset=128 (i32.add (local.get 0) (local.get 1))))
 (memory 1) (export "mem" (memory 0))
 (table 42 funcref)
 (export "jmp" (func 0))
 (elem (i32.const 0) func $lb   $xxx  $lw  $xxx $lbu  $xxx  $xxx $xxx                  ;;  0..7
                          $add  $shl  $xxx $xxx $xor  $shr  $or  $and $sb $xxx $sw     ;;  7..18
			  $mul  $xxx  $xxx $xxx $div  $divu $rem $remu                 ;; 19..26
			  $sub  $xxx  $xxx $xxx $xxx  $sra  $xxx $xxx                  ;; 27..34
			  $beq  $bne  $xxx $xxx $blt  $bge  $bltu $bgeu  )             ;; 35..42
)
