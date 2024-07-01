rpi:
	go build basic.go

basic:
	GOOS=linux GOARCH=arm GOARM=7 \
	go build basic.go
	
test:
	scp basic loragw.local://home/jakobsen/
	ssh loragw.local "sudo -S ./basic" 
	 