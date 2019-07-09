package main

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

type (
	l  = []interface{}
	d  = [2]interface{}
	iv = []int
	sv = []string
)

func TestIni(t *testing.T) {
	t.Skip()
	ini()
	st := Stats()
	if st.UsedBlocks() != 1 {
		t.Fatal()
	}
	//mk(1, 9000)
	//pfl()
	//xxd()
}

func TestK(t *testing.T) {
	//t.Skip()
	testCases := []struct {
		x, r s
	}{
		//{"x:3 3#!9;x[1 2;0]", "0 1 2"}, // TODO matrix indexing
		//{"x:3 3#!9;x[1 2;1 2]:0;x", "(0 1 2;3 0 0;6 0 0)"},
		//{"x:3 3#!9;x[1 2;1 2]:2 2#1+!4", "(0 1 2;3 1 2;6 3 4)"},
		//{"x:3 3#!9;x[1 2;1 2]*:10", "(0 1 2;3 40 50;6 70 80)"},
		{"1", "1"},
		{"`a", "`a"},
		{"`a`b", "`a`b"},
		{"*1 2 3", "1"},
		{"%4f", "0.25"},
		{"%4 5", "0.25 0.2"},
		{"%0i1", "0i-1"},
		{"1 2,4 5", "1 2 4 5"},
		{"1 2+2 3", "3 5"},
		{"3+1 2 3", "4 5 6"},
		{"1 2 3+4", "5 6 7"},
		{"1 2 3+4f", "5 6 7f"},
		{"3+0x02", "5"},
		{"(1 2;3 4)+5", "(6 7;8 9)"},
		{"2+(3f;,4 5)", "(5f;,6 7)"},
		{"1 2 3<4f", "1 1 1"},
		{"1 2>3 4f", "0 0"},
		{"`a`b`c=`c`b`a", "0 1 0"},
		{"`a`b!2 3", "`a`b!2 3"},
		{"`a!4", "(,`a)!,4"},
		{"`a!1 2 3", "(,`a)!,1 2 3"},
		{"`a`b`c!!3", "`a`b`c!0 1 2"},
		{"`a`b`c!(1;2;3)", "`a`b`c!1 2 3"},
		{"`a`b`c!(1;`a;3)", "`a`b`c!(1;`a;3)"},
		{"(`a`b!1 2),`c!3", "`a`b`c!1 2 3"},
		{"(1;2+3;4)", "1 5 4"},
		{"1;2", "2"},
		{"1 2 3[0 2]", "1 3"},
		{"(1;(`a;3);4)[1;0]", "`a"},
		{"(`a`b!1 2)[`b]", "2"},
		{"`a`b`c^`b", "`a`c"},
		{"`a`b`c^,`b", "`a`c"},
		{"2#!3", "0 1"},
		{"2#!30", "0 1"},
		{"5#!3", "0 1 2 0 1"},
		{"-2#!3", "1 2"},
		{"-5#!3", "1 2 0 1 2"},
		{"0#!3", "!0"},
		{"2#!0", "0N 0N"},
		{"-3#!0", "0N 0N 0N"},
		{"(!0)#1 2 3", "1"},
		{"2 3#!5", "(0 1 2;3 4 0)"},
		{"3 -2#!5", "(3 4;0 1;2 3)"},
		{"3 2 3#!5", "((0 1 2;3 4 0);(1 2 3;4 0 1);(2 3 4;0 1 2))"},
		{"2 3#(`a;2)", "((`a;2;`a);(2;`a;2))"},
		{"2 3#(`a;(1;`b))", "((`a;(1;`b);`a);((1;`b);`a;(1;`b)))"},
		{"-3#(1;(`a`b;4.5))", "((`a`b;4.5);1;(`a`b;4.5))"},
		{"`a#`a`c!1 2", "(,`a)!,1"},
		// {"`a`b`c#`a`c!1 2", "`a`b`c!1 0N 2"}, index error instead of na(k7)
		{"`b_`a`b`c!1 2 3", "`a`c!1 3"},
		{"`c`b_`a`b`c!1 2 3", "(,`a)!,1"},
		{"3 5_!10", "(3 4;5 6 7 8 9)"},
		{"0 3 5_!10", "(0 1 2;3 4;5 6 7 8 9)"},
		{"0 3 3 5_!6", "(0 1 2;!0;3 4;,5)"},
		{"0 3_!3", "(0 1 2;!0)"},
		{"3 3_!3", "(!0;!0)"},
		{"2 4_(1;`a;3f;\"x\";8)", `((3f;"x");,8)`},
		{"2 4_(1;2;`a;`b;8)", "(`a`b;,8)"},
		{"`i$5", "5"},
		{"`c$5", "0x05"},
		{"`f$5 6", "5 6f"},
		{"`z$5", "5i0"},
		{"`i$`xyz", "1"},
		{"`i$`", "0"},
		{"`n$\"alpha\"", "`alpha"},
		{"`n$\"x\"", "`x"},
		{"`$\"xyz\"", "`xyz"},
		{`#"\\"`, "1"},
		{`"\\"`, `"\\"`},
		{"2 3 1 4 0?3", "1"},
		{"2 3 1 4 0?3 1 0", "1 2 4"},
		{"4 2 1f?5 6f", "3 3"}, // TODO: or 0N?
		{"2+", "2+"},
		{"2~3", "0"},
		{"(*3.0)~3f", "1"},
		{"+[;3]", "+[;3]"},
		{"(2+).,3", "5"},
		{"(1;(2;3f)).(1 1)", "3f"},
		{"(1;(2;3f)).,1", "(2;3f)"},
		{"(1;(2;3f))@1", "(2;3f)"},
		{"{1+x}", "{1+x}"},
		{"1 3 4 8 10 bin 1 -3 3 5 9", "0 -1 1 3 4"},
		{"1 3 4 8 10f bin 3.5", "2"},
		{"`abc`def`ghi bin `d`e`j", "1 2 3"},
		{"x:1;y:2;x+y", "3"},
		{"x:1;x", "1"},
		{"x:3;x:4", "4"},
		{"x:3", "3"},
		{"::x:3", "3"},
		{"x:3;-x", "-3"},
		{"g:{x+y};g[3;4]", "7"},
		{"x:1", "1"},
		{"{2*x}3", "6"},
		{"+", "+"},
		{"'", "'"},
		{"~`a``b`c", "0 1 0 0"},
		{"~(;)", "1 1"},
		{"@(;)[0]", "`"},
		{"#(;)[0]", "1"},
		{"~(;)[0]", "1"},
		{"~(+)", "0"},
		{"~!0", "!0"},
		{"`p\"1+2 3\"", "(+;1;2 3)"},
		{"+/(0x01;2;3f)", "6f"},
		{"+/1 2 3", "6"},
		{"*/4 5 6", "120"},
		{"+/,3", "3"},
		{"+/3", "3"},
		{"+/1 2 3f", "6f"},
		{"+/(1;(2f;3))", "(3f;4)"},
		{"4+/1 2 3f", "10f"},
		{"3+/2 3 4", "12"},
		{"3+/(0x02;3f;4)", "12f"},
		{"(,3)+/2 3 4", ",12"},
		{"3 4 5+/2 3", "8 9 10"},
		{"1 2 3+/4 5", "10 11 12"},
		{`-\4 8 9`, "4 -4 -13"},
		{`-\(2 3;5 8;2 1)`, `(2 3;-3 -5;-5 -6)`},
		{`3-\5 9 3`, "-2 -11 -14"},
		{`(,3)-\5 9 3`, "(,-2;,-11;,-14)"},
		{`2 3-\5 9 3`, "(-3 -2;-12 -11;-15 -14)"},
		{`2 3-\(5f;9;0x03)`, "(-3 -2f;-12 -11f;-15 -14f)"},
		{`3 2 1-\,3`, ",0 -1 -2"},
		{`2 3-\5 9 3`, "(-3 -2;-12 -11;-15 -14)"},
		{`+\4 8 9`, "4 12 21"},
		{`+\(4f;0x08;9)`, "4 12 21f"},
		{"+/`a`b`c!1 2 3", "6"},
		{"+\\`a`b`c!1 2 3", "`a`b`c!1 3 6"},
		{"%:'4 5", "0.25 0.2"},
		{"%:'`a`b!4 8", "`a`b!0.25 0.125"},
		{"12%'3 4", "4 3f"},
		{"12 15%'3", "4 5f"},
		{"1 2 3{2*x+y}'4 5 6", "10 14 18"},
		{"-':8 2 5", "8 -6 3"},
		{"-':`a`b`c!3 8 1", "`a`b`c!3 5 -7"},
		{"1 2 3+/:7 8", "(8 9 10;9 10 11)"},
		{`1 2 3+\:7 8`, "(8 9;9 10;10 11)"},
		{`{4>x}{x+1}/1`, "4"},
		{`{4>x}{x+1}\1`, "1 2 3 4"},
		{`{4>#x}{x,"k"}/"o"`, `"okkk"`},
		{"3(2,)/1", "2 2 2 1"},
		{`3(2,)\1`, "(1;2 1;2 2 1;2 2 2 1)"},
		{`3(-:)\1`, "1 -1 1 -1"},
		{"(1_)\\!3", "(0 1 2;1 2;,2;!0)"},
		{"(1_)/!3", "!0"},
		{`|+\!3`, `3 1 0`},
		{`"-"\:"ab-cd--ef-gh-"`, `("ab";"cd";"";"ef";"gh";"")`},
		{"`\\:\"ab\ncd\n\ne\n\n\"", `("ab";"cd";"";,"e";"")`}, // 1 trailing nl removed
		{`"x"/:(,"a";,"b";,"c")`, `"axbxc"`},
		{"`/:(\"aa\";\"bb\")", `"aa\nbb\n"`},
		{`2\:234`, "1 1 1 0 1 0 1 0"},
		{`8\:234`, "3 5 2"},
		{`16\:234`, "14 10"},
		{`16/:14 10`, "234"},
		{"$[0;2;3]", "3"},
		{"$[0;2;1;4]", "4"},
		{"$[0;1;0;2;3;4;5]", "4"},
		{"$[0;1;2]", "2"},
		{"$[0x00;1;2]", "2"},
		{"$[`;1;2]", "2"},
		{"$[`a;1;2]", "1"},
		{"$[\"\";1;2]", "2"},
		{"$[0 0;1;2]", "1"},
		{"x:`a`b!(1;2 3);x[`b]", "2 3"},
		{"x:`a`b!(1;2 3);x[`b;0]", "2"},
		{"x:`a`b!(1;2 3);x.b", "2 3"},
		{"x:`a`b!(1;2 3);x.b 1", "3"},
		{"a:`a`b!(1;`b`c!(2;3f));a.b.c", "3f"},
		{"()", "()"},
		{`"."\:""`, "()"},
		{"*()", `""`},
		{"*!0", "0N"},
		{"*!0f", "0n"},
		{`*""`, `" "`},
		{"*0#`", "`"},
		{"(a;b):3 4;a+b", "7"},
		{"(a;b):1;a+b", "2"},
		{"x:8;x-:", "-8"},
		{"x:8;x-:;x", "8"},
		{"x:`a`b!1 2;x[`b]:3;x", "`a`b!1 3"},
		{"y:`b;x:`a`b!1 2;x[y]:3;x", "`a`b!1 3"},
		{"x:`a`b!1 2;x.b:3;x", "`a`b!1 3"},
		{"x:!3;x[0]:4;x", "4 1 2"},
		{"x:!3;x[0 2]:4;x", "4 1 4"},
		{"x:!3;x[0 2]:-1 4;x", "-1 1 4"},
		{"x:(1;2f;`a);x[0]:3;x", "(3;2f;`a)"},
		{"x:(1;2f;`a);x[0 2]:-(1 2);x", "(-1;2f;-2)"},
		{"x:!4;@[`x;1 2;*;10];x", "0 10 20 3"},
		{"x:(1;2f;0x03);@[`x;0 1;*;10];x", "(10;20f;0x03)"},
		{"x:(1;2f;`a);x[0 2]:3 4f;x", "3 2 4f"},
		{"x:2f;x*:3;x", "6f"},
		{"(x;y):2 3f;(x;y)*:3;x+y", "15f"},
		{"x:!5;x[2 3]*:2;x", "0 1 4 6 4"},
		{"x:(1;2 3 4);x[1]:4 5;x", "(1;4 5)"},
		{"x:1 2;@[`x;0;3];x", "3 2"},
		{"x:1 2;@[`x;0;3]", "`x"},
		{"x:`a`b!3 4;x[!x]", "3 4"},
		{"@[1 2;0;3]", "3 2"},
		{"(3 3#!9)[1;0 2]", "3 5"},
		{"x:3 3#!9;x[1;2]:0;x", "(0 1 2;3 4 0;6 7 8)"},
		{"@[(1;2 3 4);1;5 5]", "(1;5 5)"},
		{"@[`a`b!1 2;`a;3]", "`a`b!3 2"},
		{"x:`a`b!(1;`c`d!(2;3f));y:x;x.b.d:4;x.b.d+y.b.d", "7f"},
		{"x:3 3#!9;x[1;0 2]+:1 2;x", "(0 1 2;4 4 7;6 7 8)"},
		{"x:3 3#!9;.[`x;1 2;0];x", "(0 1 2;3 4 0;6 7 8)"},
		{"x:3 3#!9;.[`x;1 2;+;2];x", "(0 1 2;3 4 7;6 7 8)"},
		{"x:`a`b!(1;2 3);x.b[0]:4;x", "`a`b!(1;4 3)"},
		{"x:`a`b!(1;`c`d!1 2);x.b.c:4;x", "`a`b!(1;`c`d!4 2)"},
		{"x:`a`b!(1;`c`d!1 2);x.b[`c]:4;x", "`a`b!(1;`c`d!4 2)"},
		// {"(-).(1 2)", "3"},
	}
	for _, occ := range []bool{false, true} {
		for _, tc := range testCases {
			ini()
			fmt.Printf("%s → %s\n", tc.x, tc.r)
			x := prs(K([]byte(tc.x)))
			if occ {
				inc(x)
			}
			y := kst(evl(x))
			r := string(G(y).([]c))
			if r != tc.r {
				t.Fatalf("expected %s got %s\n", tc.r, r)
			}
			if occ {
				dec(x)
			}
			dec(y)
			clear()
			check(t)
		}
	}
}

func TestParse(t *testing.T) {
	ini()
	//t.Skip()
	testCases := []struct {
		x, r, t s
	}{
		{"", "", "`"},
		{"`a", ",`a", "`N"},
		{"`a`b", ",`a`b", "`."},
		{"`a / b", ",`a", "`N"},
		{".x", "`.x", "`n"},
		{"{1+x}", "{1+x}", "`1"},
		{"{x+y}", "{x+y}", "`2"},
		{"{1+(2;`a;y)}", "{1+(2;`a;y)}", "`2"},
		{"/alpha\n`a", "(`;;,`a)", "`."},
		{"/alpha\n`a /beta\n/gamma", "(`;;,`a;)", "`."},
		{"0x01", "0x01", "`c"},
		{"0xF", "0x0f", "`c"},
		{"0x", `""`, "`C"},
		{"0x1234", "0x1234", "`C"},
		{`"a"`, `"a"`, "`c"},
		{`"a\t\n\r\"xyz"`, `"a\t\n\r\"xyz"`, "`C"},
		{"10", "10", "`i"},
		{"10 20", "10 20", "`I"},
		{"-3", "-3", "`i"},
		{".1", "0.1", "`f"},
		{"2.", "2f", "`f"},
		{"-1.23e-005", "-1.23e-05", "`f"},
		{"1 2f", "1 2f", "`F"},
		{"1.23 3", "1.23 3", "`F"},
		{"1 2 3. ", "1 2 3f", "`F"},
		{"2i-3", "2i-3", "`z"},
		{"-2.0e+012i-3.6", "-2e+12i-3.6", "`z"},
		{"1 2 3. 2i-3", "1i0 2i0 3i0 2i-3", "`Z"},
		{"+", "+", "`2"},
		{"+:", "+:", "`1"},
		{"0:", "0:", "`2"},
		{"0:1", "(0::;1)", "`."},
		{"x-:", "(-:;`x)", "`."},
		{"'", "'", "`2"},
		{"1+2;", "(`;(+;1;2);)", "`."},
		{"(/)", "/", "`2"},
		{"2/3", "((/;2);3)", "`."},
		{"(/:)", "/:", "`1"},
		{"in", "in", "`2"},
		{"within", "within", "`2"},
		{"bin", "bin", "`2"},
		{"like", "like", "`2"},
		{"1 2@0", "(1 2;0)", "`."},
		{"+/", "(/;+)", "`."},
		{"+/1 2 3", "((/;+);1 2 3)", "`."},
		{"1+/3 4 5", "((/;+);1;3 4 5)", "`."},
		{"*1 2 3", "(*:;1 2 3)", "`."},
		{"1+(2;`a;3.5)", "(+;1;(;2;,`a;3.5))", "`."},
		{"1;2\n3", "(`;1;2;3)", "`."},
		{"1 2 3[0 2]", "(1 2 3;0 2)", "`."},
		{"`a`b`c[2]", "(,`a`b`c;2)", "`."},
		{"(1;(2;3);4)[1;1]", "((;1;(;2;3);4);1;1)", "`."},
		{"(1;(`a;3);4)[1;0]", "((;1;(;,`a;3);4);1;0)", "`."},
		{"x[1;2;3]", "(`x;1;2;3)", "`."},
		{"(`a`b!1 2)", "(!;,`a`b;1 2)", "`."},
		{"(`a`b!1 2)[`b]", "((!;,`a`b;1 2);,`b)", "`."},
		{"2+", "(+;2)", "`."},
		{"{x+y}[2;3]", "({x+y};2;3)", "`."},
		{"g:{x+y};g[3;4]", "(`;(:;`g;{x+y});(`g;3;4))", "`."},
		{"-':8 2 5", "((':;-);8 2 5)", "`."},
		{`|+\!3`, `(|:;((\;+);(!:;3)))`, "`."},
		{"{-x}+3", `(+;{-x};3)`, "`."},
		{"$[0;1;0;2;3]", "($;0;1;0;2;3)", "`."},
		// {"(-).(1 2)", "(.;-;1 2)", "`."},
	}
	for i, occ := range []bool{true, false} {
		for j, tc := range testCases {
			for _, at := range []bool{false, true} {
				x := K([]byte(tc.x))
				if occ {
					inc(x)
				}
				y := prs(x)
				exp := tc.r
				if at {
					fmt.Printf("@`p(%q) ~ %v\n", tc.x, tc.t)
					y = kst(tip(y))
					exp = tc.t
				} else {
					fmt.Printf("`p(%q) ~ %v\n", tc.x, tc.r)
					y = kst(y)
				}
				r := string(G(y).([]c))
				if !reflect.DeepEqual(r, exp) {
					t.Fatalf("[%d/%d]: expected: %v got %v\n", j, i, exp, r)
				}
				dec(y)
				if occ {
					dec(x)
				}
				check(t)
			}
		}
	}
}
func TestNumMonad(t *testing.T) {
	// t.Skip()
	ini()
	xv := []interface{}{c(3), []c{3, 5}, -5, iv{3, -9}, 3.2, []f{-3.5, 2.9, 0}, 2 - 4i, []z{4 - 2i, 3 + 4i}}
	testCases := []struct {
		f func(k) k
		s s
		r []interface{}
	}{
		{neg, "-", l{c(253), []c{253, 251}, 5, iv{-3, 9}, -3.2, []f{3.5, -2.9, -0}, -2 + 4i, []z{-4 + 2i, -3 - 4i}}},
		{fst, "*", l{c(3), c(3), -5, 3, 3.2, -3.5, 2 - 4i, 4 - 2i}},
		{rev, "|", l{c(3), []c{5, 3}, -5, iv{-9, 3}, 3.2, []f{0, 2.9, -3.5}, 2 - 4i, []z{3 + 4i, 4 - 2i}}},
		{not, "~", l{0, iv{0, 0}, 0, iv{0, 0}, 0, iv{0, 0, 1}, 0, iv{0, 0}}},
		{enl, ",", l{[]c{3}, l{[]c{3, 5}}, iv{-5}, l{iv{3, -9}}, []f{3.2}, l{[]f{-3.5, 2.9, 0}}, []z{2 - 4i}, l{[]z{4 - 2i, 3 + 4i}}}},
		{cnt, "#", l{1, 2, 1, 2, 1, 3, 1, 2}},
		{str, "$", l{c(3), []c{3, 5}, []c("-5"), l{[]c("3"), []c("-9")}, []c("3.2"), l{[]c("-3.5"), []c("2.9"), []c("0")}, []c("2i-4"), l{[]c("4i-2"), []c("3i4")}}},
		{kst, "`k@", l{[]c("0x03"), []c("0x0305"), []c("-5"), []c("3 -9"), []c("3.2"), []c("-3.5 2.9 0"), []c("2i-4"), []c("4i-2 3i4")}},
		{tip, "@", l{"c", "C", "i", "I", "f", "F", "z", "Z"}},
		{evl, ".", xv},
	}
	for _, occ := range []bool{true, false} {
		for j, tc := range testCases {
			for i := range xv {
				// fmt.Println("TC", xv[i])
				x := K(xv[i])
				if x == 0 {
					t.Fatalf("cannot import go type %T", xv[i])
				}
				if occ {
					inc(x)
				}
				y := tc.f(x)
				if occ {
					dec(x)
				}
				r := G(y)
				fmt.Printf("%s(%v) = %v\n", tc.s, xv[i], r)
				if !reflect.DeepEqual(r, tc.r[i]) {
					t.Fatalf("[%d/%d]: expected: %v got %v (@%d)\n", j, i, tc.r[i], r, y)
				}
				dec(y)
				check(t)
			}
		}
	}
}
func TestMonad(t *testing.T) {
	// t.Skip()
	ini()
	testCases := []struct {
		f    func(k) k
		s    s
		x, r interface{}
	}{
		{flp, "+", l{iv{1, 2}, iv{3, 4}}, l{iv{1, 3}, iv{2, 4}}},
		{flp, "+", l{iv{1, 2, 3}, []f{3, 4, 5}}, l{l{1, 3.0}, l{2, 4.0}, l{3, 5.0}}},
		{flp, "+", l{l{1, 2.2, l{3, 4}}, l{"a", []c{'x'}, 5 + 2i}}, l{l{1, "a"}, l{2.2, []c{'x'}}, l{l{3, 4}, 5 + 2i}}},
		{flp, "+", l{l{1, 2.2}, l{3, 4.4}}, l{iv{1, 3}, []f{2.2, 4.4}}},
		{til, "!", 3, iv{0, 1, 2}},
		{til, "!", -1, l{iv{1}}},
		{til, "!", -3, l{iv{1, 0, 0}, iv{0, 1, 0}, iv{0, 0, 1}}},
		{til, "!", d{sv{"a", "b"}, iv{1, 2}}, sv{"a", "b"}},
		// TODO !a(odometer)
		{fst, "*", l{3, 4, 5}, 3},
		{fst, "*", "alpha", "alpha"},
		{fst, "*", l{"alpha"}, "alpha"},
		{fst, "*", d{l{"x", "y"}, l{iv{5, 3}, 4}}, iv{5, 3}},
		{fst, "*", d{sv{"x", "y"}, iv{7, 2}}, 7},
		{inv, "%", 4, 0.25},
		{inv, "%", []f{4.0, 8.0}, []f{0.25, 0.125}},
		{str, "$", "", []c("")},
		{str, "$", "a", []c("a")},
		{str, "$", sv{"", "a", "bb", "a\t\r\nb"}, l{[]c(""), []c("a"), []c("bb"), []c("a\t\r\nb")}},
		{str, "$", l{1, c(3), l{4, 5.0}}, l{[]c("1"), c(3), l{[]c("4"), []c("5")}}},
		{str, "$", d{sv{"x", "y"}, iv{1, 2}}, d{sv{"x", "y"}, l{[]c("1"), []c("2")}}},
		{kst, "`k", iv{1, 2, 3}, []c("1 2 3")},
		{kst, "`k", l{1, 2, l{4, 5}}, []c("(1;2;(4;5))")},
		{kst, "`k", d{l{5, 5.5}, iv{1, 2}}, []c("(5;5.5)!1 2")},
		{rev, "|", l{}, l{}},
		{rev, "|", l{iv{3}}, l{iv{3}}},
		{rev, "|", l{1, 2}, l{2, 1}},
		{rev, "|", l{1, l{3, 4}}, l{l{3, 4}, 1}},
		{rev, "|", d{iv{1, 2}, iv{3, 4}}, d{iv{2, 1}, iv{4, 3}}},
		{rev, "|", d{sv{"alpha", "beta"}, l{3, iv{3, 5}}}, d{sv{"beta", "alpha"}, l{iv{3, 5}, 3}}},
		{wer, "&", iv{0, 0, 1, 1, 0, 1}, iv{2, 3, 5}},
		{wer, "&", iv{}, iv{}},
		{wer, "&", iv{2}, iv{0, 0}},
		{wer, "&", iv{1, 2, 3}, iv{0, 1, 1, 2, 2, 2}},
		{asc, "<", iv{1, 2, 3, 4}, iv{0, 1, 2, 3}},
		{asc, "<", iv{1, 4, 3, 2}, iv{0, 3, 2, 1}},
		{asc, "<", iv{4, 2, 3, 4}, iv{1, 2, 0, 3}},
		{asc, "<", []f{4, 1, 2}, iv{1, 2, 0}},
		{asc, "<", []c{6, 4, 2, 1}, iv{3, 2, 1, 0}},
		{asc, "<", []z{4, 1, 2}, iv{1, 2, 0}},
		{asc, "<", []z{0, 1 + 1i, 1, 2}, iv{0, 2, 1, 3}},
		{asc, "<", sv{"b", "ab", "a", "aa"}, iv{2, 3, 1, 0}},
		{dsc, ">", iv{1, 4, 3, 2}, iv{1, 2, 3, 0}},
		{grp, "=", []c{'c', 'b', 'a', 'c', 'a', 'b', 'c'}, d{[]c{'c', 'b', 'a'}, l{iv{0, 3, 6}, iv{1, 5}, iv{2, 4}}}},
		{grp, "=", iv{1, 2, 3}, d{iv{1, 2, 3}, l{iv{0}, iv{1}, iv{2}}}},
		{grp, "=", iv{3, 3, 1, 3, 2, 1}, d{iv{3, 1, 2}, l{iv{0, 1, 3}, iv{2, 5}, iv{4}}}}, // TODO sort keys(k7)
		{grp, "=", []f{5.5, 1, 3, 3, 2}, d{[]f{5.5, 1, 3, 2}, l{iv{0}, iv{1}, iv{2, 3}, iv{4}}}},
		{grp, "=", []z{3, 3 + 1i, 3, 3 + 1i, 3 + 1i}, d{[]z{3, 3 + 1i}, l{iv{0, 2}, iv{1, 3, 4}}}},
		{grp, "=", sv{"alpha", "beta", "alpha", "gamma", "alpha", "beta"}, d{sv{"alpha", "beta", "gamma"}, l{iv{0, 2, 4}, iv{1, 5}, iv{3}}}},
		{enl, ",", "alpha", sv{"alpha"}},
		{enl, ",", l{1, 2, l{3, 4.5}}, l{l{1, 2, l{3, 4.5}}}},
		{enl, ",", d{iv{3, 4}, sv{"x", "y"}}, l{d{iv{3, 4}, sv{"x", "y"}}}},
		{srt, "^", iv{3, 1, 2, 3, 5}, iv{1, 2, 3, 3, 5}},
		{cnt, "#", "alpha", 1},
		{cnt, "#", l{}, 0},
		{cnt, "#", l{1, 2, l{3, 4}}, 3},
		{cnt, "#", d{iv{3, 4}, sv{"x", "y"}}, 2},
		{tip, "@", l{}, "."},
		{tip, "@", d{iv{1, 2}, iv{3, 4}}, "a"},
		{evl, ".", l{uint16(2), l{uint16(2), 3}}, 3},
		{evl, ".", l{uint16(2), l{uint16(6), iv{3, 4}}}, iv{-4, -3}},
		{evl, ".", l{uint16(2), iv{3, 4}}, iv{-3, -4}},
		{unq, "?", []c{1, 2, 43, 2}, []c{1, 2, 43}},
		{unq, "?", iv{1, 2, 3, 2}, iv{1, 2, 3}},
		{unq, "?", []f{5, 0, 0, 0, 8, 0, 0, 0, 5, 0, 0, 5}, []f{5, 0, 8}},
		{unq, "?", []z{0, 4i, 5i, 4i, 0, 3}, []z{0, 4i, 5i, 3}},
		{unq, "?", l{1, 2, 3, 1}, l{1, 2, 3}},
		{unq, "?", l{1i, l{2, sv{"a"}}, l{3, "b"}, l{2, sv{"a"}}, 1i}, l{1i, l{2, sv{"a"}}, l{3, "b"}}},
	}
	for _, occ := range []bool{true, false} {
		for j, tc := range testCases {
			//fmt.Println("TC", i, j, tc.s, tc.x, "occ", occ)
			x := K(tc.x)
			_ = Stats().UsedBlocks()
			if x == 0 {
				t.Fatalf("cannot import go type %T", tc.x)
			}
			if occ {
				inc(x)
			}
			y := tc.f(x)
			fpck("1")
			if occ {
				dec(x)
			}
			r := G(y)
			fmt.Printf("%s[%v] = %v\n", tc.s, tc.x, r)
			if !reflect.DeepEqual(r, tc.r) {
				t.Fatalf("monad[%d]: expected: %v got %v (@%d)\n", j, tc.r, r, y)
			}
			dec(y)
			check(t)
		}
	}
}
func TestDyad(t *testing.T) {
	// t.Skip()
	testCases := []struct {
		f       func(k, k) k
		s       s
		x, y, r interface{}
	}{
		{atx, "@", iv{2, 1, 3, 5}, iv{2, 0, 1}, iv{3, 2, 1}},
		{ept, "^", iv{1, 5, 3, 3, 2}, iv{6, 3}, iv{1, 5, 2}},
		{drp, "_", 2, iv{1, 2, 3, 4}, iv{3, 4}},
		{drp, "_", -2, iv{1, 2, 3, 4}, iv{1, 2}},
		{drp, "_", 2, l{1, 2.0, "a"}, sv{"a"}},
		{drp, "_", -2, l{1, 2.0, "a"}, iv{1}},
		{fnd, "?", iv{1, 2}, iv{2}, iv{1}},
		{fnd, "?", []c{5, 4, 3, 2}, c(2), 3},
	}
	for _, occ := range []bool{true, false} {
		for j, tc := range testCases {
			// fmt.Println("TC", i, j, tc.s, tc.x, "occ", occ)
			x := K(tc.x)
			y := K(tc.y)
			_ = Stats().UsedBlocks()
			if x == 0 || y == 0 {
				t.Fatalf("cannot import go type %T", tc.x)
			}
			if occ {
				inc(x)
				inc(y)
			}
			z := tc.f(x, y)
			fpck("1")
			if occ {
				dec(x)
				dec(y)
			}
			r := G(z)
			fmt.Printf("%s[%v, %v] = %v\n", tc.s, tc.x, tc.y, r)
			if !reflect.DeepEqual(r, tc.r) {
				t.Fatalf("dyad[%d]: expected: %v got %v\n", j, tc.r, r)
			}
			dec(z)
			check(t)
		}
	}
}
func TestKst(t *testing.T) {
	// t.Skip()
	ini()
	testCases := []struct {
		x interface{}
		s s
	}{
		{[]c{}, `""`},
		{c('x'), `"x"`},
		{[]c{'x'}, `,"x"`},
		{c(28), "0x1c"},
		{[]c{28}, ",0x1c"},
		{[]c{0x1b, 0x5b, 0x5c}, "0x1b5b5c"},
		{[]c("alpha"), `"alpha"`},
		{[]c("a\nb\rc\t\"d"), `"a\nb\rc\t\"d"`},
		{"alpha", "`alpha"},
		{sv{}, "0#`"},
		{sv{"alpha"}, ",`alpha"},
		{sv{"a", "b", "c"}, "`a`b`c"},
		{1, "1"},
		{iv{}, "!0"},
		{iv{1}, ",1"},
		{[]f{}, "!0.0"},
		{[]f{1.2, -3.5, 4}, "1.2 -3.5 4"},
		{[]f{3, 5, 4}, "3 5 4f"},
		{13.0, "13f"},
		{1 + 2i, "1i2"},
		{[]z{}, "!0i0"},
		{[]z{2i}, ",0i2"},
		{[]z{2i, 3.5 + 7i}, "0i2 3.5i7"},
		{l{1, 2, l{4, 5}}, "(1;2;(4;5))"},
		{l{iv{1, 2}}, ",1 2"},
		{d{l{5, 5.5}, iv{1, 2}}, "(5;5.5)!1 2"},
		{nil, ``},
		{uint16(0), `::`},
		{uint16(1), `+:`},
		{uint16(50), `:`},
		{uint16(69), `.`},
		{uint16(70), `0:`},
		{uint16(74), `4:`},
	}
	for _, tc := range testCases {
		fmt.Printf("%v ?= %q\n", tc.x, tc.s)
		x := K(tc.x)
		y := kst(x)
		r := G(y).([]c)
		if reflect.DeepEqual(r, []byte(tc.s)) == false {
			t.Fatalf("expected: %s got %s (%q)\n", tc.s, string(r), string(r))
		}
		dec(y)
	}
}
func TestTo(t *testing.T) {
	// t.Skip()
	ini()
	testCases := []struct {
		x, r interface{}
		t    k
	}{
		{c(1), 1, I},
		{c(1), 1.0, F},
		{[]c{1, 2, 3, 4}, iv{1, 2, 3, 4}, I},
		{1, c(1), C},
		{5, 5.0, F},
		{1, 1 + 0i, Z},
		{iv{1, 2, 3, 4}, []z{1, 2, 3, 4}, Z},
		{2.3, c(2), C},
		{2.3, 2, I},
		{2.3, 2.3 + 0i, Z},
		{-2.3 + 4.5i, -2.3, F},
		{2.3 - 4.5i, c(2), C},
	}
	for _, occ := range []bool{true, false} {
		for _, tc := range testCases {
			x := K(tc.x)
			if occ {
				inc(x)
			}
			y := to(x, tc.t)
			if occ {
				dec(x)
			}
			r := G(y)
			fmt.Printf("to(%v,%d) = %v\n", tc.x, tc.t, tc.r)
			if !reflect.DeepEqual(r, tc.r) {
				t.Fatalf("expected: %v got %v\n", tc.r, r)
			}
			dec(y)
			check(t)
		}
	}
}
func TestStr(t *testing.T) {
	// t.Skip()
	ini()
	for _, x := range []s{"a", "b", "aa", "bb", "alpha", "betagamm"} {
		n := len(x)
		if n > 8 {
			n = 8
		}
		if r := G(K(x)); r != x[:n] {
			t.Fatalf("expected %s got %s\n", x, r)
		}
	}
	if u := sym(8 + K("abcdefgh")<<2); u != 0x6162636465666768 {
		t.Fatalf("%x\n", u)
	}
}
func check(t *testing.T) {
	// Number of used blocks after an expression should be:
	// 1(block 0) + 3(built-in dict,k,v) + 2(k-tree k,v) + number of variables (`.f`.n`.l)
	// vars := m.k[m.k[kkey]] & atom
	if u := Stats().UsedBlocks(); u != 9 {
		xxd()
		t.Fatalf("leak: %d", u)
	}
	fpck("")
}
func pfl() {
	for i := 4; i < 32; i++ {
		println(i, strconv.FormatUint(uint64(m.k[i]), 16), strconv.FormatUint(uint64(m.k[i]<<2), 16))
	}
}
func xdec(x k) {
	if m.k[x]>>28 != 0 {
		dec(x)
	}
}
func xxd() { // memory dump
	h := k(0)
	for i := k(0); i < k(len(m.k)); i += 4 {
		a, b, c, d := m.k[i+0], m.k[i+1], m.k[i+2], m.k[i+3]
		if a == 0 && b == 0 && c == 0 && d == 0 {
			continue
		}
		fmt.Printf("0x%04x %08x %08x %08x %08x", i, a, b, c, d)
		if i == h {
			tp := m.k[i] >> 28
			if tp == 0 {
				fmt.Printf("  %d", m.k[i])
				h += 1 << (m.k[i] - 2)
				nf := m.k[i+1]
				if nf > 0 && nf < 64 {
					fmt.Printf(" illegal fp")
				} else if nf > 0 && m.k[nf]>>28 != 0 {
					fmt.Printf(" fp is not free")
				}
			} else {
				atoms := "?cifzsld01234"
				vects := "?CIFZSLD01234"
				tp, n := typ(i)
				bt := bk(tp, n)
				if n == atom {
					fmt.Printf(" %c%d +%d", atoms[tp], bt, b)
				} else {
					fmt.Printf(" %c%d #%d +%d", vects[tp], bt, n, b)
				}
				h += 1 << (bt - 2)
			}
		}
		fmt.Println()
	}
}
func fpck(s s) { // check free pointers
	for i := 4; i < 32; i++ {
		nf := m.k[i]
		if nf > 0 && (nf < 64 || m.k[nf]>>28 != 0) {
			xxd()
			panic("fpck " + s + " bad pointer in free-list: @" + strconv.Itoa(int(i)))
		}
	}
	h := k(0)
	for i := k(0); i < k(len(m.k)); i += 4 {
		if i == h {
			tp := m.k[i] >> 28
			if tp == 0 {
				h += 1 << (m.k[i] - 2)
				nf := m.k[i+1]
				if nf > 0 && (nf < 64 || m.k[nf]>>28 != 0) {
					xxd()
					panic("fpck " + s + " illegal free-pointer")
				}
			} else {
				tp, n := typ(i)
				bt := bk(tp, n)
				h += 1 << (bt - 2)
			}
		}
	}
}
func pr(x k, a ...interface{}) {
	fmt.Printf(":%x ", x)
	r := kst(inc(x))
	_, n := typ(r)
	s := s(m.c[8+r<<2 : 8+n+r<<2])
	dec(r)
	fmt.Println(a, s)
}

type Bucket struct {
	Type       uint32
	Used, Free uint32 // num blocks
	Net        uint32
}
type MemStats map[uint32]Bucket

func (b Bucket) Overhead() uint32 {
	return b.Used*uint32(1<<b.Type) - b.Net
}
func (s MemStats) UsedBlocks() (t uint32) {
	for _, b := range s {
		t += b.Used
	}
	return t
}
func Stats() MemStats {
	st := make(MemStats)
	a := uint32(0)
	o := uint32(0)
	for a < 1<<(m.k[2]-2) {
		tp := m.k[a] >> 28
		if tp == 0 {
			t := m.k[a]
			if t < 4 || t > 31 {
				xxd()
				fmt.Printf("free block at %x with bt %d\n", a, t)
				panic("size")
			}
			b := st[t]
			b.Type = t
			b.Free++
			st[t] = b
			o = 1 << (t - 2)
		} else {
			tt, n := typ(a)
			t := bk(tt, n)
			if t < 4 || t > 31 {
				println(a, t)
				panic("size")
			}
			b := st[t]
			b.Type = t
			b.Used++
			if n == atom {
				n = 1
			}
			b.Net += n * lns[tp]
			st[t] = b
			o = 1 << (t - 2)
		}
		a += o
	}
	return st
}

// type conversions between go and k (used by k_test.go)

func K(x interface{}) k { // convert go value to k type, returns 0 on error
	if x == nil {
		return mk(N, atom)
	}
	kstr := func(dst k, s string) {
		mys(dst, btou([]c(s)))
	}
	var r k
	switch a := x.(type) {
	case bool:
		r = mkc(0)
		if a {
			m.c[8+r<<2] = 1
		}
	case byte:
		r = mkc(a)
	case int:
		r = mki(k(a))
	case uint16: // function index
		if a < 20 {
			r = mk(N+1, atom)
		} else {
			r = mk(N+2, atom)
		}
		m.k[2+r] = k(a)
	case float64:
		r = mk(F, atom)
		m.f[1+r>>1] = a
	case complex128:
		r = mk(Z, atom)
		m.z[1+r>>2] = a
	case string:
		r = mks(a)
	case []bool:
		buf := make([]byte, len(a))
		for i, v := range a {
			if v {
				buf[i] = 1
			}
		}
		return K(buf)
	case []byte:
		r = mk(C, k(len(a)))
		for i, v := range a {
			m.c[8+i+int(r<<2)] = v
		}
	case []int:
		r = mk(I, k(len(a)))
		for i, v := range a {
			m.k[2+i+int(r)] = k(v)
		}
	case []float64:
		r = mk(F, k(len(a)))
		for i, v := range a {
			m.f[1+i+int(r>>1)] = v
		}
	case []complex128:
		r = mk(Z, k(len(a)))
		for i, v := range a {
			m.z[1+i+int(r>>2)] = v
		}
	case []string:
		r = mk(S, k(len(a)))
		for i := range a {
			kstr(8+8*k(i)+r<<2, a[i])
		}
	case []interface{}:
		if len(a) == 1 { // collapse list of atom to single element vector
			rr := K(a[0])
			t, n := typ(rr)
			if n == atom { // TODO: allow ,d?
				r = rr
				m.k[r] = t<<28 | 1
				return r
			} else {
				dec(rr)
			}
		}
		r = mk(L, k(len(a)))
		for i, v := range a {
			u := K(v)
			m.k[2+i+int(r)] = u
		}
	case [2]interface{}:
		key := K(a[0])
		val := K(a[1])
		_, nk := typ(key)
		_, nv := typ(val)
		if nk != nv {
			return 0
		}
		r = mk(D, atom)
		m.k[2+r] = key
		m.k[3+r] = val
	}
	return r
}
func G(x k) interface{} { // convert k value to go type (returns nil on error)
	t, n := typ(x)
	str := func(xp k) s {
		r := mk(C, 0)
		rc := 8 + r<<2
		n := stS(rc, xp)
		dec(r)
		return string(m.c[rc : rc+n])
	}
	if n == atom {
		switch t {
		case C:
			return c(m.c[8+x<<2])
		case I:
			return int(i(m.k[2+x]))
		case F:
			return m.f[1+x>>1]
		case Z:
			return m.z[1+x>>2]
		case S:
			return str(1 + x>>1)
		case D:
			return [2]interface{}{G(m.k[2+x]), G(m.k[3+x])}
		case N:
			return nil
		case N + 1, N + 2:
			return uint16(m.k[2+x])
		}
	} else {
		switch t {
		case C:
			r := make([]byte, n)
			for i := range r {
				r[i] = c(m.c[8+i+int(x<<2)])
			}
			return r
		case I:
			r := make([]int, n)
			for i := range r {
				r[i] = int(int32(m.k[2+i+int(x)]))
			}
			return r
		case F:
			r := make([]f, n)
			for i := range r {
				r[i] = m.f[1+i+int(x>>1)]
			}
			return r
		case Z:
			r := make([]complex128, n)
			for i := range r {
				r[i] = m.z[1+i+int(x>>2)]
			}
			return r
		case S:
			r := make([]string, n)
			for i := range r {
				r[i] = str(1 + k(i) + x>>1)
			}
			return r
		case L:
			r := make([]interface{}, n)
			for i := range r {
				r[i] = G(m.k[2+i+int(x)])
			}
			return r
		}
	}
	return nil
}
