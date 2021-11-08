package customlogger

import "fmt"

func Infof(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	fmt.Printf("[INFO]: %s", s)
}

func Errorf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	fmt.Printf("[ERROR]: %s", s)
}

func Error(err error, format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	fmt.Printf("[ERROR]: Message: %s, Error: %s", s, err.Error())
}
