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
	"math"
	"log"

	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/disk"
)

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

func virtual_memory() {
	memory, err := mem.VirtualMemory()

	if err != nil {
		log.Fatal("Error while accessing virtual memory diagnostics")
	}

	fmt.Println("VIRTUAL MEMORY:")
	fmt.Printf("total Memory: %v \n", memory.Total)
	fmt.Printf("free Memory: %v \n", memory.Free)
	fmt.Printf("usage in percent: %v%% \n", math.Floor(memory.UsedPercent * 100) / 100)
	divider()
}

func cpu_data() {

	//general information about cpu
	cpuData, err := cpu.Info()

	if err != nil {
		log.Fatal("Error while accessing cpu diagnostics")
	}

	//usage data of cpu
	cpuUsageData, err := cpu.Percent(100, false)
	if err != nil {
		log.Fatal("Error while accessing cpu usage")
	}

	//cpu models
	fmt.Println("CPU:")
	for i, cpus := range cpuData {
		fmt.Printf("%v - %v \n", (i + 1), cpus.ModelName)
		fmt.Printf("  - cores: %v \n", cpus.Cores)
		fmt.Printf("  - mhz: %v \n", cpus.Mhz)
		fmt.Printf("  - cacheSize: %v \n", cpus.CacheSize)
	}

	divider()

	cores_total()

	divider()
}

func load_data() {
	loadData, err := load.Misc()

	if err != nil {
		log.Fatal("Error while accessing load data")
	}

	fmt.Println(loadData)
}

func process_data() {

}

func physical_partitions() {
	fmt.Println("PHYSICAL PARTITIONS:")
	partitions, err := disk.Partitions(false)

	if err != nil {
		log.Fatal("Error getting system partitions")
	}

	for i, partition := range partitions {


		usage, err := disk.Usage(partition.Mountpoint)

		if err != nil {
			log.Fatalf("Error getting usage stats for partition %v", partition.Device)
		}

		fmt.Printf("%v - %v \n", i, partition.Device)
		fmt.Printf("  - mountpoint: %v \n", partition.Mountpoint)
		fmt.Printf("  - free: %v \n", usage.Free)
		fmt.Printf("  - used: %v \n", usage.Used)
		fmt.Printf("  - usage in percent: %v%% \n", math.Floor(usage.UsedPercent * 100) / 100)
	}

	divider()
}

func cores_total() {
	phys, err := cpu.Counts(false)
	if err != nil {
		log.Fatal("Error getting physical number of cores")
	}

	logic, err := cpu.Counts(true)
	if err != nil {
		log.Fatal("Error getting logical number of cores")
	}

	fmt.Printf("physical cores total: %v \n", phys)
	fmt.Printf("logical cores total: %v \n", logic)
}

func main() {
	header()
	cpu_data()
	physical_partitions()
	virtual_memory()
}