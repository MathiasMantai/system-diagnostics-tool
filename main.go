package main 

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/mem"
	"log"
)

func main() {
	memory()
}

func memory() {
	memory, error := mem.VirtualMemory()

	if error != nil {
		log.Fatal("Error while accessing virtual memory")
	}

	fmt.Printf("Total Memory:%v \n", memory.Total)
	fmt.Printf("%v", memory.Free)
	fmt.Printf("Usage: %v ", memory.UsedPercent)
}