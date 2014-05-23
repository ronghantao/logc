package main

import (
	"flag"
	"fmt"
	"github.com/ronghantao/logc/client"
	"github.com/ronghantao/logc/constant"
	"github.com/ronghantao/logc/server"
	//	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

//Init 初始化参数列表
func init() {
	constant.Conf.Protocol = "tcp"
	var confFile string
	var baseDir string
	//初始化本地基础路径
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		baseDir = ""
	} else {
		baseDir = dir
	}

	fmt.Println("基础路径：", baseDir)

	var cmdMode string
	var cmdPort int
	flag.StringVar(&confFile, "conf", "", "配置文件地址")
	flag.StringVar(&cmdMode, "type", "", "启动类型，client=客户端；server=服务器。默认client")
	flag.IntVar(&cmdPort, "port", 0, "绑定端口，当type=server时有效")
	flag.Parse()
	tailArgs := flag.Args()
	if size := len(tailArgs); size > 0 {
		fmt.Println("以下参数无法解析，我们将忽略这些参数：", tailArgs)
	}
	//如果配置了-conf参数，从conf文件中读取配置
	if confFile == "" {
		//判断默认配置文件是否存在
		defConfPath := baseDir + string(os.PathSeparator) + constant.LOGC_CONF_DEFAULT
		if finfo, err := os.Stat(defConfPath); os.IsNotExist(err) || finfo.IsDir() {
			//如果默认的配置文件不存在，或者是一个目录，就不再分析配置文件
			fmt.Printf("默认配置文件'%s'不存在或者不是一个合法的文件\n", defConfPath)
			fmt.Println("跳过默认配置文件")
		} else {
			confFile = defConfPath
			fmt.Printf("发现默认配置文件'%s'\n", confFile)
		}
	}
	if confFile != "" {
		//添加了-conf参数或者默认配置文件存在，再次验证文件是否存在，并且是合法的文件
		err := constant.Conf.ParseConfFile(confFile)
		if nil != err {
			fmt.Println(err.Error())
			os.Exit(constant.ERROR_CONF_FILE_PARSER)
		}
	}
	//再从命令行读取数据
	if cmdMode != "" {
		if constant.Conf.Mode != "" {
			fmt.Println("命令行发现'type'参数，这将覆盖配置文件中的设置.")
		}
		constant.Conf.Mode = cmdMode
	}
	if cmdPort != 0 {
		if constant.Conf.Port != 0 {
			fmt.Println("命令行发现'port'参数，这将覆盖配置文件中的设置.")
		}
		constant.Conf.Port = cmdPort
	}

	if constant.Conf.Mode == constant.STYPE_SERVER && constant.Conf.Port == 0 {
		fmt.Printf("启动类型为'%s'时，必须提供port参数.\n", constant.Conf.Mode)
		os.Exit(constant.ERROR_PARAM_PORT_ERROR)
	} else if constant.Conf.Mode == constant.STYPE_CLIENT && constant.Conf.Port > 0 {
		fmt.Printf("启动类型为'%s'时，port参数无效，我们将忽略port参数.\n", constant.Conf.Mode)
	}
	fmt.Println(constant.Conf)
}

func main() {
	go server.StartServer()
	time.Sleep(time.Second * 10)
	client.C.SendEcho()
}
