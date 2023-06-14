/*
	System Tool Version 0.0.1
	Author: Mathias Mantai <mmmantaibusiness@gmail.com>

	shown diagnostics:
	- virtual memory usage
	- cpu models
	- physical partitions
	- network interfaces
	- 
*/


package main 

import (
	"fmt"
	"math"
	"log"
	"os"

	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/docker/docker/api/types"
    "github.com/docker/docker/client"
    "golang.org/x/net/context"
)

//helper functions

func header() {
	divider()
	fmt.Println("SYSTEM DIAGNOSTICS")
	divider()
}

func divider() {
	i := 0
	for i < 40 {
		fmt.Print("=")
		i++
	}
	fmt.Print("\n")
}

//parse command line arguments
func parse_args(functions map[string]func()) {
	args := get_param()

	//if args are empty just display everything
	if len(args) == 0 {
		display_all(functions)
		return
	}

	//check if help or all argument was put in
	if check_utility_param(args, functions) {
		return 
	}

	//display specific information
	header()
	for _, arg := range args {
		_, exists := functions[arg]
		if exists {
			functions[arg]()
		}
	}
}

func check_utility_param(args []string, functions map[string]func()) bool {
	for _, arg := range args {
		if arg == "-h" || arg == "-help" {
			help_menu()
			return true
		} else if arg == "-a" || arg == "-all" {
			display_all(functions)
			return true
		}
	}

	return false

}

func get_param() []string{
	return os.Args[1:]
}

//display help menu with all arguments
func help_menu() {
	help := "-h/-help  - open this help window  \n-a        - display all information (omitting arguments will also display everything) \n-cp       - cpu information \n-vm       - virtual memory \n-ni       - display information about network interfaces \n-c        - display information for docker container \n"
	fmt.Print(help)
}

func display_all(functions map[string]func()) {
	header()
	for _, function := range functions {
		function()
	}
}


// System Diagnostics

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

	//cpu models
	fmt.Println("CPU:")
	for i, cpus := range cpuData {
		fmt.Printf("%v - %v \n", (i + 1), cpus.ModelName)
		fmt.Printf("   - cores: %v \n", cpus.Cores)
		fmt.Printf("   - mhz: %v \n", cpus.Mhz)
		fmt.Printf("   - cacheSize: %v \n", cpus.CacheSize)
	}

	//display total number of physical and logical cores
	divider()
	cores_total()
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

		fmt.Printf("%v - name: %v \n", i, partition.Device)
		fmt.Printf("   - mountpoint: %v \n", partition.Mountpoint)
		fmt.Printf("   - free: %v \n", usage.Free)
		fmt.Printf("   - used: %v \n", usage.Used)
		fmt.Printf("   - usage in percent: %v%% \n", math.Floor(usage.UsedPercent * 100) / 100)
	}

	divider()
}

//gives information about network interfaces
func net_interfaces() {
	netInterfaces, err := net.Interfaces()

	if err != nil {
		log.Fatal("Error getting network interface information")
	}

	fmt.Println("NETWORK INTERFACES")

	for i, netInterface := range netInterfaces {
		fmt.Printf("%v - name: %v \n", i+1, netInterface.Name)

		//hardware address
		fmt.Printf("  - hardware address: %v \n", netInterface.HardwareAddr)

		//display addresses (ipv4/ipv6) for each network interface
		fmt.Println("  Addresses:")
		for _, address := range netInterface.Addrs {
			fmt.Printf("    - %v \n", address.Addr)
		}

		//MTU 
		fmt.Printf("  - MTU: %v \n", netInterface.MTU)
	}

	divider()
}

/*
	get docker container data
*/
func container_data() {
	cli, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		log.Fatal(err)
	}
    options := types.ContainerListOptions{
        All: true, // Include stopped containers as well
    }

	containers, err := cli.ContainerList(context.Background(), options)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("DOCKER CONTAINER")
	for i, container := range containers {
		fmt.Printf("%v - container id: %v \n", i, container.ID)
		fmt.Printf("   - image: %v \n", container.Image)
		fmt.Printf("   - state: %v \n", container.State)
		fmt.Printf("   - status: %v \n", container.Status)

		if len(container.Ports) > 0 {
			fmt.Println("   - Ports")
			for _, port := range container.Ports {
				fmt.Printf("      - %v (%v) \n", port.PublicPort, port.Type)
			}
		}
	}

	divider()
}



func main() {

	//put all functions inside a map
	functions := map[string]func(){
		"-cp": cpu_data,
		"-pp": physical_partitions,
		"-vm": virtual_memory,
		"-ni": net_interfaces,
		"-c" : container_data,
	}

	//first parse all cli argument entered
	parse_args(functions)


	// header()
	// cpu_data()
	// physical_partitions()
	// virtual_memory()
	// net_interfaces()
	// container_data()
}