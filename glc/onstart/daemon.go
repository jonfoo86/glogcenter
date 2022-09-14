package onstart

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"time"
)

func init() {

	// 仅支持linux
	if runtime.GOOS != "linux" {
		return
	}

	// pid 目录、文件
	pidpath := "."
	pidfile := "glc.pid"
	u, err := user.Current()
	if err == nil {
		pidpath = filepath.Join(u.HomeDir, ".glogcenter")
		os.MkdirAll(pidpath, 0766)
	}
	pidpathfile := filepath.Join(pidpath, pidfile)

	// 操作系统是linux时，支持以命令行参数【-d】后台方式启动，【stop】停止
	daemon := false
	stop := false
	restart := false
	for index, arg := range os.Args {
		if index == 0 {
			continue
		}
		if arg == "-d" {
			daemon = true
		}
		if arg == "stop" {
			stop = true
		}
		if arg == "restart" {
			restart = true
		}
	}

	rs := checkPidFile(pidpathfile)
	if rs != nil {
		if stop || restart {
			// 退出/重启
			cmd := exec.Command("sh", "-c", "kill "+rs.Pid)
			cmd.Start()
		} else {
			// 禁止重复启动
			fmt.Printf("%s\n", rs.Pid)
			os.Exit(0)
		}
	}

	if stop {
		os.Exit(0)
	}

	if daemon {
		// cmd := exec.Command(os.Args[0], flag.Args()...)
		cmd := exec.Command(os.Args[0]) // 不再需要启动参数了

		err := cmd.Start()
		for i := 0; i < 60; i++ {
			if err != nil {
				time.Sleep(time.Duration(1) * time.Second) // 原进程没退出的话会导致启动失败，等待1秒后再试
			} else {
				break
			}
		}
		if err != nil {
			// 最多等1分钟，仍旧启动失败就放弃
			fmt.Printf("start %s failed, error: %v\n", os.Args[0], err)
			os.Exit(1)
		}

		fmt.Printf("%d\n", cmd.Process.Pid)
		os.Exit(0)
	} else {
		npid := fmt.Sprintf("%d", os.Getpid())
		nerr := savePid(pidpathfile, npid)
		if nerr != nil {
			cmd := exec.Command("sh", "-c", "kill "+npid)
			cmd.Start()
			os.Exit(1)
		}
	}

}
