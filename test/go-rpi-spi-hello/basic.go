package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
	
	"lora"
)
const (
	CodingRate4_5 = 0x01 //  7     0     LoRa coding rate: 4/5
	CodingRate4_6 = 0x02 //  7     0                       4/6
	CodingRate4_7 = 0x03 //  7     0                       4/7
	CodingRate4_8 = 0x04 //  7     0                       4/8
)
const (
	SpreadingFactor5  = 0x05
	SpreadingFactor6  = 0x06
	SpreadingFactor7  = 0x07
	SpreadingFactor8  = 0x08
	SpreadingFactor9  = 0x09
	SpreadingFactor10 = 0x0A
	SpreadingFactor11 = 0x0B
	SpreadingFactor12 = 0x0C
)
const (
	Bandwidth_7_8   = iota // 7.8 kHz
	Bandwidth_10_4         // 10.4 kHz
	Bandwidth_15_6         // 15.6 kHz
	Bandwidth_20_8         // 20.8 kHz
	Bandwidth_31_25        // 31.25 kHz
	Bandwidth_41_7         // 41.7 kHz
	Bandwidth_62_5         // 62.5 kHz
	Bandwidth_125_0        // 125.0 kHz
	Bandwidth_250_0        // 250.0 kHz
	Bandwidth_500_0        // 500.0 kHz
)

func main() {
	interrupt := make(chan os.Signal, 1) // for catching ctrl+c
	stop := make(chan bool, 1)
	signal.Notify(interrupt, os.Interrupt)
	var lora lora.LORA
	
	err := lora.Open("/dev/spidev0.0")
	if err != nil {
		log.Fatal(err)
	}
	
	// Set lora mode + standby 
	fmt.Println("Setting mode LORA + standby" )
    lora.SetLORAmode(0x0)
	lora.SetLORAmode(0x1)
	fmt.Println("Setting mode LORA + standby [done]")

	fmt.Println("Force reset")
	lora.Reset()

	lora.SetLORAmode(0x0)
	lora.SetLORAmode(0x1)
	
	lora.SetFrequency(436000000)
 	lora.SetBanddwidth(Bandwidth_31_25)
	lora.SetSpreadingFactor(SpreadingFactor10)
	lora.SetCodingRate(CodingRate4_5)
	lora.SetSyncWord(0x12)
	
	lora.SetLORAmode(0x5)
	
	var debug int = 0	

	ticker := time.NewTicker(1000 * time.Millisecond)
	go func() {
		for {
			select {
			case <-stop:
				log.Println("exiting loop")
				return
			case <-ticker.C:
				irq , _ := lora.Read(0x12) 
				if irq & 0x50 == 0x50 {
					if (debug == 1) {
						lora.Read(0x0d)
						lora.Read(0x10)
						lora.Read(0x13)
						lora.Read(0x15)
						lora.Read(0x18)
					}
					lora.Write(0x12, 0x50)
					lora.Write(0x12, 0x50)
					lora.Getpacket()
				}
				//lora.SetLORAmode(0x05)
			}
		}
	}()
	// more stuff could go here, e.g. push data to database
	<-interrupt
	stop <- true
	lora.Close()
}
