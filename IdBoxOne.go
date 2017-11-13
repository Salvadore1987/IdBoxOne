package main

import (
	"IdBoxOne/lib"
	"github.com/ebfe/scard"
	"regexp"
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/astaxie/beego/logs"
	"sync"
	"github.com/astaxie/beego/config"
	"encoding/base64"
)

var (
	mrtd *lib.MRTD
	response *lib.Response
	log *logs.BeeLogger
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {return true},
	}
	wg *sync.WaitGroup
)

func main()  {
	wg = new(sync.WaitGroup)
	log = logs.NewLogger(10000)
	log.SetLogger(logs.AdapterConsole)
	log.SetLogger(logs.AdapterFile, `{"filename":"logs/project.log"}`)
	conf, err := config.NewConfig("ini", "conf/app.ini")
	if err != nil {
		log.Error("%s", err)
		return
	}
	log.Info("Start application on port %s", conf.String("port"))

	http.HandleFunc("/info", info)
	log.Error("%s", http.ListenAndServe(":" + conf.String("port"), nil))
}

func info(w http.ResponseWriter, r *http.Request) {

	log.Info("%s", "Start echo")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("%s", err)
	}

	wg.Add(2)
	response = lib.NewResponse()
	go readPassportInfo()
	go readICCID()
	wg.Wait()

	if err = conn.WriteJSON(response); err != nil {
		log.Error("%s", err)
		return
	}

}

func getDG1Info() *lib.PersonalInfo {
	dg1 := lib.NewPersonalInfo()
	dg1.Name = mrtd.GetName()
	dg1.Surname = mrtd.GetSurname()
	dg1.Sex = mrtd.GetSex()
	dg1.BirthDate = mrtd.GetBirthDate()
	dg1.ValidityDate = mrtd.GetValidityDate()
	dg1.DocNumber = mrtd.GetDocNum()
	dg1.DocType = mrtd.GetDocumentType()
	dg1.OptionalData = mrtd.GetOptionalData()
	dg1.Country = mrtd.GetCountryName(mrtd.GetNationality())
	dg1.IssuingState = mrtd.GetCountryName(mrtd.GetIssuingState())
	return dg1
}

func getDG2Info() []byte {
	return mrtd.GetDg2()
}

func readPassportInfo()  {

	defer wg.Done()

	// Подключаемся к COM - порту
	err := lib.Connect(); if err != nil {
		log.Error("%s", err)
	}

	// Освобождаем COM - порт
	defer lib.Disconnect()

	// Посылаем команду сканеру через COM - порт, для извлечения MRZ
	err = lib.Inquire(); if err != nil {
		log.Error("%s", err)
	}

	// Считываем MRZ
	mrz, err := lib.ReadMRZ(); if err != nil {
		log.Error("%s", err)
	}

	// Создаем объект MRZ - парсера.
	mrz_info := lib.NewMRZ(mrz)

	// Проверяем контрольную сумму MRZ
	if mrz_info.GetChecksum() {
		var readerCL string
		ctx, err := scard.EstablishContext(); if err != nil {
			log.Error("%s", err)
		}
		defer ctx.Release()
		listReaders, err := ctx.ListReaders(); if err != nil {
			log.Error("%s", err)
		}
		for i := range listReaders {
			matched, err := regexp.MatchString("ELYCTIS CL reader", listReaders[i])
			if err != nil {
				log.Error("%s", err)
				break
			}
			if matched {
				readerCL = listReaders[i]
			}
		}
		mrtd := lib.NewMRTD()
		res := mrtd.ConnectReader(readerCL)
		defer lib.DisconnectReader()
		if res == 0 {
			log.Info("%s", "ICAO")
			if !(mrtd.EstablishBAC(mrz_info.GetMRZInfo())) {
				log.Error("%s", "BAC not established")
				return
			}
			log.Info("%s", "BAC Established")
			iRes := mrtd.ReadDG1()
			if iRes >= 0 {
				dg1 := getDG1Info()
				response.PassportData = dg1
			}
			iRes = mrtd.ReadDG2()
			if iRes >= 0 {
				dg2 := getDG2Info()
				sEnc := base64.URLEncoding.EncodeToString(dg2)
				response.Photo = sEnc
			}
		} else if res == -1 {
			log.Error("%s", "Error when connect reader")
		}
	} else {
		log.Error("%s", "Неверные данные")
	}
	log.Info("%s", "Read Passport Info Finished!")
}

func readICCID()  {

	defer wg.Done()

	var readerCL string
	ctx, err := scard.EstablishContext(); if err != nil {
		log.Error("%s", err)
	}
	defer ctx.Release()
	listReaders, err := ctx.ListReaders(); if err != nil {
		log.Error("%s", err)
	}
	for i := range listReaders {
		matched, err := regexp.MatchString("ELYCTIS CNT reader 0", listReaders[i])
		if err != nil {
			log.Error("%s", err)
			return
		}
		if matched {
			readerCL = listReaders[i]
		}
	}
	card, err := ctx.Connect(readerCL, scard.ShareShared, scard.ProtocolAny)
	if err != nil {
		log.Error("Error Connect: ", err)
		return
	}
	defer card.Disconnect(scard.LeaveCard)
	var cmdSelectFile = []byte{
		0xA0,
		0xA4,
		0x00,
		0x00,
		0x02,
		0x2F,
		0xE2}
	_, err = card.Transmit(cmdSelectFile)
	if err != nil {
		log.Error("Error send command: ", err)
		return
	}
	var cmdGetResponse = []byte{
		0xA0,
		0xC0,
		0x00,
		0x00,
		0x0F}
	_, err = card.Transmit(cmdGetResponse)
	if err != nil {
		log.Error("Error send command: ", err)
		return
	}
	var cmdReadBinary = []byte{
		0xA0,
		0xB0,
		0x00,
		0x00,
		0x0A}
	rsp, err := card.Transmit(cmdReadBinary)
	if err != nil {
		log.Error("Error send command: ", err)
		return
	}
	iccid := swapNibbles(rsp)
	response.ICCID = fmt.Sprintf("%x", iccid)
	log.Info("%s", "Read ICCID Finished!")
}

func swapNibbles(nibbles []byte) []byte {
	var resp []byte
	for i := 0; i < 9; i++ {
		nibbleOld := nibbles[i]
		nibbleNew := (nibbleOld >> 4) | (nibbleOld << 4)
		resp = append(resp, nibbleNew)
	}
	return resp
}
