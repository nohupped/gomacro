/*
 * gomacro - A Go interpreter with Lisp-like macros
 *
 * Copyright (C) 2019 Massimiliano Ghilardi
 *
 *     This Source Code Form is subject to the terms of the Mozilla Public
 *     License, v. 2.0. If a copy of the MPL was not distributed with this
 *     file, You can obtain one at http://mozilla.org/MPL/2.0/.
 *
 *
 * gen_op2.go
 *
 *  Created on Jan 28, 2019
 *      Author Massimiliano Ghilardi
 */

package main

import (
	"fmt"
	"io"
	"os"

	arch "github.com/cosmos72/gomacro/jit/amd64"
)

type genOp2 struct {
	opname, opName string
	w              io.Writer
}

func NewGenOp2(w io.Writer, opname string) *genOp2 {
	return &genOp2{
		opname: opname,
		opName: string(opname[0]-'a'+'A') + opname[1:],
		w:      w,
	}
}

func GenOp2() {
	for _, opname := range [...]string{"add", "and", "mov", "or", "sub", "xchg", "xor", "imul"} {
		f, err := os.Create("_gen_" + opname + ".s")
		if err != nil {
			panic(err)
		}
		g := NewGenOp2(f, opname)
		g.generate()
		f.Close()
	}
}

func (g *genOp2) generate() {
	g.fileHeader()

	switch g.opname {
	case "xchg":
	case "imul":
	default:
		g.opConstReg()
		g.opConstMem()
	}
	g.opRegReg()
	g.opMemReg()

	switch g.opname {
	case "imul":
		g.opRegConst8Reg()
		g.opMemConst8Reg()
	default:
		g.opRegMem()
	}
}

func (g *genOp2) fileHeader() {
	fmt.Fprintf(g.w,
		`	.file	"%s.s"
	.text
`, g.opname)
}

func (g *genOp2) funcHeader(funcName string) {
	fmt.Fprintf(g.w,
		`
	.p2align 4,,15
	.globl	%s%s
	.type	%s%s, @function
%s%s:
	.cfi_startproc
`, g.opName, funcName, g.opName, funcName, g.opName, funcName)
}

func (g *genOp2) funcFooter() {
	fmt.Fprint(g.w, `	ret
	.cfi_endproc

`)
}

func (g *genOp2) opConstReg() {
	for _, k := range [...]arch.Kind{arch.Uint8, arch.Uint16, arch.Uint32, arch.Uint64} {
		g.opConstRegKind(k, 8)
		if k.Size() != 1 {
			g.opConstRegKind(k, k.Size()*8)
		}
	}
}

func (g *genOp2) opConstRegKind(k arch.Kind, constbits arch.Size) {
	kbits := k.Size() * 8
	conststr := map[arch.Size]string{8: "$0x33", 16: "$0xaabb", 32: "$0x11223344", 64: "$0x55667788"}[constbits]
	g.funcHeader(fmt.Sprintf("Const%dReg%d", constbits, kbits))
	for r := arch.RLo; r <= arch.RHi; r++ {
		fmt.Fprintf(g.w, "\t%s\t%s,%v\n", g.opname, conststr, arch.MakeReg(r, k))
	}
	g.funcFooter()
}

func (g *genOp2) opRegConst8Reg() {
	g.funcHeader("RegConst8Reg")
	for _, k := range [...]arch.Kind{arch.Uint8, arch.Uint16, arch.Uint32, arch.Uint64} {
		if g.opname == "imul" && k.Size() == 1 {
			continue
		}
		fmt.Fprintf(g.w, "\t// reg%d = reg%d OP const8\n", k.Size()*8, k.Size()*8)
		for src := arch.RLo; src <= arch.RHi; src++ {
			for dst := arch.RLo; dst <= arch.RHi; dst++ {
				for _, cval := range [...]int8{-128, 0, 127} {
					fmt.Fprintf(g.w, "\t%s\t$%#x,%v,%v\n", g.opname, cval, arch.MakeReg(src, k), arch.MakeReg(dst, k))
				}
			}
			fmt.Fprint(g.w, "\tnop\n")
		}
		fmt.Fprint(g.w, "\tnop\n")
	}
	g.funcFooter()
}

func (g *genOp2) opRegReg() {
	g.funcHeader("RegReg")
	for _, k := range [...]arch.Kind{arch.Uint8, arch.Uint16, arch.Uint32, arch.Uint64} {
		if g.opname == "imul" && k.Size() == 1 {
			continue
		}
		fmt.Fprintf(g.w, "\t// reg%d OP= reg%d\n", k.Size()*8, k.Size()*8)
		for src := arch.RLo; src <= arch.RHi; src++ {
			for dst := arch.RLo; dst <= arch.RHi; dst++ {
				fmt.Fprintf(g.w, "\t%s\t%v,%v\n", g.opname, arch.MakeReg(src, k), arch.MakeReg(dst, k))
			}
			fmt.Fprint(g.w, "\tnop\n")
		}
		fmt.Fprint(g.w, "\tnop\n")
	}
	g.funcFooter()
}

func (g *genOp2) opMemConst8Reg() {
	for _, k := range [...]arch.Kind{arch.Uint8, arch.Uint16, arch.Uint32, arch.Uint64} {
		g.opMemConst8RegKind(k)
	}
}

func (g *genOp2) opMemConst8RegKind(k arch.Kind) {
	klen := k.Size() * 8
	if g.opname == "imul" && klen == 8 {
		return
	}
	g.funcHeader(fmt.Sprintf("MemConst8Reg%d", klen))
	offstr := [...]string{"", "0x7F", "0x78563412"}
	for i, offlen := range [...]uint8{0, 8, 32} {
		fmt.Fprintf(g.w, "\t// reg%d = mem%d[off%d] OP const8\n", klen, klen, offlen)
		for src := arch.RLo; src <= arch.RHi; src++ {
			for dst := arch.RLo; dst <= arch.RHi; dst++ {
				for _, cval := range [...]int8{-128, 0, 127} {
					fmt.Fprintf(g.w, "\t%s\t$%#x,%s(%v),%v\n", g.opname, cval,
						offstr[i], arch.MakeReg(src, arch.Uintptr),
						arch.MakeReg(dst, k))
				}
				fmt.Fprint(g.w, "\tnop\n")
			}
			fmt.Fprint(g.w, "\tnop\n")
		}
		fmt.Fprint(g.w, "\tnop\n")
		fmt.Fprint(g.w, "\tnop\n")
	}
	g.funcFooter()
}

func (g *genOp2) opMemReg() {
	for _, k := range [...]arch.Kind{arch.Uint8, arch.Uint16, arch.Uint32, arch.Uint64} {
		g.opMemRegKind(k)
	}
}

func (g *genOp2) opMemRegKind(k arch.Kind) {
	klen := k.Size() * 8
	if g.opname == "imul" && klen == 8 {
		return
	}
	g.funcHeader(fmt.Sprintf("MemReg%d", klen))
	offstr := [...]string{"", "0x7F", "0x78563412"}
	for i, offlen := range [...]uint8{0, 8, 32} {
		fmt.Fprintf(g.w, "\t// reg%d OP= mem%d[off%d]\n", klen, klen, offlen)
		for src := arch.RLo; src <= arch.RHi; src++ {
			for dst := arch.RLo; dst <= arch.RHi; dst++ {
				fmt.Fprintf(g.w, "\t%s\t%s(%v),%v\n", g.opname,
					offstr[i], arch.MakeReg(src, arch.Uintptr),
					arch.MakeReg(dst, k))
			}
			fmt.Fprint(g.w, "\tnop\n")
		}
		fmt.Fprint(g.w, "\tnop\n")
	}
	g.funcFooter()
}

func (g *genOp2) opConstMem() {
	for _, k := range [...]arch.Kind{arch.Uint8, arch.Uint16, arch.Uint32, arch.Uint64} {
		g.opConstMemKind(k, 8)
		if k.Size() != 1 {
			g.opConstMemKind(k, k.Size()*8)
		}
	}
}

func (g *genOp2) opConstMemKind(k arch.Kind, constbits arch.Size) {
	kbits := k.Size() * 8
	g.funcHeader(fmt.Sprintf("Const%dMem%d", constbits, kbits))
	suffixstr := map[arch.Size]string{1: "b", 2: "w", 4: "l", 8: "q"}[k.Size()]
	offstr := [...]string{"", "0x7F", "0x78563412"}
	conststr := map[arch.Size]string{8: "$0x33", 16: "$0xaabb", 32: "$0x11223344", 64: "$0x55667788"}[constbits]
	for i, offlen := range [...]uint8{0, 8, 32} {
		fmt.Fprintf(g.w, "\t// mem%d[off%d] OP= const%d\n", kbits, offlen, kbits)
		for dst := arch.RLo; dst <= arch.RHi; dst++ {
			fmt.Fprintf(g.w, "\t%s%s\t%v,%s(%v)\n", g.opname, suffixstr,
				conststr, offstr[i], arch.MakeReg(dst, arch.Uintptr))
		}
		fmt.Fprint(g.w, "\tnop\n")
	}
	g.funcFooter()
}

func (g *genOp2) opRegMem() {
	for _, k := range [...]arch.Kind{arch.Uint8, arch.Uint16, arch.Uint32, arch.Uint64} {
		g.opRegMemKind(k)
	}
}

func (g *genOp2) opRegMemKind(k arch.Kind) {
	klen := k.Size() * 8
	g.funcHeader(fmt.Sprintf("RegMem%d", klen))
	offstr := [...]string{"", "0x7F", "0x78563412"}
	for i, offlen := range [...]uint8{0, 8, 32} {
		fmt.Fprintf(g.w, "\t// mem%d[off%d] OP= reg%d\n", klen, offlen, klen)
		for src := arch.RLo; src <= arch.RHi; src++ {
			for dst := arch.RLo; dst <= arch.RHi; dst++ {
				fmt.Fprintf(g.w, "\t%s\t%v,%s(%v)\n", g.opname,
					arch.MakeReg(src, k),
					offstr[i], arch.MakeReg(dst, arch.Uintptr))
			}
			fmt.Fprint(g.w, "\tnop\n")
		}
		fmt.Fprint(g.w, "\tnop\n")
	}
	g.funcFooter()
}
