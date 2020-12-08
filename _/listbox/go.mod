module listbox

go 1.15

require (
	github.com/ktye/plot v0.0.0
	github.com/ktye/plot/plotui v0.0.0
	github.com/ktye/pptx v0.0.0
	github.com/ktye/pptx/pptxt v0.0.0
	github.com/lxn/walk v0.0.0-20201125094449-2a61ddb5a2b8
	github.com/lxn/win v0.0.0-20201111105847-2a20daff6a55 // indirect
	golang.org/x/tools v0.0.0-20201206230334-368bee879bfd // indirect
	gopkg.in/Knetic/govaluate.v3 v3.0.0 // indirect
)

replace (
	github.com/ktye/pptx => c:/k/pptx
	github.com/ktye/pptx/pptxt => c:/k/pptx/pptxt
	github.com/ktye/plot => c:/k/plot
	github.com/ktye/plot/plotui => c:/k/plot/plotui
)
