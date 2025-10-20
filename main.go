package main

import (
	"github.com/qiulaidongfeng/netdisk/netdisk"
	"github.com/gin-gonic/gin"
)

func main() {
	s := gin.Default()
	netdisk.Route(s)
	s.Run(":80")
}
