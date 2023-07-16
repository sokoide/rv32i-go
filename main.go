package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/sokoide/rv32i-go/pkg/rv32i"
)

func chkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	log.Info("* Started")
	// test binary
	// 80000000 <boot>:
	// 80000000: 93 00 00 00   li      ra, 0
	// # x1 == ra == 0
	// 80000004: 13 04 00 00   li      s0, 0
	// # x8 == s0/fp == 0
	// 80000008: 37 45 00 00   lui     a0, 4
	// # x10 == a0 == 4<<12 == 16384
	p, err := rv32i.ReadBinary("./data/sample-binary-001.txt")
	chkerr(err)

	f := rv32i.NewFetcher(p)
	d := rv32i.NewDecoder()

	log.Infof("f: %+v", f)
	log.Infof("d: %+v", d)

	for i := 0; i < 3; i++ {
		u32instr, _ := f.Fetch()
		log.Infof("u32instr: %08x", u32instr)
		instr := d.Decode(u32instr)
		log.Infof("instr: %+v", instr)
	}

	log.Info("* Completed")
}
