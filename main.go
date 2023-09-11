// main.go

package main

import (
	"fmt"
	"log"
	"nacos-go/nacos_registry"
	"net/http"
)

func main() {
	// 创建 NacosRegistry 实例
	nacosRegistry, err := nacos_registry.NewNacosRegistry()
	if err != nil {
		log.Fatal(err)
	}

	// 注册服务
	if err := nacosRegistry.RegisterService(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Nacos!")
	})

	// 启动 HTTP 服务器

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
	// 启动服务
	go func() {
		nacosRegistry.Run()
	}()

}
