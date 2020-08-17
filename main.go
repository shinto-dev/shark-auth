package main

import (
	"shark-auth/apis"
)

func main() {
	router := apis.Router()
	router.Run(":8080")
}
