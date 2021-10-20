package main

import (
	"log"

	"aitsuki.com/pixiv-capture/app"
)

func main() {
	// var port int
	// var username string
	// var password string
	// flag.IntVar(&port, "P", 8080, "端口")
	// flag.StringVar(&username, "u", "", "用户名")
	// flag.StringVar(&password, "p", "", "密码")
	// flag.Parse()
	// if len(username) == 0 || len(password) == 0 {
	// 	log.Fatal("请设置管理员账户和密码")
	// }
	log.Fatal(app.Run(8080, "./pixiv.db"))
}
