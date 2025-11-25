package netdisk

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func Route(s *gin.Engine) {
	s.GET("/", func(ctx *gin.Context) {
		ctx.File("./netdisk/index.html")
	})
	s.GET("/user.html", func(ctx *gin.Context) {
		ctx.File("./netdisk/user.html")
	})
	s.GET("/register.html", func(ctx *gin.Context) {
		ctx.File("./netdisk/register.html")
	})
	s.GET("/register_result.html", func(ctx *gin.Context) {
		ctx.File("./netdisk/register_result.html")
	})
	s.GET("/login.html", func(ctx *gin.Context) {
		ctx.File("./netdisk/login.html")
	})
	s.GET("/upload.html", func(ctx *gin.Context) {
		ctx.File("./netdisk/upload.html")
	})
	s.Static("/assets", "./netdisk/assets")
	s.POST("/register", func(ctx *gin.Context) {
		type data struct {
			Name     string `form:"name" binging:"required"`
			Password string `form:"password1" binging:"required"`
		}
		f := &data{}
		err := ctx.ShouldBind(f)
		if err != nil { //TODO:更好的处理
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		id := new_user(f.Name, f.Password)
		login(ctx, id, f.Password)
		ctx.Redirect(303, "/register_result.html")
	})
	s.POST("/login", func(ctx *gin.Context) {
		tmp, err := ctx.Request.Cookie("session")
		if err != http.ErrNoCookie {
			_, err, _ := sessionControl.CheckLogined(ctx.ClientIP(), ctx.Request.UserAgent(), tmp)
			if err != nil {
				tmp.MaxAge = -1
				http.SetCookie(ctx.Writer, tmp)
				ctx.String(http.StatusUnauthorized, err.Error())
				return
			}
			ctx.Redirect(303, "/")
			return
		}
		type data struct {
			ID       string `form:"id" binging:"required"`
			Password string `form:"password1" binging:"required"`
		}
		f := &data{}
		err = ctx.ShouldBind(f)
		if err != nil { //TODO:更好的处理
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}

		ok, _ := login(ctx, f.ID, f.Password)
		if !ok {
			ctx.String(http.StatusUnauthorized, "用户名或密码错误")
			return
		}
		ctx.Redirect(303, "/")
	})
	s.POST("/list", func(ctx *gin.Context) {
		id, shouldReturn := check_login(ctx)
		if shouldReturn {
			return
		}

		filedb := new_fileDb()
		s := filedb.List(id)
		ctx.JSON(200, s)
	})
	s.POST("/upload", func(ctx *gin.Context) {
		id, shouldReturn := check_login(ctx)
		if shouldReturn {
			return
		}

		f, err := ctx.FormFile("file")
		if err != nil { //TODO:更好的处理
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}

		// 打开文件
		src, err := f.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer src.Close()

		filedb := new_fileDb()

		path := ctx.Query("path")
		if path == "" {
			ctx.String(http.StatusBadRequest, "没有 path")
			return
		}

		if !filedb.Set(id, path, src) {
			ctx.String(http.StatusForbidden, "上传文件太大，超出您的可用剩余空间上限！")
		}
	})
	s.POST("/setname", func(ctx *gin.Context) {
		id, shouldReturn := check_login(ctx)
		if shouldReturn {
			return
		}
		update(ctx, "name", id)
		ctx.Redirect(303, "/user.html")
	})
	s.POST("/set_password", func(ctx *gin.Context) {
		id, shouldReturn := check_login(ctx)
		if shouldReturn {
			return
		}
		update(ctx, "password", id)
		ctx.Redirect(303, "/user.html")
	})
	s.POST("/stat", func(ctx *gin.Context) {
		id, shouldReturn := check_login(ctx)
		if shouldReturn {
			return
		}
		u := user{}
		result := db.Model(&user{}).Select("limit", "used").Where("id = ?", id).First(&u)
		if result.Error != nil {
			panic(result.Error)
		}
		ctx.JSON(200, u)
	})
	s.GET("/download/*filepath", func(ctx *gin.Context) {
		id, path, shouldReturn := getIdAndPath(ctx)
		if shouldReturn {
			return
		}
		filedb := new_fileDb()

		r := filedb.Get(id, path)
		if r == nil {
			ctx.String(404, "")
			return
		}

		io.Copy(ctx.Writer, r)
	})
	s.GET("/delete/*filepath", func(ctx *gin.Context) {
		id, path, shouldReturn := getIdAndPath(ctx)
		if shouldReturn {
			return
		}
		filedb := new_fileDb()

		r := filedb.Delete(id, path)
		if !r {
			ctx.String(404, "")
			return
		}

		ctx.Redirect(303, "/")
	})
}

func getIdAndPath(ctx *gin.Context) (id string, path string, shouldreutrn bool) {
	id, shouldReturn := check_login(ctx)
	if shouldReturn {
		return "", "", true
	}
	path = ctx.Param("filepath")
	if path == "" {
		ctx.String(404, "空 filepath")
		return "", "", true
	}
	path, err := url.PathUnescape(path)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Invalid file path encoding: %v", err)
		return "", "", true
	}
	path = strings.TrimSuffix(path, "/")
	path = strings.TrimPrefix(path, "/")
	return id, path, false
}

func update(ctx *gin.Context, field string, userid string) {
	if field == "password" {
		field = "password1"
	}
	v := ctx.PostForm(field)
	if v == "" {
		//TODO:更好的处理
		ctx.String(http.StatusBadRequest, "")
		return
	}
	if field == "password1" {
		field = "password"
	}

	if field == "password" {
		v = hash(v)
	}

	result := db.Model(&user{}).Where("id = ?", userid).Update(field, v)
	if result.Error != nil {
		panic(result.Error)
	}

	if field == "name" {
		add_login_cookie(ctx, "", v)
	}
}

// check_login 检查是否登录，如果是返回用户ID
func check_login(ctx *gin.Context) (string, bool) {
	tmp, err := ctx.Request.Cookie("session")
	if err == http.ErrNoCookie {
		ctx.String(http.StatusUnauthorized, "未登录")
		return "", true
	}
	_, err, se := sessionControl.CheckLogined(ctx.ClientIP(), ctx.Request.UserAgent(), tmp)
	if err != nil {
		http.SetCookie(ctx.Writer, &http.Cookie{
			Name:   "session",
			MaxAge: -1,
		})
		http.SetCookie(ctx.Writer, &http.Cookie{
			Name:   "Name",
			MaxAge: -1,
		})
		ctx.String(http.StatusUnauthorized, err.Error())
		return "", true
	}
	return se.Name, false
}

func add_login_cookie(ctx *gin.Context, id string, Name string) {
	if id != "" {
		http.SetCookie(ctx.Writer, &http.Cookie{
			Name:   "ID",
			Value:  id,
			Secure: true,
		})
	}
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:   "Name",
		Value:  Name,
		Secure: true,
	})
}
