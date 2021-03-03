package main

import (
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	// Universal markup builders.
	menu = &tb.ReplyMarkup{ResizeReplyKeyboard: true}

	// Reply buttons.
	btnadd    = menu.Text("⊕ 추가")
	btnremove = menu.Text("Θ 제거")
	btnview   = menu.Text("Θ 보기")
	btntoday   = menu.Text("Θ 정보")
)
var (
	// Universal markup builders.
	off = &tb.ReplyMarkup{ResizeReplyKeyboard: true}

	// Reply buttons.
	offbutton = off.Text("● 끄기")
)

var bp *tb.Bot
var user = &tb.User{ID: USER_ID}
var poller = &tb.LongPoller{Timeout: 15 * time.Second}
var spamProtected = tb.NewMiddlewarePoller(poller, func(upd *tb.Update) bool {
	if upd.Message == nil {
		return true
	}

	if strings.Contains(upd.Message.Text, "추가") {
		return true
	} else if strings.Contains(upd.Message.Text, "제거") {
		return true
	} else if strings.Contains(upd.Message.Text, "끄기") {
		return true
	} else if strings.Contains(upd.Message.Text, "보기") {
		return true
	}

	if strings.Contains(upd.Message.Text, "정보") == true {
		var nyat string
		if today == true{
			nyat = "활성화"
		}
		if today == false{
			nyat = "비활성화"
		}
		_, _ = bp.Send(user,  nyat + "상태입니다.")
		_, _ = bp.Send(user,  "활성화된 고루틴은" + strconv.Itoa(runtime.NumGoroutine()) + "개 입니다.")

	}

	dd := strings.Split(upd.Message.Text, " ")
	for _, a := range dd {
		_, err := strconv.Atoi(a)
		if err != nil {
			return true
		}
	}
	// date 파일이 비었는지 확인 후 비었으면 n_Date 파일에서 내용 끌어옴 , 최적화 완료
	datc, _ := ioutil.ReadFile("/date.txt")
	if string(datc) == "" {
		err := ioutil.WriteFile("/date.txt", []byte(upd.Message.Text), os.FileMode(644))
		if err != nil {
			_, _ = bp.Send(user, "추가하는데 오류가 발생하였습니다.")
		}
	} else {
		err := ioutil.WriteFile("/n_date.txt", []byte(upd.Message.Text), os.FileMode(644))
		if err != nil {
			_, _ = bp.Send(user, "추가하는데 오류가 발생하였습니다.")
		}
	}
	_, _ = bp.Send(user, upd.Message.Text+"일 추가 완료")
	return true
})
var today = false
var today_nat = false
func main() {
	_, err := os.Open("/off_m.txt")
	if err != nil {
		_, _ = os.Create("/off_m.txt")
	}
	_, err = os.Open("/off_hm.txt")
	if err != nil {
		_, _ = os.Create("/off_hm.txt")
	}
	_, err = os.Open("/date.txt")
	if err != nil {
		_, _ = os.Create("/date.txt")
	}
	_, err = os.Open("/n_date.txt")
	if err != nil {
		_, _ = os.Create("/n_date.txt")
	}
	_, err = os.Open("/month.txt")
	if err != nil {
		_, _ = os.Create("/month.txt")
	}
	_, err = os.Open("/last_date_m.txt")
	if err != nil {
		_, _ = os.Create("/last_date_m.txt")
	}
	bot, err := tb.NewBot(tb.Settings{
		Token:  "TOKEN:TOKEN",
		Poller: spamProtected,
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	bp = bot
	menu.Reply(
		menu.Row(btnadd),
		menu.Row(btnremove),
		menu.Row(btnview),
		menu.Row(btntoday),
	)
	off.Reply(
		off.Row(offbutton),
	)
	_, _ = bp.Send(user, "", menu)
	// On reply button pressed (message)
	bp.Handle(&btnadd, func(m *tb.Message) {
		_, _ = bp.Send(m.Sender, "쉬는 날짜를 입력해주세요")
	})
	bp.Handle(&btnremove, func(m *tb.Message) {
		_ = ioutil.WriteFile("/n_date.txt", []byte(""), os.FileMode(644))
		_ = ioutil.WriteFile("/date.txt", []byte(""), os.FileMode(644))
		_, _ = bp.Send(m.Sender, "제거되었음")
	})
	bp.Handle(&btnview, func(m *tb.Message) {
		data, _ := ioutil.ReadFile("/n_date.txt")
		data2, _ := ioutil.ReadFile("/date.txt")
		_, _ = bp.Send(m.Sender, "이번달 휴일 :"+string(data2)+"\n\n 다음달 휴일 :"+string(data))
	})
	bp.Handle(&offbutton, func(m *tb.Message) {
		var timea = time.Now()
		if timea.Hour() >= 11 {
			err = ioutil.WriteFile("/off_hm.txt", []byte("1"), os.FileMode(644))
		} else {
			err = ioutil.WriteFile("/off_m.txt", []byte("1"), os.FileMode(644))
		}
		if err != nil {
			_, _ = bp.Send(user, "알람을 끄는데 실패하였습니다.", menu)
		}
		_, _ = bp.Send(m.Sender, "꺼짐!", menu)
	})
	go NeverExit(mezamashi)
	go NeverExit(checkmonthdate)
	bp.Start()
}
func mezamashi() {
	time.Sleep(5*time.Second)
	for {
		timea := time.Now()
		if today == true {
			for {
				if timea.Hour() == 7 {
					offm, err := ioutil.ReadFile("/off_m.txt")
					if err != nil {
						break
					}
					if timea.Minute() >= 0 {
						if string(offm) == "1" {
							today = false
							today_nat = true
							break
						}
						_, _ = bp.Send(user, "일어나세요!", off)
					}
				} else {
					break
				}
				time.Sleep(500 * time.Millisecond)
			}
		}
		if today_nat == true{
			for {
				if timea.Hour() >= 11 {
					offc, err := ioutil.ReadFile("/off_hm.txt")
					if err != nil {
						break
					}
					if timea.Minute() >= 30 {
						if string(offc) == "1" {
							today_nat = false
							break
						}
						_, _ = bp.Send(user, "일어나세요!", off)
					}
				} else {
					break
				}
				time.Sleep(500 * time.Millisecond)
			}
		}
		time.Sleep(29500 * time.Millisecond)
	}
}



func checkmonthdate() {
	for {
		timea := time.Now()
		montha := strconv.Itoa(int(timea.Month()))
		date := strconv.Itoa(timea.Day())

		lastMonth, err := ioutil.ReadFile("/month.txt")
		if err != nil {
			panic(err)
		}
		lastDate, err := ioutil.ReadFile("/last_date_m.txt")
		if err != nil {
			panic(err)
		}
		datc, err := ioutil.ReadFile("/date.txt")
		if err != nil {
			panic(err)
		}
		// 월이 맞지 않으면
		if montha != string(lastMonth) {
			data, err := ioutil.ReadFile("/n_date.txt")
			if err != nil {
				_, err = bp.Send(user, "다음달 일을 불러오는데 실패하였습니다.", menu)
				panic(err)
			}
			err = ioutil.WriteFile("/month.txt", []byte(montha), os.FileMode(644))
			if err != nil {
				_, err = bp.Send(user, "월을 덧쓰는데 실패하였습니다.", menu)
				panic(err)
			}
			err = ioutil.WriteFile("/date.txt", data, os.FileMode(644))
			if err != nil {
				_, err = bp.Send(user, "날을 덧쓰는데 실패하였습니다.", menu)
				panic(err)
			}
			err = ioutil.WriteFile("/n_date.txt", []byte(""), os.FileMode(644))
			if err != nil {
				_, err = bp.Send(user, "다음달 예정을 초기화하는데 실패하였습니다.", menu)
				panic(err)
			}
			_, _ = bp.Send(user, "월이 바뀌었습니다", menu)
		}
		// 일이 맞지 않으면
		if date != string(lastDate) {
			err = ioutil.WriteFile("/off_m.txt", []byte("0"), os.FileMode(644))
			if err != nil {
				_, err = bp.Send(user, "아침 알람을 초기화하는데 실패하였습니다.", menu)
				panic(err)
			}
			err = ioutil.WriteFile("/off_hm.txt", []byte("0"), os.FileMode(644))
			if err != nil {
				_, err = bp.Send(user, "점심 알람을 초기화하는데 실패하였습니다.", menu)
				panic(err)
			}
			err := ioutil.WriteFile("/last_date_m.txt", []byte(date), os.FileMode(644))
			if err != nil {
				_, err = bp.Send(user, "날짜를 덧쓰는데 실패하였습니다.", menu)
				panic(err)
			}

		}
		// 오늘이 알람 울릴 날인가
		dateS := strings.Split(string(datc), " ")
		cache := true
		for _, dates := range dateS {
			if strconv.Itoa(timea.Day()) == dates {
				cache = false
			}
		}
		today = cache
		time.Sleep(6 * time.Hour)
	}
}

func NeverExit(f func()) {
	defer func() { if v := recover(); v != nil {
		log.Println(v)
		go NeverExit(f)
		 }
	}()
	f()
}
