package dd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Config struct {
	AuthToken    string
	BarkId       string
	FloorId      int //1,普通商品 2,全球购保税 3,特殊订购自提 4,大件商品 5,厂家直供商品 6,特殊订购商品 7,失效商品
	DeliveryType int //1 急速达，2， 全程配送
	Longitude    string
	Latitude     string
	Deviceid     string
	Trackinfo    string
}

type DingdongSession struct {
	Conf               Config
	Address            Address            `json:"address"`
	Uid                string             `json:"uid"`
	Capacity           Capacity           `json:"capacity"`
	Channel            string             `json:"channel"` //0 => wechat  1 =>alipay
	SettleInfo         SettleInfo         `json:"settleInfo"`
	DeliveryInfoVO     DeliveryInfoVO     `json:"deliveryInfoVO"`
	SettleDeliveryInfo SettleDeliveryInfo `json:"settleDeliveryInfo"`
	GoodsList          []Goods            `json:"goods"`
	FloorInfo          FloorInfo          `json:"floorInfo"`
	StoreList          []Store            `json:"store"`
	OrderInfo          OrderInfo          `json:"orderInfo"`
	Client             *http.Client       `json:"client"`
	Cart               Cart               `json:"cart"`
}

func (s *DingdongSession) InitSession(conf Config) error {
	fmt.Println("########## 初始化 ##########")
	s.Client = &http.Client{Timeout: 60 * time.Second}
	s.Conf = conf

	err, addrList := s.GetAddress()
	if err != nil {
		return err
	}
	if len(addrList) == 0 {
		return errors.New("未查询到有效收货地址，请前往app添加或检查cookie是否正确！")
	}
	fmt.Println("########## 选择收货地址 ##########")
	for i, addr := range addrList {
		fmt.Printf("[%v] %s %s %s %s %s \n", i, addr.Name, addr.DistrictName, addr.ReceiverAddress, addr.DetailAddress, addr.Mobile)
	}
	var index int
	for true {
		fmt.Println("请输入地址序号（0, 1, 2...)：")
		stdin := bufio.NewReader(os.Stdin)
		_, err := fmt.Fscanln(stdin, &index)
		if err != nil {
			fmt.Printf("输入有误：%s!\n", err)
		} else if index >= len(addrList) {
			fmt.Println("输入有误：超过最大序号！")
		} else {
			break
		}
	}
	s.Address = addrList[index]

	fmt.Println("########## 选择支付方式 ##########")
	for true {
		fmt.Println("请输入支付方式序号（0：微信 1：支付宝)：")
		stdin := bufio.NewReader(os.Stdin)
		_, err := fmt.Fscanln(stdin, &index)
		if err != nil {
			fmt.Printf("输入有误：%s!\n", err)
		} else if index == 0 {
			s.Channel = "wechat"
			break
		} else if index == 1 {
			s.Channel = "alipay"
			break
		} else {
			fmt.Println("输入有误：序号无效！")
		}
	}
	return nil
}

func (s *DingdongSession) NewRequest(method, url string, dataStr []byte) *http.Request {

	var body io.Reader = nil
	if dataStr != nil {
		body = bytes.NewReader(dataStr)
	}
	req, _ := http.NewRequest(method, url, body)

	req.Header.Set("Host", "api-sams.walmartmobile.cn")
	req.Header.Set("content-type", "application/json;charset=UTF-8")
	//req.Header.Set("accept", "*/*")
	req.Header.Set("auth-token", s.Conf.AuthToken)
	req.Header.Set("longitude", s.Conf.Longitude)
	req.Header.Set("latitude", s.Conf.Latitude)
	req.Header.Set("device-id", s.Conf.Deviceid)
	req.Header.Set("app-version", "5.0.47.0")
	req.Header.Set("device-type", "ios")
	req.Header.Set("Accept-Language", "zh-Hans-CN;q=1")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("apptype", "ios")
	req.Header.Set("device-name", "iPhone14,5")
	req.Header.Set("device-os-version", "15.4.1")
	req.Header.Set("User-Agent", "SamClub/5.0.47 (iPhone; iOS 15.4.1; Scale/3.00)")
	req.Header.Set("system-language", "CN")

	return req
}
