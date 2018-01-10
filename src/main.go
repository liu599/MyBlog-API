package main

import "os"

func main() {

	configuration()

	a := App{}

	a.Initialize(os.Getenv("NEKOHAND_ADMINISTRATOR"), os.Getenv("NEKOHAND_PASSWORD"), os.Getenv("NEKOHAND_DATABASE_NAME"))

	a.Run(os.Getenv("NEKOHAND_APP_ADDR"))
}
