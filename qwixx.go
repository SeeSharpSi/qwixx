package main

import (
	qxm "seesharpsi/qwixx/model"
)

// bubble tea programs have:
// model three methods (init, update, view)

// note: use `ssh -p <port#> localhost` to connect

func main() {
	app := qxm.NewApp()
	app.Start()
}

