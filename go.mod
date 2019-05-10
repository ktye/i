module github.com/ktye/i

go 1.12

require (
	github.com/eaburns/T v0.0.0-20190217122806-dbc7887ff15c
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/ktye/iv v0.0.0-20190506155233-ea1af1e66e83
	github.com/ktye/plot v0.0.0
	github.com/ktye/ui v1.0.0 // this is wrong: should be v0.0.0
	golang.org/x/exp v0.0.0-20190510051728-0e2d8f6bf8da
	golang.org/x/image v0.0.0-20190507092727-e4e5bf290fec
	golang.org/x/mobile v0.0.0-20190509164839-32b2708ab171
)

replace github.com/ktye/ui => ../ui

replace github.com/ktye/plot => ../plot
