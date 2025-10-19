package netdisk

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/qiulaidongfeng/safesession"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var db *gorm.DB

func init() {
	User, password, addr := getDsnInfo()
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/netdisk?charset=utf8mb4&parseTime=True&loc=Local", User, password, addr)
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(100)
	err = db.AutoMigrate(&user{}, &safesession.Session{})
	if err != nil {
		panic(err)
	}
}

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

// init_fileDb 为一位用户初始化它的存储数据库
func init_fileDb(id string) {
	err := db.Scopes(tableForUser(id)).AutoMigrate(&fileEntry{})
	if err != nil {
		panic(err)
	}
}

// new_fileDb 获取管理所有用户文件的数据库
func new_fileDb() fileDb {
	return &mysqlFileDb{}
}

// TODO: 使用更好的读写方法
// gorm不支持流式传输
// 可改用 SUBSTRING 和 LENGTH 流式读取数据库
// 改用 database/sql + go-sql-driver/mysql 流式写入数据库
type mysqlFileDb struct{}

var _ fileDb = (*mysqlFileDb)(nil)

// fileEntry 是一个存在数据库的文件
type fileEntry struct {
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Path      string    `gorm:"primaryKey"`
	Size      int
	Data      []byte `gorm:"type:LONGBLOB"`
}

func (m *mysqlFileDb) SetLimit(id string, size int) {
	result := db.Model(&user{}).Where("id = ?", id).Update("limit", size)
	if result.Error != nil {
		panic(result.Error)
	}
}

func (m *mysqlFileDb) Get(id, path string) io.Reader {
	db := db.Scopes(tableForUser(id))
	f := &fileEntry{Path: path}
	result := db.First(f)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil
		}
		panic(result.Error)
	}
	return bytes.NewReader(f.Data)
}

func (m *mysqlFileDb) Delete(id, path string) bool {
	f := &fileEntry{Path: path}
	err := db.Transaction(func(tx *gorm.DB) error {
		result := tx.Scopes(tableForUser(id)).First(f)
		if result.Error != nil {
			return result.Error
		}

		// TODO: use SET -=
		var u user
		u.Id = id
		result = tx.First(&u)
		if result.Error != nil {
			return result.Error
		}

		u.Used -= f.Size
		result = tx.Model(&u).Where(&u).Update("used", u.Used)
		if result.Error != nil {
			return result.Error
		}

		result = tx.Scopes(tableForUser(id)).Delete(f)
		return result.Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false
		}
		panic(err)
	}

	return true
}

func (m *mysqlFileDb) Set(id, path string, data io.Reader) bool {
	var u user
	u.Id = id
	result := db.Select("limit").First(&u)
	if result.Error != nil {
		panic(result.Error)
	}
	s, err := io.ReadAll(data)
	if err != nil {
		panic(err)
	}
	if u.Used+len(s) > u.Limit {
		return false
	}
	err = db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&user{}).Where("id = ?", id).Update("used", u.Used+len(s))
		if result.Error != nil {
			return result.Error
		}

		tx = tx.Scopes(tableForUser(id))
		f := &fileEntry{Path: path, Size: len(s), Data: s}
		//文件不存在就保存，存在就修改
		result = tx.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(f)
		return result.Error
	})
	if err != nil {
		panic(err)
	}
	return true
}

func (m *mysqlFileDb) List(id string) []fileEntry {
	db := db.Scopes(tableForUser(id))
	var s []fileEntry
	result := db.Select("updated_at", "path", "size").Find(&s)
	if result.Error != nil {
		panic(result.Error)
	}
	return s
}
