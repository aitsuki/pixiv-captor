package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/aitsuki/pixiv-captor/services"
)

func main() {
	var port int
	var username string
	var password string
	var dbPath string
	var logPath string
	flag.IntVar(&port, "P", 8080, "端口")
	flag.StringVar(&username, "u", "", "用户名")
	flag.StringVar(&password, "p", "", "密码")
	flag.StringVar(&dbPath, "db", relativePath("pixiv.db"), "数据库保存位置")
	flag.StringVar(&logPath, "log", relativePath("pixiv.log"), "日志保存位置")
	flag.Parse()
	if len(username) == 0 || len(password) == 0 {
		log.Fatal("为了安全性考虑，请设置管理员账户和密码。")
	}
	log.Fatal(services.Run(port, dbPath, logPath, username, password))
}

func relativePath(name string) string {
	if filepath.IsAbs(name) {
		return name
	}
	e, _ := os.Executable()
	return filepath.Join(filepath.Dir(e), name)
}
