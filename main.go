package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sevlyar/go-daemon"
)

//go:embed static/* templates/*
var embeddedFiles embed.FS

// 定义全局变量来存储 Basic Auth 账号密码
var username string
var password string

type CrontabEntry struct {
	Command  string `form:"command" binding:"required"`
	Schedule string `form:"schedule" binding:"required"`
}

func basicAuth() gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{
		username: password, // 修改为你的用户名和密码
	})
}

func getCrontab() ([]string, error) {
	cmd := exec.Command("crontab", "-l")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	// log.Printf(" crontab Output: %s", lines)
	return lines, nil
}

func setCrontab(lines []string) error {
	tmpFile := "/tmp/cron.tmp"
	newCron := strings.Join(lines, "\n")

	err := os.WriteFile(tmpFile, []byte(newCron+"\n"), 0644)

	if err != nil {
		return err
	}

	cmd := exec.Command("crontab", tmpFile)
	if err := cmd.Run(); err != nil {
		return err
	}
	err = os.Remove(tmpFile)
	return err
}
func split(input string, sep string) []string {
	parts := strings.Fields(input) // 自动去除多余空格并拆分
	if len(parts) == 0 {
		return []string{""} // 避免空切片
	}
	return parts
}

// join 用于将字符串切片连接为单一字符串
func join(input []string, sep string) string {
	return strings.Join(input, sep)
}

func runServer(port string) {

	funcMap := template.FuncMap{
		"split": split,
		"join":  join,
	}
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(basicAuth())
	r.StaticFS("/static", http.FS(embeddedFiles))
	// r.LoadHTMLGlob("templates/*")
	r.SetHTMLTemplate(template.Must(template.New("").Funcs(funcMap).ParseFS(embeddedFiles, "templates/*")))

	r.GET("/", func(c *gin.Context) {
		entries, err := getCrontab()
		if err != nil {
			c.String(http.StatusInternalServerError, "Error getting crontab: %v", err)
			return
		}
		if len(entries) == 0 {
			entries = nil
		}
		c.HTML(http.StatusOK, "index.html", gin.H{"entries": entries})
	})

	r.POST("/add", func(c *gin.Context) {
		var entry CrontabEntry
		if err := c.ShouldBindWith(&entry, binding.Form); err != nil {
			c.String(http.StatusBadRequest, "Invalid input")
			return
		}
		entries, err := getCrontab()
		if err != nil {
			c.String(http.StatusInternalServerError, "Error getting crontab")
			return
		}
		entries = append(entries, fmt.Sprintf("%s %s", entry.Schedule, entry.Command))
		if err := setCrontab(entries); err != nil {
			c.String(http.StatusInternalServerError, "Failed to set crontab")
			return
		}
		c.Redirect(http.StatusSeeOther, "/")
	})

	r.POST("/delete", func(c *gin.Context) {
		index := c.PostForm("index")
		entries, err := getCrontab()
		if err != nil {
			c.String(http.StatusInternalServerError, "Error getting crontab")
			return
		}
		idx := -1
		for i, _ := range entries {
			if fmt.Sprintf("%d", i) == index {
				idx = i
				break
			}
		}
		if idx != -1 {
			entries = append(entries[:idx], entries[idx+1:]...)
			if err := setCrontab(entries); err != nil {
				c.String(http.StatusInternalServerError, "Failed to set crontab")
				return
			}
		}
		c.Redirect(http.StatusSeeOther, "/")
	})

	r.Run(":" + port)
}

func main() {
	port := flag.String("p", "10010", "Web 监听端口")
	user := flag.String("u", "admin", "Basic Auth 用户名")
	pass := flag.String("pwd", "password", "Basic Auth 密码")
	daemonMode := flag.Bool("d", false, "是否以守护进程模式运行")
	flag.Parse()

	username = *user
	password = *pass

	// 守护进程模式
	if *daemonMode {
		cntxt := &daemon.Context{
			PidFileName: "crontab_manager.pid",
			PidFilePerm: 0644,
			LogFileName: "crontab_manager.log",
			LogFilePerm: 0640,
			WorkDir:     "./",
			Umask:       027,
			Args:        []string{"[CrontabManager]"},
		}

		d, err := cntxt.Reborn()
		if err != nil {
			fmt.Println("无法启动守护进程:", err)
			os.Exit(1)
		}
		if d != nil {
			return
		}
		defer cntxt.Release()
		fmt.Println("守护进程启动成功！")
		runServer(*port)
	} else {
		// 直接前台运行
		runServer(*port)
	}
}
