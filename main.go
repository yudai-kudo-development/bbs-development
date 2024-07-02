package main

import (
	"fmt"
	"bbs-development/app/controllers"
)

func main () {
	fmt.Println("サーバー起動処理始めます")
	controllers.StartMainServer()
}

