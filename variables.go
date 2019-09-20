package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"
)

//包级变量管理器
//统一init函数
//避免golang同一包下多个init函数按照文件字母序调用出现问题

var (
	now         = time.Now()
	nextWorkDay time.Time
	params      *cmd //命令行参数保存
	//testFlag       = false     //使用这个标志在单元测试阶段关闭出发init函数
	lg             = lgVerbose //日志函数定义,默认采用罗嗦模式
	templateFile   *os.File    //模板文件
	jiraHttpClient http.Client //http客户端-访问jira
	userHome       string
)

const (
	docDir = "Documents"
	drDir  = "dailyReport"
)

func init() {
	//if testFlag {
	//	return
	//}
	mainInit()

	//初始化isHoliday,isTX函数
	holidayJudgeInit()

	nextWorkDayInit()

	err := templateInit()
	if err != nil {
		warn("程序初始化失败: %v", err)
		return
	}
	jiraInit()
}

//初始化下个工作日
func nextWorkDayInit() {
	nextWorkDay = func() time.Time {
		nwd := now
		for true {
			nwd = nwd.AddDate(0, 0, 1)
			if isWorkDay(nwd) {
				break
			}
		}
		lg("下个工作日推断为%s\n", nwd.Format(dateFormat))
		return nwd
	}()
}

//初始化isHoliday,isTX函数
func holidayJudgeInit() {
	isHoliday, isTX = func() (func(d time.Time) bool, func(d time.Time) bool) {
		isHoliday_ := func(date time.Time) bool {
			y := date.Year()
			holidaysOfYear, ok := holidaysMap[y]
			if !ok {
				//实际上是不会走到这里的
				warn("请在假期表中维护%d年的假期及调休数据", y)
				return false
			}
			for _, h := range holidaysOfYear {
				if time.Month(h.m) == date.Month() && h.d == date.Day() {
					//找到当前日期
					return h.t == rest
				}
			}
			return false
		}
		isHoliday__ := func(time.Time) bool { return false }
		//判断是否是调休
		isTX_ := func(date time.Time) bool {
			y := date.Year()
			holidaysOfYear, ok := holidaysMap[y]
			if !ok {
				//实际上是不会走到这里的
				warn("请在假期表中维护%d年的假期及调休数据 ", y)
				return false
			}
			for _, h := range holidaysOfYear {
				if time.Month(h.m) == date.Month() && h.d == date.Day() {
					//找到当前日期
					return h.t == work
				}
			}
			return false
		}
		isTX__ := func(time.Time) bool { return false }

		balanceFlag := holidayBalanceByNow()

		switch balanceFlag {
		case enough: //维护的假期充足,应用正常的假期/调休判断方法
			lg("假期库充足")
			return isHoliday_, isTX_
		case exhausting: //维护的假期即将耗尽,应用正常的假期/调休判断方法,但给出警告
			warn("假期库即将耗尽, 请尽快维护")
			return isHoliday_, isTX_
		default: //exhausted 维护的假期库已经耗尽,假期判断和调休判断总是返回false,并且给出警告
			warn("假期库已经耗尽, 法定节假日的推断逻辑将被禁用, 下个工作日的推断可能会不准确,请尽快维护假期库")
			return isHoliday__, isTX__
		}
	}()
}

//初始化main
func mainInit() {
	params = parseCmd()
	//控制日志输出模式
	if !params.Verbose {
		//未开启啰嗦模式,启用静默模式
		lg = lgSilence
	}
}

//日志函数有两种定义:啰嗦模式于静默模式
//啰嗦模式下
//若args长度为0,message后会主动追加换行
//若args长度非零,message只会进行格式化,不会进行换行符追加
//啰嗦日志输出
func lgVerbose(message string, args ...interface{}) {
	outputLock.Lock()
	defer outputLock.Unlock()
	fmt.Printf("\r")//先清除前方输出(主要是干掉spinner输出)
	if len(args) == 0 {
		log.Println(message)
	} else {
		log.Printf(message, args...)
	}
}

//静默模式日志输出
func lgSilence(message string, args ...interface{}) {
	//什么也不做
}

//模板文件初始化
func templateInit() error {
	//检查userhome/Documents/dailyReport/dr.template文件是否存在
	//若不存在,创建之
	home, err := os.UserHomeDir()
	if err != nil {
		warn("用户目录获取失败,程序终止: %v", err)
		return err
	}
	userHome = home
	const TemplateFileName = "dr.template"
	targetFile := fmt.Sprintf("%s%c%s%c%s%c%s",
		userHome, os.PathSeparator,
		docDir, os.PathSeparator,
		drDir, os.PathSeparator,
		TemplateFileName)
	templateFile, err = os.Open(targetFile)
	if templateFile != nil {
		lg("发现模板文件%s\n", targetFile)
		return nil //找到文件,直接结束函数
	}
	if !os.IsNotExist(err) {
		//如果返回的err不是文件存在,停止进程
		warn("打开文件异常: %v", err)
		return err
	}
	//文件不存在,创建之,并写入默认模板
	templateFile, err = os.Create(targetFile)
	if err != nil {
		warn("创建文件异常: %v", err)
		return err
	}

	const templateDefaultContent = `收件人：
{{.Receiver}}

抄送：
{{.Cc}}

主题：
{{.DepName}}日报 - {{.Date}} - {{.ReporterName}}

内容：

{{.DepName}}日报 - {{.Date}} - {{.ReporterName}}
----------------------------------------------------------------------------------------------------------------------
当天工作成果
{{.Achievements}}
----------------------------------------------------------------------------------------------------------------------
碰到问题&解决方案
{{.PAndS}}
----------------------------------------------------------------------------------------------------------------------
明天工作计划
{{.Targets}}
`
	_, err = templateFile.WriteString(templateDefaultContent)
	if err != nil {
		warn("模板文件写入异常: %v", err)
		return err
	}
	lg("已创建模板文件%s\n", targetFile)
	//写入后重新打开
	templateFile.Close()
	templateFile, err = os.Open(targetFile)
	if err != nil {
		warn("创建并重新打开模板文件异常: %v", err)
		return err
	}
	return nil
}

//初始化httpClient,并设置cookie管理
func jiraInit() {
	jiraCookieJar, _ := cookiejar.New(nil) //根据源代码,这个函数并不会产生err,这里忽略返回值中的err
	jiraHttpClient = http.Client{Jar: jiraCookieJar}
}
