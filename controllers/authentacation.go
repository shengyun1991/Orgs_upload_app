package controllers

import (
	"Orgs_upload_app/models"
	"encoding/csv"
	"fmt"
	"Orgs_upload_app/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/toolbox"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type AuthControllers struct{
	beego.Controller
}

func (this * AuthControllers) RecodeAuth(){
	beego.Debug("---- Auth Recode File ----")

	file,header,err := this.GetFile("auth")
	if err!=nil{
		utils.HandleResponse(this.Ctx.ResponseWriter,400,err.Error())
		return
	}

	defer file.Close()

	filename := header.Filename
	fmt.Print("文件名称:")
	fmt.Println(filename)

	// 保存文件到指定的位置
	// static/uploadfile
	err = this.SaveToFile("auth", path.Join("static/uploadfile", filename))
	if err != nil {
		utils.HandleResponse(this.Ctx.ResponseWriter, 500, err.Error())
	}

	// 为了保证文件能够上传完毕,设置任务 5 秒后开启
	t := time.Now().Add(5 * time.Second)
	second := t.Second()
	minute := t.Minute()
	hour := t.Hour()

	spec := fmt.Sprintf("%d %d %d * * *", second, minute, hour)
	tk := toolbox.NewTask("AuthTask", spec,
		func() error {
			beego.Info("Start task")
			defer toolbox.StopTask()
			return AuthTask(filename)
		})

	toolbox.AddTask("AuthTask", tk)
	toolbox.StartTask()

	utils.HandleResponse(this.Ctx.ResponseWriter, 200, "ok")
}

// 编写我们的任务函数
func AuthTask(fileName string) error {
	// 任务的主要职责就是 解析文件并将数据写入区块链
	var (
		channelID   = beego.AppConfig.String("channel_id_gaj")
		chaincodeId = beego.AppConfig.String("chaincode_id_auth")
		userID      = beego.AppConfig.String("user_id_auth")
	)

	fmt.Println("---- channelID = :",channelID)
	fmt.Println("---- chaincodeID = :",chaincodeId)
	fmt.Println("---- UserID = :",userID)
	ccs, err := models.Initialize(channelID, chaincodeId, userID,"CORE_OGAJ_CONFIG_FILE")
	if err != nil {
		beego.Error(err.Error())
	}
	defer ccs.Close()

	file, _ := os.Open(path.Join("static/uploadfile", fileName))
	reader := csv.NewReader(file)

	// 我们在上传大量数据的时候,有可能少部分上传会出现问题,我们不能因为这一小部分问题结束上传的进度
	// 出现问题的情况:1.某些数据读取错误 2.某些数据不完整.....
	// 解决办法:我们要做的就是记录出错的数据信息,回头单独处理
	var line = 0
	var lines []string
	for {
		line += 1
		linestr := strconv.Itoa(line)
		record, err := reader.Read()
		if err == io.EOF {
			//文件结尾
			break
		}
		// 读取 文件出错
		if err != nil {
			lines = append(lines, linestr)
			continue
		}
		// 数据不完整
		if len(record) != 3 {
			lines = append(lines, linestr)
			continue
		}

		var args [][]byte
		for _, str := range record {
			args = append(args, []byte(str))
		}

		_, err = ccs.ChainCodeUpdate("add", args)
		if err != nil {

			lines = append(lines, linestr)
		}

		fmt.Printf("第 %d 条数据已成功写入...",line)
		fmt.Println()
	}

	if len(lines) > 0 {
		beego.Error("Error lines:", strings.Join(lines, "."))
	} else {
		beego.Info("write data success!")
	}

	return nil
}