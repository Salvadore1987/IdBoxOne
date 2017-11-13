package lib

import (
	"syscall"
	"fmt"
	"unsafe"
)

type Error uint32
type MRTD struct {}

var (
	modmrtd = syscall.NewLazyDLL("lib/ElyMRTD.dll")
	procCalculateChecksum = modmrtd.NewProc("calculateChecksum")
	procConnect = modmrtd.NewProc("connect")
	procDisconnect = modmrtd.NewProc("disconnect")
	procInit = modmrtd.NewProc("init")
	procEstablishBAC = modmrtd.NewProc("establishBAC")
	procReadDG1 = modmrtd.NewProc("readDG1")
	procReadDG2 = modmrtd.NewProc("readDG2")
	procGetName = modmrtd.NewProc("getName")
	procGetSurname = modmrtd.NewProc("getSurname")
	procGetSex = modmrtd.NewProc("getSex")
	procGetBirthDate = modmrtd.NewProc("getBirthDate")
	procGetValidityDate = modmrtd.NewProc("getValidityDate")
	procGetDocNum = modmrtd.NewProc("getDocNum")
	procGetDocumentType = modmrtd.NewProc("getDocumentType")
	procGetOptionalData = modmrtd.NewProc("getOptionalData")
	procGetNationality = modmrtd.NewProc("getNationality")
	procGetIssuingState = modmrtd.NewProc("getIssuingState")
	procGetCountryName = modmrtd.NewProc("getCountryNameL")
	procGetDG2 = modmrtd.NewProc("getDG2")
)

func (e Error) Error() string {
	err := syscall.Errno(e)
	return fmt.Sprintf("mrtd: error(%0x): %s", uintptr(e), err.Error())
}

func NewMRTD() *MRTD {
	_, _, _ = procInit.Call()
	return &MRTD{}
}

func CalculateChecksum(str string) bool {
	pasPtr, err := syscall.BytePtrFromString(str)
	if err != nil {
		return false
	}
	r, _, _ := procCalculateChecksum.Call(uintptr(unsafe.Pointer(pasPtr)))
	if r == 0 {
		return false
	}
	return true
}

func (t *MRTD) ConnectReader(reader string) int {
	readerPtr, err := syscall.BytePtrFromString(reader)
	if err != nil {
		return -1
	}
	r, _, _ := procConnect.Call(uintptr(unsafe.Pointer(readerPtr)))
	return int(r)
}

func (t *MRTD) EstablishBAC(mrz_info string) bool {
	mrzPtr, err := syscall.BytePtrFromString(mrz_info)
	if err != nil {
		return false
	}
	r, _, _ := procEstablishBAC.Call(uintptr(unsafe.Pointer(mrzPtr)))
	if r == 0 {
		return false
	}
	return true
}

func DisconnectReader()  {
	_, _, _ = procDisconnect.Call()
}

func (t* MRTD) ReadDG1() int {
	r, _, _ := procReadDG1.Call()
	return int(r)
}

func (t* MRTD) ReadDG2() int {
	r, _, _ := procReadDG2.Call()
	return int(r)
}

func (t* MRTD) GetName() string {
	return getName()
}

func (t* MRTD) GetSurname() string {
	return getSurname()
}

func (t* MRTD) GetSex() string {
	return getSex()
}

func (t* MRTD) GetBirthDate() string {
	return getBirthDate()
}

func (t* MRTD) GetValidityDate() string {
	return getValidityDate()
}

func (t* MRTD) GetDocNum() string {
	return getDocNum()
}

func (t* MRTD) GetDocumentType() string {
	return getDocumentType()
}

func (t* MRTD) GetOptionalData() string {
	return getOptionalData()
}

func (t* MRTD) GetNationality() string {
	return getNationality()
}

func (t* MRTD) GetIssuingState() string {
	return getIssuingState()
}

func (t* MRTD) GetCountryName(info string) string {
	return getCountryName(info)
}

func (t* MRTD) GetDg2() []byte {
	return getDG2()
}

func getDG2() []byte {
	var (
		length uint = 0
		handle *syscall.Handle
	)
	_, _, _ = procGetDG2.Call(uintptr(unsafe.Pointer(handle)), uintptr(unsafe.Pointer(&length)))
	data := make([]byte, length)
	_, _, _ = procGetDG2.Call(uintptr(unsafe.Pointer(&data[0])), uintptr(unsafe.Pointer(&length)))
	return data
}

func getCountryName(info string) string {
	var (
		length uint = 0
		handle *syscall.Handle
	)
	infoBytes := []byte(info)
	_, _, _ = procGetCountryName.Call(uintptr(unsafe.Pointer(&infoBytes[0])), uintptr(unsafe.Pointer(handle)), uintptr(unsafe.Pointer(&length)))
	data := make([]byte, length)
	_, _, _ = procGetCountryName.Call(uintptr(unsafe.Pointer(&infoBytes[0])), uintptr(unsafe.Pointer(&data[0])), uintptr(unsafe.Pointer(&length)))
	return getStringFromByteArray(data)
}

func getIssuingState() string {
	var (
		length uint = 0
		handle *syscall.Handle
	)
	_, _, _ = procGetIssuingState.Call(uintptr(unsafe.Pointer(handle)), uintptr(unsafe.Pointer(&length)))
	data := make([]byte, length)
	_, _, _ = procGetIssuingState.Call(uintptr(unsafe.Pointer(&data[0])), uintptr(unsafe.Pointer(&length)))
	return getStringFromByteArray(data)
}

func getNationality() string {
	var (
		length uint = 0
		handle *syscall.Handle
	)
	_, _, _ = procGetNationality.Call(uintptr(unsafe.Pointer(handle)), uintptr(unsafe.Pointer(&length)))
	data := make([]byte, length)
	_, _, _ = procGetNationality.Call(uintptr(unsafe.Pointer(&data[0])), uintptr(unsafe.Pointer(&length)))
	return getStringFromByteArray(data)
}

func getOptionalData() string {
	var (
		length uint = 0
		handle *syscall.Handle
	)
	_, _, _ = procGetOptionalData.Call(uintptr(unsafe.Pointer(handle)), uintptr(unsafe.Pointer(&length)))
	data := make([]byte, length)
	_, _, _ = procGetOptionalData.Call(uintptr(unsafe.Pointer(&data[0])), uintptr(unsafe.Pointer(&length)))
	return getStringFromByteArray(data)
}

func getDocumentType() string {
	var (
		length uint = 0
		handle *syscall.Handle
	)
	_, _, _ = procGetDocumentType.Call(uintptr(unsafe.Pointer(handle)), uintptr(unsafe.Pointer(&length)))
	data := make([]byte, length)
	_, _, _ = procGetDocumentType.Call(uintptr(unsafe.Pointer(&data[0])), uintptr(unsafe.Pointer(&length)))
	return getStringFromByteArray(data)
}

func getDocNum() string {
	var (
		length uint = 0
		handle *syscall.Handle
	)
	_, _, _ = procGetDocNum.Call(uintptr(unsafe.Pointer(handle)), uintptr(unsafe.Pointer(&length)))
	data := make([]byte, length)
	_, _, _ = procGetDocNum.Call(uintptr(unsafe.Pointer(&data[0])), uintptr(unsafe.Pointer(&length)))
	return getStringFromByteArray(data)
}

func getValidityDate() string {
	var (
		length uint = 0
		handle *syscall.Handle
	)
	_, _, _ = procGetValidityDate.Call(uintptr(unsafe.Pointer(handle)), uintptr(unsafe.Pointer(&length)))
	data := make([]byte, length)
	_, _, _ = procGetValidityDate.Call(uintptr(unsafe.Pointer(&data[0])), uintptr(unsafe.Pointer(&length)))
	return getStringFromByteArray(data)
}

func getBirthDate() string {
	var (
		length uint = 0
		handle *syscall.Handle
	)
	_, _, _ = procGetBirthDate.Call(uintptr(unsafe.Pointer(handle)), uintptr(unsafe.Pointer(&length)))
	birth := make([]byte, length)
	_, _, _ = procGetBirthDate.Call(uintptr(unsafe.Pointer(&birth[0])), uintptr(unsafe.Pointer(&length)))
	return getStringFromByteArray(birth)
}

func getSex() string {
	var (
		length uint = 0
		handle *syscall.Handle
	)
	_, _, _ = procGetSex.Call(uintptr(unsafe.Pointer(handle)), uintptr(unsafe.Pointer(&length)))
	sex := make([]byte, length)
	_, _, _ = procGetSex.Call(uintptr(unsafe.Pointer(&sex[0])), uintptr(unsafe.Pointer(&length)))
	return getStringFromByteArray(sex)
}

func getSurname() string {
	var (
		length uint = 0
		handle *syscall.Handle
	)
	_, _, _ = procGetSurname.Call(uintptr(unsafe.Pointer(handle)), uintptr(unsafe.Pointer(&length)))
	surname := make([]byte, length)
	_, _, _ = procGetSurname.Call(uintptr(unsafe.Pointer(&surname[0])), uintptr(unsafe.Pointer(&length)))
	return getStringFromByteArray(surname)
}

func getName() string {
	var (
		length uint = 0
		handle *syscall.Handle
	)
	_, _, _ = procGetName.Call(uintptr(unsafe.Pointer(handle)), uintptr(unsafe.Pointer(&length)))
	name := make([]byte, length)
	_, _, _ = procGetName.Call(uintptr(unsafe.Pointer(&name[0])), uintptr(unsafe.Pointer(&length)))
	return getStringFromByteArray(name)
}

func getStringFromByteArray(arr []byte) string {
	var newArr []byte
	for i := range arr {
		if arr[i] != 0 {
			newArr = append(newArr, arr[i])
		}
	}
	return string(newArr)
}