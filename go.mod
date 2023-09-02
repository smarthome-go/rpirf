module github.com/smarthome-go/rpirf

go 1.17

require (
	github.com/stianeikeland/go-rpio/v4 v4.6.0
	github.com/warthog618/gpiod v0.8.2
)

require golang.org/x/sys v0.10.0 // indirect

replace rpirf => ../rpirf/
