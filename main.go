package main

import "nyooom/logging"

func main() {
	logging.ReadConfig() // Setup printing colors
	var _ AdvancedDB = SetupDB()
	logging.Println("Hello, World")
}
