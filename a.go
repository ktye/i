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
	Functions(64, Asn, Dex, Add, Sub, Mul, Div, Min, Max, Les, Mor, Eql, Mtc, Key, Cat, Cut, Tak, Drp, Cst, Fnd, Atx, Cal, Ech, Whl, Rdc, nyi, Scn, nyi, com, prj, Otu, In, Find, Hyp, Cpx, nyi, Rot, Enc, Dec, nyi, nyi, Bin, Mod, Pow, Lgn, nyi, nyi, Rtp)
	Functions(193, tchr, tnms, tvrb, tpct, tvar, tsym)
	Functions(211, Amd, Dmd)

	Functions(220, negi, negf, negz) //220
	Functions(223, absi, absf, nyi)  //227
	Functions(226, addi, addf, addz) //234
	Functions(229, subi, subf, subz) //245
	Functions(232, muli, mulf, mulz) //256
	Functions(235, divi, divf, divz) //267
	Functions(238, mini, minf, minz) //278
	Functions(241, maxi, maxf, maxz) //289
	Functions(244, modi, sqrf, nyi)  //300

	Functions(247, eqi, eqf, eqz, eqC, eqI, eqI, eqF, eqZ) //308
	Functions(255, lti, ltf, ltz, guC, guI, guI, guF, guZ) //323
	Functions(263, gti, gtf, gtz, gdC, gdI, gdI, gdF, gdZ) //338

	Functions(271, guC, guC, guI, guI, guF, guZ, guL, gdC, gdC, gdI, gdI, gdF, gdZ, gdL) //353

	Functions(285, sum, rd0, prd, rd0, min, max) //367
	Functions(291, sums, rd0, prds, rd0, rd0)    //374
}

// kvm
var IP int32 //instruction pointer
var RS K     //return stack

func exec(x K) K {
	//println("exec>sp", sp)
	//println("exec IP", IP)
	println("exec#", nn(x), sK(x))
	if nn(x) < 2 {
		return 0
	}

	rpush(IP)
	rpush(0)
	IP = int32(x)
	for step() != 0 {
	}

	//rpop() //0
	IP = rpop()
	r := pop()

	dx(x)
	//println("st", K(I64(256)))
	//println("exec<sp", sp)
	println("exec return", sK(r))
	return r
}
func step() int32 {
	var x, a, b, c K
	for {
		x = K(I64(IP)) //fetch instruction
		//fmt.Println("fetch", sK(x))
		if x == 320 {
			//rest()
			//println("retn")
			IP = rpop()
			//println("jmp IP", IP)
			if IP != 0 {
				//println("    IP", IP, ":", sK(K(Lt)<<59|K(IP)))
			}
			if IP == 0 {
				//println("halt")
				return 0 //halt
			}
			continue
		}
		break
	}

	/*
		if tp(x) != 0 {
			println("u", tp(x), sK(x))
		} else {
			println("v", int32(x))
		}
	*/

	if tp(x) != 0 { //noun
		push(rx(x))
		//println("push", sK(x))
	} else {
		v := int32(x)
		//println("v", v, v>>6)
		switch v >> 6 {
		case 0: //monadic
			push(Func[v].(f1)(pop()))
		case 1: //dyadic
			a = pop()
			push(Func[v].(f2)(a, pop()))
		case 2: //dyadic indirect
			a = pop()
			b = pop()
			push(Cal(a, l2(b, pop())))
		case 3: //tetradic
			a = pop()
			b = pop()
			c = pop()
			push(Func[v].(f4)(a, b, c, pop()))
		case 4: // drop
			if v == 256 {
				dx(pop())
			}
			if v == 257 {
				cnd(pop())
			}
		case 5:
			trap(Nyi) //todo RET
		default: //quoted
			push(x - 448)
		}
	}
	IP += 8
	return 1
}
func rpush(x int32) {
	//fmt.Println("rpush #", nn(RS)+1, RS, x)
	RS = cat1(RS, Ki(x))
	/*
		if RS != 10952754293765050768 {
			println("RS push", RS, "relocation")
			panic("RS")
		}
	*/
}
func rpop() int32 {
	//fmt.Println("rpop  #", nn(RS)-1, RS)
	if nn(RS) == 13 {
		panic("rs underflow")
	}
	r := I32(int32(RS) + 4*(nn(RS)-1))

	RS = ndrop(-1, RS)
	/*
		if RS != 10952754293765050768 {
			println("RS pop", RS, "relocation")
			panic("RS")
		}
	*/
	//fmt.Println("rpop ", r, "#", nn(RS))
	return r
}
