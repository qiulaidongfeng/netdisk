package netdisk

import (
	"fmt"
	"sync/atomic"

	"github.com/fsnotify/fsnotify"
	"github.com/go-viper/encoding/ini"
	"github.com/spf13/viper"
)

var v *viper.Viper

func init() {
	v = newv()
	loadConfig()
}

func newv() *viper.Viper {
	codecRegistry := viper.NewCodecRegistry()
	codecRegistry.RegisterCodec("ini", ini.Codec{})
	v := viper.NewWithOptions(viper.WithCodecRegistry(codecRegistry))
	v.SetConfigFile("config.ini")
	v.OnConfigChange(func(e fsnotify.Event) {
		loadConfig()
		fmt.Println("Config file changed:", e.Name)
	})
	v.WatchConfig()
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	return v
}

var config = struct {
	mysqluser     atomic.Pointer[string]
	mysqlpassword atomic.Pointer[string]
	mysqladdr     atomic.Pointer[string]
}{}

func loadConfig() {
	config.mysqluser.Store(ptr(v.GetString("mysql.user")))
	config.mysqlpassword.Store(ptr(v.GetString("mysql.password")))
	config.mysqladdr.Store(ptr(v.GetString("mysql.addr")))
}

func ptr(v string) *string {
	return &v
}

func getDsnInfo() (user, password, addr string) {
	return *config.mysqluser.Load(), *config.mysqlpassword.Load(), *config.mysqladdr.Load()
}
