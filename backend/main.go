package main

import "github.com/raynard2/backend/router"

func main()  {
	e := router.New()

	e.Logger.Fatal(e.Start(":5001"))
}
