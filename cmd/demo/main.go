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
	var err error

	log.SetOutput(os.Stdout)
	// log.SetLevel(log.TraceLevel)
	log.SetLevel(log.InfoLevel)

	log.Info("* Started")
	log.Info("* Running a small program")

	emu := rv32i.NewEmulator()
	emu.Reset()

	// test binary in sample-binary-001.txt
	// 80000000 <boot>:
	// 80000000: 93 00 00 00   li      ra, 0
	// # x1 == ra == 0
	// 80000004: 13 04 00 00   li      s0, 0
	// # x8 == s0/fp == 0
	// 80000008: 37 45 00 00   lui     a0, 4
	// # x10 == a0 == 4<<12 == 16384
	err = emu.Load("./data/sample-binary-001.txt")
	chkerr(err)

	emu.Step()
	emu.Step()
	emu.Step()
	emu.Dump()

	log.Info("* Converting txt to bin")
	l := rv32i.NewLoader()
	l.TextToBinary("./data/sample-binary-003.txt", "./data/sample-binary-003.bin")

	log.Info("* Running the bin")
	emu.Reset()
	err = emu.Load("./data/sample-binary-003.bin")
	chkerr(err)

	// emu.Run()
	emu.StepUntil(0x18)
	emu.Dump()
	emu.StepUntil(0x1c)

	log.Info("* Completed")
}
