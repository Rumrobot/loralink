rpi:
	go build hello.go

hello:
	GOOS=linux GOARCH=arm GOARM=7 \
	go build hello.go
	
test:
	scp hello loragw.local://home/jakobsen/
	ssh loragw.local ./hello 
	 