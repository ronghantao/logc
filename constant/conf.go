package constant

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type ConfInfo struct {
	Mode       string
	Port       int
	MaxConnCnt int
	Protocol   string
	ServerHost string
}

var Conf ConfInfo

func (c *ConfInfo) ParseConfFile(filepath string) error {
	fmt.Printf("使用配置文件'%s'初始化系统\n", filepath)
	if finfo, err := os.Stat(filepath); os.IsNotExist(err) || finfo.IsDir() {
		fmt.Println("配置文件'%s'不存在或者不是一个合法的文件", filepath)
		return err
	}
	f, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("读取配置文件出错，错误信息:", err.Error())
		return err
	}

	//开始使用配置文件初始化参数
	dec := json.NewDecoder(bufio.NewReader(f))
	var confMap map[string]interface{}
	if err := dec.Decode(&confMap); err != nil {
		//这里可能是一个非法的json
		fmt.Println("配置文件不是一个合法的json数据")
		return err
	}
	if val, ok := confMap["mode"]; ok {
		Conf.Mode = val.(string)
	}
	if val, ok := confMap["port"]; ok {
		floatVal := val.(float64)
		Conf.Port = int(floatVal)
	}
	if val, ok := confMap["server_host"]; ok {
		Conf.ServerHost = val.(string)
	} else {
		Conf.ServerHost = "127.0.0.1"
	}
	if Conf.IsServer() {
		if val, ok := confMap["max_conn_count"]; ok {
			floatVal := val.(float64)
			Conf.MaxConnCnt = int(floatVal)
		} else {
			Conf.MaxConnCnt = CONN_MAX_COUNT
		}

	}
	return nil
}

func (c ConfInfo) IsServer() bool {
	return c.Mode == STYPE_SERVER
}

func (c ConfInfo) IsClient() bool {
	return c.Mode == STYPE_CLIENT
}

func (c ConfInfo) String() string {
	s := "启动类型: " + c.Mode + "\n"
	if c.IsServer() {
		s += "启动端口：" + strconv.Itoa(c.Port) + "\n"
		s += "最大连接数：" + strconv.Itoa(c.MaxConnCnt) + "\n"
	}
	return s
}
