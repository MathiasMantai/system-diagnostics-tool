/*
	System Tool Version 0.0.1
	Author: Mathias Mantai <mmmantaibusiness@gmail.com>

	shown diagnostics:
	- virtual memory usage
	- cpu models

*/


package main 

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/cpu"
	"log"
)

func main() {
	header()
	memory_data()
	cpu_data()
}

func header() {
	divider()
	fmt.Println("SYSTEM DIAGNOSTICS")
	divider()
}

func divider() {
	i := 0
	for i < 30 {
		fmt.Print("=")
		i++
	}
	fmt.Print("\n")
}

func memory_data() {
	memory, err := mem.VirtualMemory()

	if err != nil {
		log.Fatal("Error while accessing virtual memory diagnostics")
	}

	fmt.Println("MEMORY:")
	fmt.Printf("Total Memory: %v \n", memory.Total)
	fmt.Printf("Free Memory: %v \n", memory.Free)
	fmt.Printf("Usage: %v%% \n", memory.UsedPercent)
	divider()
}

func cpu_data() {
	cpuData, err := cpu.Info()

	if err != nil {
		log.Fatal("Error while accessing cpu diagnostics")
	}
	
	//cpu models
	fmt.Println("CPU:")
	for i, cpus := range cpuData {
		fmt.Printf("%v - %v \n", (i + 1), cpus.ModelName)
	}
	divider()
}