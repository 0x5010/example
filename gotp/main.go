package main

import (
	"fmt"
	"time"

	"github.com/xlzd/gotp"
)

func main() {
	defaultTOTPUsage()
	defaultHOTPUsage()
}

func defaultTOTPUsage() {
	otp := gotp.NewDefaultTOTP("4S62BZNFXXSZLCRO") // NewTOTP(secret, 6, 30, nil)

	password := otp.Now()
	fmt.Println("current ont-time password is:", password)
	fmt.Println("one-time password of timestamp 0 is:", otp.At(0))
	fmt.Println(otp.ProvisioningUri("demoAccountName", "issuerName"))

	fmt.Println(otp.Verify("179394", 1524485781))

	time.Sleep(1 * time.Second)
	fmt.Println(otp.Verify(password, int(time.Now().Unix())))
	time.Sleep(30 * time.Second)
	fmt.Println(otp.Verify(password, int(time.Now().Unix())))
}

func defaultHOTPUsage() {
	otp := gotp.NewDefaultHOTP("4S62BZNFXXSZLCRO")

	password0 := otp.At(0)
	password1 := otp.At(1)
	fmt.Println("one-time password of counter 0 is:", password0)
	fmt.Println(otp.ProvisioningUri("demoAccountName", "issuerName", 1))

	fmt.Println(otp.Verify(password0, 0))
	fmt.Println(otp.Verify(password0, 1))
	fmt.Println(otp.Verify(password1, 0))
	fmt.Println(otp.Verify(password1, 1))
	fmt.Println(otp.Verify(password1, 2))
}
