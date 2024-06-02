package main

import (
	"flag"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/sokoide/rv32i-go/pkg/rv32i"
)

func chkerr(err error) {
	if err != nil {
		panic(err)
	}
}

type options struct {
	logLevel   string
	demo       bool
	sourcePath string
	end        string
}

var opts options = options{
	demo: false,
}

func parseArgs() {
	flag.BoolVar(&opts.demo, "demo", false, "Run demo")
	flag.StringVar(&opts.logLevel, "logLevel", opts.logLevel, "Log level (trace, debug, info, warn, error, fatal, panic)")
	flag.StringVar(&opts.sourcePath, "sourcePath", opts.sourcePath, "Source path")
	flag.StringVar(&opts.end, "end", opts.end, "End address")
	flag.Parse()
}

func runDemo() {
	var err error

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

	log.Info("*** code ***")
	for i := 0; i < 0x108; i += 4 {
		log.Infof("0x%08x: 0x%02x%02x%02x%02x,", i, emu.Memory[i+3], emu.Memory[i+2], emu.Memory[i+1], emu.Memory[i])
	}
	log.Info("*** code end ***")

	emu.StepUntil(0x18)
	emu.Dump()
	emu.StepUntil(0x1c)
}

func run(sourcePath string, end string) {
	var err error

	emu := rv32i.NewEmulator()
	emu.Reset()

	err = emu.Load(sourcePath)
	chkerr(err)

	if len(end) > 0 {
		var uintEnd uint32
		var uintEnd64 uint64
		if strings.HasPrefix(end, "0x") {
			uintEnd64, err = strconv.ParseUint(end[2:], 16, 32)
			uintEnd = uint32(uintEnd64)
		} else {
			uintEnd64, err = strconv.ParseUint(end, 10, 32)
			uintEnd = uint32(uintEnd64)
		}

		chkerr(err)

		startTime := time.Now()
		err = emu.StepUntil(uintEnd)
		endTime := time.Now()
		log.Infof("elapsed time: %v", endTime.Sub(startTime))
	} else {
		panic("please specify the end address")
	}
}

func main() {
	var err error

	parseArgs()

	log.SetOutput(os.Stdout)
	ll := log.InfoLevel
	if len(opts.logLevel) > 0 {
		ll, err = log.ParseLevel(opts.logLevel)
		chkerr(err)
	}
	log.SetLevel(ll)

	log.Info("* Started")
	log.Info("* Running a small program")

	if opts.demo {
		runDemo()
	} else {
		run(opts.sourcePath, opts.end)
	}

	log.Info("* Completed")
}
