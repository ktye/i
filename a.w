def //* ;; %
def ! unreachable
def ret return
def con i32.const
def 0i con 0
def 1i con 1
def 4i con 4
def 32i con 32
def i. i32.load
def i: i32.store
def = i32.eq
def != i32.ne
def < i32.lt_u
def <s i32.lt_s
def > i32.gt_u
def >s i32.gt_s
def >= i32.ge_u
def >=s i32.ge_s
def =0 i32.eqz
def + i32.add
def - i32.sub
def * i32.mul
def / i32.div_u
def /s i32.div_s
def % i32.rem_u
def %s i32.rem_s
def while* block loop % br_if 1
def do* % br 0 end end
def ndo* $n eqz if loop % $n 1 sub tee $n if_br 0 end end
def & i32.and
def | i32.or
def ^ i32.xor
def >> i32.shl
def << i32.shr_u
def <<s i32.shr_s
def (f1* (func $% (param $x i32) (result i32)
def (f2* (func $% (param $x i32) (param $y i32) (result i32)
def export* (export "%" (func $%))
def var' (local $% i32)

// .x .y + :.z
// (f2 fname

(f2 add
	var a b c
	.x .y + :.c
)

export add
