package server

import (
	"errors"
	"fmt"
	"github.com/ronghantao/logc/constant"
	"log"
	"net"
	"strconv"
	"time"
)

type server struct {
	protocol   string        //协议，默认tcp
	port       int           //端口号
	maxConnCnt int           //最大连接数
	curConnCnt int           //当前连接数
	quitFlag   bool          //是否开始退出
	connChan   chan net.Conn //连接管道
}

var chConnChange = make(chan int)
var ser server

func StartServer() error {
	defer func() {
		log.Println("Server is closing...")
		log.Println("Server closed.")
	}()
	ser.protocol = constant.Conf.Protocol
	ser.port = constant.Conf.Port
	ser.maxConnCnt = constant.Conf.MaxConnCnt
	ser.curConnCnt = 0
	ser.quitFlag = false
	ser.connChan = make(chan net.Conn, ser.maxConnCnt)
	err := ser.listen()
	return err
}

func StopServer() {
	ser.quitFlag = true
}

func (s *server) listen() (err error) {
	listener, err := net.Listen(s.protocol, ":"+strconv.Itoa(s.port))
	if err != nil {
		fmt.Println("error listener: ", err.Error())
		return errors.New("create listener failed.")
	}
	defer listener.Close()

	defer func() {
		//这里等待所有的conn关闭
		if !s.quitFlag {
			s.quitFlag = true
		}
		ticker := time.NewTicker(time.Second)
		for _ = range ticker.C {
			if s.curConnCnt <= 0 {
				break
			}
			log.Printf("wait all process quit...\n")
		}
	}()

	fmt.Println("running ...")

	s.curConnCnt = 0

	go func() {
		for connChange := range chConnChange {
			s.curConnCnt += connChange
		}
	}()

	// for i := 0; i < s.maxConnCnt; i++ {
	// 	go func() {
	// 		for conn := range connChan {
	// 			chConnChange <- 1
	// 			EchoFunc(conn)
	// 			chConnChange <- -1
	// 		}
	// 	}()
	// }

	go func() {
		//监听connection
		for conn := range s.connChan {
			go EchoFunc(conn)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Error accept:", err.Error())
			return err
		}
		if s.quitFlag {
			conn.Close()
			return nil
		}
		s.connChan <- conn
		chConnChange <- 1
	}
	return nil
}

func EchoFunc(conn net.Conn) {
	defer func() { chConnChange <- -1 }()
	defer conn.Close()
	log.Println("echo func.")
}
