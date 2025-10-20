# netdisk

一个简易网盘，设计在多个低价vps运行的分布式网盘

## 设计目标

1. 能在低价vps(小于5美元1个月)运行
2. 支持分布式部署

## 当前最低配置要求

内存：300mb
CPU：单核
数据库：mysql
go语言：go1.26

## 快速开始

克隆代码

创建配置文件config.ini
内容要求

```
[mysql]
user=mysql用户名
password=mysql密码
addr=mysql地址

[db]
mode=os
```

```
npm install
npm run build

cd ./netdisk
mv ./netdisk ./netdisk1
mv ./dist ./netdisk 

go run main.go
```

## 代办实现

- 流式读写mysql
- AES-256-GCM加密存储数据
- 支持使用MongoDB
- 支持存储文件夹
- 支持使用MinIO

## 实现原理

前端基于 Web 实现，可在主流平台浏览器（Windows/macOS/Linux/iOS/Android）上使用。

go语言实现后端。

通过下列存储接口可以灵活适配不同存储方案（mysql,MongoDB,MinIO等）。

```go
// fileDb 表示网盘的存储数据库实现
// 所有方法并发调用是安全的
type fileDb interface {
	// 下列id指的是用户id

	// SetLimit 设置用户的可用空间上限
	// size单位字节
	SetLimit(id string, size int)
	// Get 获取用户存在网盘的一个文件
	// 如果返回nil，表示没有这个文件
	Get(id, path string) io.Reader
	// Set 将用户的一个文件存入网盘
	// 如果返回false表示保存失败，因为保存会超过用户可用空间上限
	Set(id, path string, data io.Reader) bool
	// List 查询用户存的所有文件
	// 返回值仅包含大小等元数据，没有保存的文件
	List(id string) []fileEntry
	// Delete 删除用户保存的一个文件
	// 如果返回false，表示没有这个文件
	Delete(id, path string) bool
}

// fileEntry 是一个存在数据库的文件
type fileEntry struct {
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Path      string    `gorm:"primaryKey"`
	Size      int
	Data      []byte `gorm:"type:LONGBLOB"`
}
```

目前使用mysql存储数据，可以通过ySQL InnoDB Cluster实现多地自动备份。
未来支持MongoDB+副本集，可以更高效的存储并备份数据。
使用MinIO也是一个值得研究的方向。

go后端使用gin处理根据（`/list /stat /upload /download`）等api接口，与前端交互。
通过上述存储接口实现网盘所有功能。

## 适用场景

在私人或熟人之间的小范围使用的网盘，对数据自动多地备份有需求。

此处备份=**一个文件+至少一个副本**

如果只是要做单机自建网盘，可以去使用[蓝眼云盘](https://github.com/eyebluecn/tank)，可以快速搭起网盘。

## 如何自己做新的存储接口实现

实现上述存储接口，如何将db.go文件的`init_fileDb`和`new_fileDb`两个函数改为新接口即可。