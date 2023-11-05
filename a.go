package main

import (
	. "github.com/ktye/wg/module"
)

func init() {
	Memory(1)
	Memory2(1)
	Data(132, "\x00\x01@\x01\x01\x01\x01\t\x10`\x01\x01\x01\x01\x01\t\xc4\xc4\xc4\xc4\xc4\xc4\xc4\xc4\xc4\xc4\x01 \x01\x01\x01\x01\x01BBBBBBBBBBBBBBBBBBBBBBBBBB\x10\t`\x01\x01\x00\xc2\xc2\xc2\xc2\xc2\xc2BBBBBBBBBBBBBBBBBBBB\x10\x01`\x01") // k_test.go: TestClass
	Data(227, ":+-*%&|<>=~!,^#_$?@.':/:\\:vbcisfzldtmdplx00BCISFZLDT0")
	Export(main, Asn, Atx, Cal, cs, dx, Kc, Kf, Ki, kinit, l2, mk, nn, repl, rx, sc, src, Srcp, tp, trap, Val)

	//            0    :    +    -    *    %    &    |    <    >    =10   ~    !    ,    ^    #    _    $    ?    @    .20  '    ':   /    /:   \    \:                  30                       35                       40                       45
	Functions(00, nul, Idy, Flp, Neg, Fst, Sqr, Wer, Rev, Asc, Dsc, Grp, Not, Til, Enl, Srt, Cnt, Flr, Str, Unq, Typ, Val, ech, nyi, rdc, nyi, scn, nyi, lst, Kst, Out, nyi, nyi, Abs, Img, Cnj, Ang, nyi, Uqs, nyi, Tok, Fwh, Las, Exp, Log, Sin, Cos, Prs)
	Functions(64, Asn, Dex, Add, Sub, Mul, Div, Min, Max, Les, Mor, Eql, Mtc, Key, Cat, Cut, Tak, Drp, Cst, Fnd, Atx, Cal, Ech, nyi, Rdc, nyi, Scn, nyi, com, prj, Otu, In, Find, Hyp, Cpx, fdl, Rot, Enc, Dec, nyi, nyi, Bin, Mod, Pow, Lgn)
	Functions(193, tchr, tnms, tvrb, tpct, tvar, tsym, pop)
	Functions(211, Amd, Dmd)

	Functions(220, negi, negf, negz)
	Functions(223, absi, absf, nyi)
	Functions(226, addi, addf, addz)
	Functions(229, subi, subf, subz)
	Functions(232, muli, mulf, mulz)
	Functions(235, divi, divf, divz)
	Functions(238, mini, minf, minz)
	Functions(241, maxi, maxf, maxz)
	Functions(244, modi, sqrf, nyi)

	Functions(247, cmi, cmi, cmi, cmF, cmZ, cmC, cmI, cmI, cmF, cmZ, cmL)
	Functions(258, sum, rd0, prd, rd0, min, max)
	Functions(264, mtC, mtC, mtC, mtF, mtF, mtL)
	Functions(270, inC, inI, inI, inF, inZ)
	Functions(275, exp1, log1, sin1, cos1, pow2)
}

func trap() {
	s := src()
	if srcp == 0 {
		write(Ku(2608)) // 0\n
	} else {
		a := maxi(srcp-30, 0)
		b := mini(nn(s), srcp+30)
		for i := a; i < b; i++ {
			if I8(int32(s)+i) == 10 {
				if i < srcp {
					a = 1 + i
				} else {
					b = i
				}
			}
		}
		Write(0, 0, int32(s)+a, b-a)
		if srcp > a {
			write(Cat(Kc(10), ntake(srcp-a-1, Kc(32))))
		}
	}
	write(Ku(2654)) // ^\n
	panic(srcp)
}
func Srcp() int32 { return srcp }
