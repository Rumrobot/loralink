package lora

import (
	
	"fmt"
	"log"
	"time"
	
	"periph.io/x/host/v3"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/gpio"
    "periph.io/x/conn/v3/gpio/gpioreg"

)

const ( Mode = 0x01
	    OpMode = 0x07
)

type LORA struct 
{
	port spi.PortCloser
	dev		     spi.Conn 
	Thermocouple float64
	Internal     float64
	Timestamp    time.Time
}

// Open establishes the SPI connection with the sensor.
func (m *LORA) Open(name string) error {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Use gpioreg GPIO port registry to find the first available GPIO pin.	
 	p := gpioreg.ByName("GPIO22")
    if p == nil {
        log.Fatal("Failed to find GPIO22")
    }
    fmt.Printf("%s: %s\n", p, p.Function())
	if err := p.In(gpio.Float,gpio.NoEdge); err != nil {
		log.Fatal(err)
    }
	fmt.Printf("%s: %s\n", p, p.Function())
	
	// Use spireg SPI port registry to find the first available SPI bus.
	port, err := spireg.Open("")
	if err != nil {
		log.Fatal(err)
	}
	//defer p.Close()

	// Convert the spi.Port into a spi.Conn so it can be used for communication.
	dev, err := port.Connect(physic.MegaHertz, spi.Mode3, 8)
	if err != nil {
		log.Fatal(err)
	} 
	
	write := []byte{0x0a, 0x00}
	read := make([]byte, len(write))
	
	if err := dev.Tx(write, read); err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("%x\n", read[1:])
	
	m.dev = dev

	return err
}
func (m *LORA) Getpacket() {
	
	write := []byte{0x13, 0x00}
	read := make([]byte, len(write))
	if err := m.dev.Tx(write, read); err != nil {
		log.Fatal(err)
	}
	n := read[1] 
	
	// Get RSSI value
	write = []byte{0x1a, 0x00}
	read  = make([]byte, len(write))
	if err := m.dev.Tx(write, read); err != nil {
		log.Fatal(err)
	}
	rssi := int(read[1])

	//fmt.Printf("Number of bytes received: %d\n", n)
	
	// Stop rx mode 
	// Set LORA mode + standby
	//m.SetLORAmode(0x01)
	
	//Get fifo rx current addr pointer
	// write = []byte{0x10, 0x00}	
	// read  = make([]byte, len(write))
	// if err := m.dev.Tx(write, read); err != nil {
	 	// log.Fatal(err)
	// }
	// current_rx_addr := read[1] 
	// fmt.Printf("Current RX addr: %d\n", current_rx_addr)
	// m.Write(0x0d, current_rx_addr)
    // //res := make([]byte, n)
	if n>0 { 
	//for i:=0; i<int(n); i++ {
		write := make([]byte, n+1) 
		read  := make([]byte, len(write))
		if err := m.dev.Tx(write, read); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s : RSSI (%d)\n", read[1:], int(-164+rssi))
  	}
}

func (m *LORA) Reset() {
	
	p := gpioreg.ByName("GPIO22")
    if p == nil {
        log.Fatal("Failed to find GPIO22")
    }
    fmt.Printf("%s: %s\n", p, p.Function())
	if err := p.Out(gpio.Low); err != nil {
		log.Fatal(err)
    } 
	time.Sleep(1 * time.Second) 
	fmt.Printf("%s: %s\n", p, p.Function())
	if err := p.In(gpio.Float,gpio.NoEdge); err != nil {
		log.Fatal(err)
    }
}

func (m *LORA) SetLORAmode(mode byte) {
	err := m.Write(0x01, 0x80 | mode) 
	if err != nil {
		fmt.Println("Error setting LORA mode")
	}
}

func (m *LORA) SetFrequency(frequency uint32) {
	var frf = (uint64(frequency) << 19) / 32000000
	m.Write(SX127X_REG_FRF_MSB, uint8(frf>>16))
	m.Write(SX127X_REG_FRF_MID, uint8(frf>>8))
	m.Write(SX127X_REG_FRF_LSB, uint8(frf>>0))
}
func (m *LORA) SetSyncWord(SyncWord byte) {
	err := m.Write(SX127X_REG_SYNC_WORD, SyncWord)
	if err != nil {
		fmt.Println("Error setting syncword")
	}
}

func (m *LORA) SetSpreadingFactor(value byte) error {
	pos := 4;
	err := m.Write(SX127X_REG_MODEM_CONFIG_2, value<<pos)
	if err != nil {
		fmt.Println("Error setting Spreading Factor")
		return err
	}
	return nil
}

func (m *LORA) SetCodingRate(value byte) error {
	pos := 1;	
	write := []byte{SX127X_REG_MODEM_CONFIG_1, 0x00}
	read := make([]byte, len(write))
	err := m.dev.Tx(write,read)
    if err != nil {
		fmt.Println("Error setting SetCodingRate")
		return err
    }
	// Complete value in read[1]
	read[1] = (read[1] & ^(0x07 << pos)) | (value << pos)
	newvalue := read[1]
	fmt.Printf("New Value in %x : %x\n",SX127X_REG_MODEM_CONFIG_1, newvalue)
	err = m.Write(SX127X_REG_MODEM_CONFIG_1, newvalue)
	if err != nil {
		fmt.Println("Error setting CodingRate")
		return err
	}
	return nil
}

func (m *LORA) SetBanddwidth(value byte) error {
	pos := 4;
	//len := 4;
	write := []byte{SX127X_REG_MODEM_CONFIG_1, 0x00}
	read := make([]byte, len(write))
	err := m.dev.Tx(write,read)
    if err != nil {
		fmt.Println("Error setting Bandwidth")
		return err
    }
	// Complete value in read[1]
	read[1] = (read[1] & ^(0x0f << pos)) | (value << pos)
	newvalue := read[1]
	fmt.Printf("New Value in %x : %x\n",SX127X_REG_MODEM_CONFIG_1, newvalue)
	err = m.Write(SX127X_REG_MODEM_CONFIG_1, newvalue)
	if err != nil {
		fmt.Println("Error setting Bandwidth")
		return err
	}
	return nil
}

// Read gets the latest values from the sensor,
// and updates the thermocouple, internal and Timestamp fields.
func (m *LORA) Read(addr byte) (byte, error) {
	write := []byte{addr, 0x00}
	read := make([]byte, len(write))
	
	err := m.dev.Tx(write,read)
	
	if err != nil {
		fmt.Printf("Error reading from SPI %x\n", addr)
		return 0x00 , err
    }
	//fmt.Printf("Read (%02x):%02x\n", addr ,read[1:] )
	return read[1], err
}

func (m *LORA) Write(addr byte, value byte) error {
	
	err := m.dev.Tx([]byte{addr | 0x80 , value}, nil )
	if err != nil {
		fmt.Println("Error writing to SPI")
		return err
	}
	return nil
}

// Close ends the SPI connection to thermocouple
func (m *LORA) Close() error {
	m.port.Close()
	return nil 
}
