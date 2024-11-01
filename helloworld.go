package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// CouponInfo 结构体用于表示 dataObj 中的每个优惠券信息
type CouponInfo struct {
	CouponID           int    `json:"couponId"`
	CouponGUID         string `json:"couponGuid"`
	Type               int    `json:"type"`
	TypeName           string `json:"typeName"`
	LineType           int    `json:"lineType"`
	CouponName         string `json:"couponName"`
	PayFee             int    `json:"payFee"`
	BeginDate          string `json:"beginDate"`
	ValidDate          string `json:"validDate"`
	ValidDays          int    `json:"validDays"`
	Description        string `json:"description"`
	CouponState        int    `json:"couponState"`
	GetState           int    `json:"getState"`
	ButtonState        int    `json:"buttonState"`
	CouponContent      string `json:"couponContent"`
	LimitAmountContent string `json:"limitAmountContent"`
	Link               string `json:"link"`
	ApplyType          int    `json:"applyType"`
	IsSuperpose        bool   `json:"isSuperpose"`
}

var globalMap = map[string][]string{
	"9597A0EE8282571A2379FB006F5E4AE6": []string{"15059546210", "小玉"},
	"D66096A73F66A082779A5A9CDB1186F2": []string{"18960432226", "小何"},
	"AC3FCFBA3C0753DF61838A95088F2A98": []string{"13559506500", "小何"},
	"89EA9514F00FA7A28A8EFD1397740B9D": []string{"13505918710", "小花"},
	"D2A8E4860A74B7B46A2B3B7855802FEC": []string{"15280286253", "小林"},
	"30BC02C5F03D869EDA4F989287435AA2": []string{"15959534510", "小苏"},
	"A621754E276E560F838891D0760E737E": []string{"15260500473", "王銮坚"},
	"7B28820E41C35B78E8A7A464FE489128": []string{"18859909987", "谢玮琼"},
	"0CE92D8786B1EF68F5D66B85F53246A0": []string{"13615966669", "李燕玲"},
	"F6FBC8920964BEFDB1529F2162C3EE84": []string{"13505039597", "9597"},
	"F20FF965DD8758555F1248EF0CC2D058": []string{"15375775673", "5673"},
}

// 按需求抢券，未配置放空为全抢
var mRobCouponList = map[string][]int32{
	"9597A0EE8282571A2379FB006F5E4AE6": []int32{5, 10, 20, 50, 100},      // 林燕玉
	"D66096A73F66A082779A5A9CDB1186F2": []int32{5, 10, 20},               // 小何
	"AC3FCFBA3C0753DF61838A95088F2A98": []int32{5, 10, 20},               // 小何
	"89EA9514F00FA7A28A8EFD1397740B9D": []int32{5, 10, 20, 50},           // 苏添花
	"D2A8E4860A74B7B46A2B3B7855802FEC": []int32{5, 10},                   // 林望琛
	"30BC02C5F03D869EDA4F989287435AA2": []int32{5, 10},                   // 苏丽娇
	"A621754E276E560F838891D0760E737E": []int32{5, 10},                   // 王銮坚
	"F6FBC8920964BEFDB1529F2162C3EE84": []int32{5, 10, 20, 50, 100, 200}, // 锦裕
	"F20FF965DD8758555F1248EF0CC2D058": []int32{5, 10, 20, 50, 100, 200}, // 锦裕
	//"7B28820E41C35B78E8A7A464FE489128": []int32{5, 10, 100}, // 谢玮琼
	//"0CE92D8786B1EF68F5D66B85F53246A0": []int32{},         // 李燕玲
}

const mTopicId int64 = 149857372

func main() {
	// 获取当前时间
	now := time.Now()

	// 计算距离下一个上午十点和下午三点的时间间隔
	nextMorning := time.Date(now.Year(), now.Month(), now.Day(), 9, 59, 55, 0, now.Location())
	if now.After(nextMorning) {
		nextMorning = nextMorning.Add(24 * time.Hour) // 下一个上午十点
	}

	nextAfternoon := time.Date(now.Year(), now.Month(), now.Day(), 14, 59, 55, 0, now.Location())
	if now.After(nextAfternoon) {
		nextAfternoon = nextAfternoon.Add(24 * time.Hour) // 下一个下午三点
	}

	for k, _ := range mRobCouponList {

		// 启动定时器，每天上午十点执行一次
		morningTimer := time.NewTimer(nextMorning.Sub(now))
		go func(userId string) {
			for {
				<-morningTimer.C
				fmt.Println("It's 10 o'clock in the morning!")
				// 执行你的任务
				startSeizeCoupons(userId)

				// 重新计算下一个上午十点的时间
				nextMorning = nextMorning.Add(24 * time.Hour)
				morningTimer.Reset(nextMorning.Sub(time.Now()))
			}
		}(k)

		// 启动定时器，每天下午三点执行一次
		afternoonTimer := time.NewTimer(nextAfternoon.Sub(now))
		go func(userId string) {
			for {
				<-afternoonTimer.C
				fmt.Println("It's 3 o'clock in the afternoon!")
				// 执行你的任务
				startSeizeCoupons(userId)

				// 重新计算下一个下午三点的时间
				nextAfternoon = nextAfternoon.Add(24 * time.Hour)
				afternoonTimer.Reset(nextAfternoon.Sub(time.Now()))
			}
		}(k)

	}

	// 等待程序结束
	select {}
}

func startSeizeCoupons(userId string) {
	// 准备要传递的参数
	params := map[string]string{
		"encryptUserId": userId,
	}

	// 将参数转换为 JSON 格式
	jsonParams, err := json.Marshal(params)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", "https://liveauth.vzan.com/api/v1/login/get_wx_token", bytes.NewBuffer(jsonParams))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 设置 Header 头部
	setUrlHeader(req, "", userId)

	// 发送请求
	resp := sendReq(req)
	defer resp.Body.Close()

	// 读取回复消息的内容
	data := getDataObj(resp)

	var dataObjMap map[string]interface{}
	dataObjMap = data["dataObj"].(map[string]interface{})

	// token
	token := dataObjMap["token"]

	// 创建一个每秒钟触发一次的定时器
	ticker := time.NewTicker(100 * time.Millisecond)

	// 定义一个计数器，用于记录经过的时间
	counter := 0

	// 在一个无限循环中，等待定时器的触发事件
	for range ticker.C {
		// 执行你的任务
		fmt.Println("Task executed at", time.Now())
		toGetCoupon(fmt.Sprintf("%v", token), userId)

		// 每次任务执行后，增加计数器
		counter++

		// 如果经过了一分钟，停止定时器并退出循环
		if counter >= 90 {
			ticker.Stop()
			break
		}
	}
}

// 请求抢券
func toGetCoupon(token string, userId string) {

	iphone := globalMap[userId][0]
	name := globalMap[userId][1]

	// 抢的券码
	robList, ok := mRobCouponList[userId]
	if !ok {
		return
	}

	robMap := make(map[int32]int32)
	for _, v := range robList {
		couponMoney := v * 100
		robMap[couponMoney] = 1
	}

	//--------------------------------------
	//创建消费券列表请求
	couponListUrl := fmt.Sprintf("https://live-marketapi.vzan.com/api/v1/coupon/getmenucouponlist?topicId=%d", mTopicId)
	reqConsume, err := http.NewRequest("GET", couponListUrl, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	setUrlHeader(reqConsume, token, userId)
	respConsume := sendReq(reqConsume)
	defer respConsume.Body.Close()
	// 回复内容
	data := getDataObj(respConsume)
	if data == nil {
		return
	}
	// 有消费券的情况下，遍历消费券请求领取
	var coupons []interface{}
	if len(data["dataObj"].([]interface{})) > 0 {
		coupons = data["dataObj"].([]interface{})
	}

	curTime := time.Now().Format("2006-01-02 15:04:05")

	for _, v := range coupons {
		couponItem, ok := v.(map[string]any)
		if !ok {
			fmt.Println("Failed to convert item to CouponInfo")
			fmt.Printf("coupon list value: %v\n", v)
			continue
		}

		// 只抢有用券
		if len(robMap) > 0 {
			// 处理 couponMoney 断言并确保其为 int32 类型
			couponMoney, ok := couponItem["couponMoney"].(float64) // 服务器返回可能是 float64
			if !ok {
				fmt.Println("Failed to get couponMoney")
				continue
			}

			// 类型转换为 int32 与 robMap 匹配
			if _, ok := robMap[int32(couponMoney)]; !ok {
				continue
			}
		}

		//fmt.Printf("coupon item: %v \n", v)
		couponId := couponItem["couponId"]
		fmt.Println("===========================couponId:", couponId, " curTime:", curTime)
		params1 := map[string]any{
			"couponId":      couponId,
			"getConditions": 0,
			"isFromMenu":    true,
			"sourceId":      mTopicId,
			"sourceType":    1,
			"topicId":       mTopicId,
			"userName":      name,
			"userPhone":     iphone,
			"zbId":          629144760,
		}

		// 将参数转换为 JSON 格式
		jsonParams1, err := json.Marshal(params1)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		reqCoupon, err := http.NewRequest("POST", "https://live-marketapi.vzan.com/api/v1/coupon/GetCoupon", bytes.NewBuffer(jsonParams1))
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		setUrlHeader(reqCoupon, fmt.Sprintf("%v", token), userId)
		respCoupon := sendReq(reqCoupon)
		if respCoupon == nil {
			continue
		}
		// 回复内容
		data = getDataObj(respCoupon)
		fmt.Println("---------success to get coupon:", data)
	}
}

// 设置请求头
func setUrlHeader(req *http.Request, token string, userId string) {
	// 设置 Header 头部
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Buid", userId)

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
}

// 发送请求
func sendReq(req *http.Request) *http.Response {
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	return resp
}

// 获取基础信息
func getDataObj(resp *http.Response) map[string]interface{} {

	// 读取回复消息的内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	// 打印回复消息的内容
	//fmt.Println("Response:", string(body))

	return data
}
