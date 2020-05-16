module github.com/junland/conveyor

go 1.13

require (
	github.com/julienschmidt/httprouter v1.3.0
	github.com/justinas/alice v1.2.0
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/pflag v1.0.5
	github.com/junland/conveyor/server v1.0.0
	github.com/junland/conveyor/cmd v1.0.0
)

replace github.com/junland/conveyor/server => ./conveyor/server
replace github.com/junland/conveyor/cmd => ./conveyor/cmd