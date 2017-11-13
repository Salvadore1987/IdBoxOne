package lib

import (
	"go.bug.st/serial.v1"
	"errors"
	"time"
)

var port serial.Port

func Connect() error {
	return connect()
}

func Inquire() error {
	return inquire()
}

func ReadMRZ() (string, error) {
	return readMRZ()
}

// Метод освобождает COM - порт
func Disconnect() error {
	err := port.Close(); if err != nil {
		return err
	}
	return nil
}

// Метод предназначен для считывания данных MRZ с COM - порта сканера.
func readMRZ() (string, error) {
	buff := make([]byte, 256)
	var (
		mrz_str string
	)
	time.Sleep(time.Second * 3)
	n, err := port.Read(buff)
	if err != nil {
		return "", err
	}
	mrz_str = string(buff[:n])
	return mrz_str, nil
}

// Метод служит для соединения с COM - портом сканера
func connect() error {
	ports, err := serial.GetPortsList(); if err != nil {
		return err
	}
	if len(ports) == 0 {
		return errors.New("No serial ports found!")
	}
	mode := &serial.Mode{
		Parity: serial.NoParity,
		DataBits: 8,
		BaudRate: 9600}
	port, err = serial.Open(ports[0], mode); if err != nil {
		return err
	}
	port.SetDTR(true)
	err = port.ResetInputBuffer(); if err != nil {
		return err
	}
	err = port.ResetOutputBuffer(); if err != nil {
		return err
	}
	return nil
}

// Метод отсылает на сканер команду через COM - порт для считывания MRZ
func inquire() error {
	inquire := []byte{0x49, 0, 0}
	err := port.ResetInputBuffer(); if err != nil {
		return err
	}
	_, err = port.Write(inquire); if err != nil {
		return err
	}
	return nil
}
