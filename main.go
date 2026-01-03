package main

import (
	"duifene_auto_sign/backend"
	"log"
)

func main() {
	router := backend.SetupRouter()

	// 启动服务器
	log.Println("服务器启动在 http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
