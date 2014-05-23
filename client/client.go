package client

import (
	"fmt"
	"github.com/ronghantao/logc/constant"
	"log"
	"net"
	"strconv"
)

type Client struct {
	Conn net.Conn
}

var C Client

func (c *Client) Read(buffer []byte) bool {
	_, err := c.Conn.Read(buffer)
	if nil != err {
		log.Println(err.Error())
		return false
	}
	log.Println("Read from connection: ", buffer)
	return true
}

func (c *Client) SendEcho() error {
	log.Println("begin send echo...")
	conn, err := net.Dial(constant.Conf.Mode, constant.Conf.ServerHost+":"+strconv.Itoa(constant.Conf.Port))
	if nil != err {
		log.Println("create connection to ", constant.Conf.ServerHost+":"+strconv.Itoa(constant.Conf.Port), " failed")
		return err
	}
	defer conn.Close()
	//发送echo命令
	fmt.Fprintf(conn, "hello")
	log.Println("end send echo...")
	return nil
}
