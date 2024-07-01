
build:
	GOOS=linux GOARCH=arm GOARM=7 \
	go build cp_lora_rx.go
	
run:
	scp cp_lora_rx loragw.local://home/jakobsen/
	ssh loragw.local "sudo -S ./cp_lora_rx"
	
	 
	 