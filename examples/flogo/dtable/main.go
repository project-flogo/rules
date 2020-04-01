package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"

	_ "github.com/project-flogo/core/data/expression/script"
	"github.com/project-flogo/core/engine"
)

var (
	cpuProfile    = flag.String("cpuprofile", "", "Writes CPU profile to the specified file")
	memProfile    = flag.String("memprofile", "", "Writes memory profile to the specified file")
	cfgJson       string
	cfgCompressed bool
)

func main() {

	flag.Parse()
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create CPU profiling file: %v\n", err)
			os.Exit(1)
		}
		if err = pprof.StartCPUProfile(f); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to start CPU profiling: %v\n", err)
			os.Exit(1)
		}
		defer pprof.StopCPUProfile()
	}

	cfg, err := engine.LoadAppConfig(cfgJson, cfgCompressed)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create engine: %v\n", err)
		os.Exit(1)
	}

	e, err := engine.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create engine: %v\n", err)
		os.Exit(1)
	}

	code := engine.RunEngine(e)

	if *memProfile != "" {
		f, err := os.Create(*memProfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create memory profiling file: %v\n", err)
			os.Exit(1)
		}

		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write memory profiling data: %v", err)
			os.Exit(1)
		}
		_ = f.Close()
	}

	os.Exit(code)
}
