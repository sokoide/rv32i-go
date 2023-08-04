package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/sokoide/rv32i-go/pkg/rv32iasm"
)

type Options struct {
	source string
}

var o *Options = &Options{}

func parseFlags() {
	flag.StringVar(&o.source, "source", "", "source file")
	flag.Parse()
}

func main() {
	var err error

	log.SetOutput(os.Stderr)
	// log.SetLevel(log.TraceLevel)
	log.SetLevel(log.InfoLevel)

	parseFlags()

	log.Info("asm started")

	reader, err := os.Open(o.source)
	if err != nil {
		log.Fatalf("faled to open %s", o.source)
	}

	ev := rv32iasm.NewEvaluator()

	code, err := ev.Assemble(reader)
	if err != nil {
		log.Fatalf("EvaluateProgram error: %v", err)
	}

	log.Info("asm completed")
	for _, line := range code {
		fmt.Println(line)
	}
}
