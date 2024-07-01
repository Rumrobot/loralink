package main

import (
	"fmt"
	"log"
	"strings"
	// "os/signal"
	"time"
//	"bufio"
//	"net"
//    "os"	
	"bytes"
	"encoding/json"
	"net/http"
        
	"lora"
)

type RequestData struct {
	auth string `json:"auth"`
	data  string `json:"json:"data"`
}
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
	var lora lora.LORA
	fmt.Println("Starting LORA gateway")
	postData := RequestData{
   		auth: "k2j39s92k!",
	}

	//defer res.Body.Close()

	//c, err := net.Dial("tcp", "anyvej11.dk:8087")
    //if err != nil {
	// 		fmt.Println(err)
	// 		return
	// }
	// fmt.Fprintf(c, "Hello from LORAgateway\n")
	
	err := lora.Open("/dev/spidev0.0")
	if err != nil {
		log.Fatal(err)
	}
	
	 
	client := &http.Client{}
	// client.Do(r)
	if err != nil {
		panic(err)
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
	var data string 
	
	var debug int = 0	
    for (true ) {	
		fmt.Print(".")
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
			data = lora.Getpacket()
		    fmt.Printf("Data: %s\n", data)
			// split on space 
			data_arr := strings.Split(data, " ")
			if len(data_arr) == 3 {	
				fmt.Printf("Data: %s\n", data_arr[2]) 
				
				
				// Insert the variable into the data field
				postData.data = data_arr[2]
			
				//jsonData, err := json.Marshal(postData)
				if err != nil {
					fmt.Println("Error marshalling JSON data:", err)
					return
				}
			
				posturl := "https://rumrobot.dk/temp" // Replace with the actual URL
			
				req, err := http.NewRequest("POST", posturl, bytes.NewBuffer(postData))
				if err != nil {
					fmt.Println("Error creating HTTP request:", err)
					return
				}
			
				req.Header.Add("Content-Type", "application/json")
				fmt.Printf("%s ", postData)
				client.Do(req)
				if err != nil {
				panic(err)
				
				}
			}
		}
	
		time.Sleep(8 * time.Second) 
	} 
}
