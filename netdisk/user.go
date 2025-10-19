package netdisk

import (
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/qiulaidongfeng/nonamevote/nonamevote"
	"github.com/qiulaidongfeng/safesession"
	"gorm.io/gorm"
)

type user struct {
	Name     string `gorm:"index"`
	Id       string `gorm:"primaryKey"`
	Password string // 存储sha512的base64哈希值
	Session1 string
	Session2 string
	Limit    int
	Used     int
}

func new_user(name, password string) string {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	ps := hash(password)
	id_base64 := base64.StdEncoding.EncodeToString(id[:])
	//默认每位用户最多可用50Mb
	u := user{Name: name, Id: id_base64, Password: ps, Limit: 50 * 1024 * 1024}
	result := db.Create(&u)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) { //如果id已经存在
			new_user(name, password)
		} else {
			panic(result.Error)
		}
	}
	init_fileDb(id_base64)
	return id_base64
}

func tableForUser(id string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table(id + "_file")
	}
}

func hash(password string) string {
	tmp := sha512.Sum512([]byte(password))
	ps := base64.StdEncoding.EncodeToString(tmp[:])
	return ps
}

func login(ctx *gin.Context, userid, password string) (bool, string) {
	ps := hash(password)
	u := user{Id: userid}
	result := db.First(&u)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, ""
		}
		panic(result.Error)
	}
	if u.Password != ps {
		return false, ""
	}
	se := sessionControl.NewSession(ctx.ClientIP(), ctx.Request.UserAgent(), userid)
	sessionControl.SetSession(&se, ctx.Writer)
	add_session(&u, se.ID)
	add_login_cookie(ctx, userid, u.Name)
	return true, userid
}

var sessionAge = 365 * 24 * time.Hour

var sessionControl = safesession.NewControl(nonamevote.GetAeskey(), sessionAge, 0, func(clientIp string) safesession.IPInfo {
	//TODO:实现这里
	return safesession.IPInfo{}
}, safesession.DB{
	Store: func(ID string, CreateTime time.Time) bool {
		s := safesession.Session{ID: ID, CreateTime: CreateTime}
		result := db.Create(&s)
		if result.Error == nil {
			return true
		}
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return false
		}
		panic(result.Error)
	},
	Delete: func(ID string) {
		s := safesession.Session{ID: ID}
		result := db.Delete(&s)
		if result.Error != nil {
			panic(result.Error)
		}
	},
	Exist: func(ID string) bool {
		s := safesession.Session{ID: ID}
		result := db.First(&s)
		if result.Error == nil && !s.CreateTime.IsZero() {
			return true
		}
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false
		}
		panic(result.Error)
	},
	Valid: func(UserName string, SessionID string) error {
		u := user{Id: UserName}
		result := db.First(&u)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return errors.New("没有这个用户id" + UserName)
			} else {
				panic(result.Error)
			}
		}

		if u.Session1 != SessionID && u.Session2 != SessionID {
			return safesession.LoginExpired
		}

		s := safesession.Session{ID: SessionID}
		result = db.First(&s)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				delete_session(&u, SessionID)
				return safesession.LoginExpired
			} else {
				panic(result.Error)
			}
		}

		if time.Now().After(s.CreateTime.Add(sessionAge)) {
			delete_session(&u, SessionID)
			return safesession.LoginExpired
		}
		return nil
	},
})

func delete_session(u *user, SessionID string) {
	if u.Session1 == SessionID {
		u.Session1 = ""
	}
	if u.Session2 == SessionID {
		u.Session2 = ""
	}
	if result := db.Save(u); result.Error != nil {
		panic(result.Error)
	}
}

func add_session(u *user, SessionID string) {
	if u.Session1 == "" {
		u.Session1 = SessionID
	} else if u.Session2 == "" {
		u.Session2 = SessionID
	} else {
		//如果登录了两台设备，就覆盖最先登录的设备
		u.Session1 = SessionID
	}
	if result := db.Save(u); result.Error != nil {
		panic(result.Error)
	}

}
