package main

import (
	. "github.com/ktye/wg/module"
)

func init() {
	Memory(1)
	Memory2(1)
	Data(132, "\x00\x01@\x01\x01\x01\x01\t\x10`\x01\x01\x01\x01\x01\t\xc4\xc4\xc4\xc4\xc4\xc4\xc4\xc4\xc4\xc4\x01 \x01\x01\x01\x01\x01BBBBBBBBBBBBBBBBBBBBBBBBBB\x10\t`\x01\x01\x00\xc2\xc2\xc2\xc2\xc2\xc2BBBBBBBBBBBBBBBBBBBB\x10\x01`\x01") // k_test.go: TestClass
	Data(228, ":+-*%&|<>=~!,^#_$?@.':/:\\:")
	Data(520, "vbcisfzldtmdplx00BCISFZLDT") //546
	Export(main, Asn, Atx, Cal, cs, dx, Kc, Kf, Ki, kinit, l2, mk, nn, repl, rx, sc, src, Srcp, tp, trap, Val)

	//            0    :    +    -    *    %    &    |    <    >    =10   ~    !    ,    ^    #    _    $    ?    @    .20  '    ':   /    /:   \    \:                  30                       35                       40                       45
	Functions(00, nul, Idy, Flp, Neg, Fst, Sqr, Wer, Rev, Asc, Dsc, Grp, Not, Til, Enl, Srt, Cnt, Flr, Str, Unq, Typ, Val, ech, nyi, rdc, nyi, scn, nyi, lst, Kst, Out, nyi, nyi, Abs, Img, Cnj, Ang, nyi, Uqs, nyi, Cos, Fwh, Las, Exp, Log, Sin, Tok, Prs)
	Functions(64, Asn, Dex, Add, Sub, Mul, Div, Min, Max, Les, Mor, Eql, Mtc, Key, Cat, Cut, Tak, Drp, Cst, Fnd, Atx, Cal, Ech, nyi, Rdc, nyi, Scn, nyi, com, prj, Otu, In, Find, Hyp, Cpx, nyi, Rot, Enc, Dec, nyi, nyi, Bin, Mod, Pow, Lgn, nyi, nyi, Rtp)
	Functions(193, tchr, tnms, tvrb, tpct, tvar, tsym)
	Functions(211, Amd, Dmd)

	Functions(220, negi, negf, negz)
	Functions(227, absi, absf)
	Functions(234, addi, addf, addz)
	Functions(245, subi, subf, subz)
	Functions(256, muli, mulf, mulz)
	Functions(267, divi, divf, divz)
	Functions(278, mini, minf, minz)
	Functions(289, maxi, maxf, maxz)
	Functions(300, modi, sqrf, nyi)

	Functions(308, eqi, eqf, eqz, eqC, eqI, eqI, eqF, eqZ)
	Functions(323, lti, ltf, ltz, guC, guI, guI, guF, guZ)
	Functions(338, gti, gtf, gtz, gdC, gdI, gdI, gdF, gdZ)

	Functions(353, guC, guC, guI, guI, guF, guZ, guL, gdC, gdC, gdI, gdI, gdF, gdZ, gdL)

	Functions(367, sum, rd0, prd, rd0, min, max)
	Functions(374, sums, rd0, prds, rd0, rd0)
}
