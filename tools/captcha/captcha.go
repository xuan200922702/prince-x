package captcha

import (
	"github.com/google/uuid"
	"github.com/mojocn/base64Captcha"
	"image/color"
)

var store = base64Captcha.DefaultMemStore
var PrinceDriverDigit = NewDriverDigit(80, 240, 4, 0.7, 80)

//DriverDigit config for captcha-engine-digit.
type DriverDigit struct {
	// Height png height in pixel.
	Height int
	// Width Captcha png width in pixel.
	Width int
	// DefaultLen Default number of digits in captcha solution.
	Length int
	// MaxSkew max absolute skew factor of a single digit.
	MaxSkew float64
	// DotCount Number of background circles.
	DotCount int
}

//configJsonBody json request body.
type ConfigJsonBody struct {
	Id            string
	CaptchaType   string
	VerifyValue   string
	DriverAudio   *base64Captcha.DriverAudio
	DriverString  *base64Captcha.DriverString
	DriverChinese *base64Captcha.DriverChinese
	DriverMath    *base64Captcha.DriverMath
	DriverDigit   *base64Captcha.DriverDigit
}

func DriverStringFunc() (id, b64s string, err error) {
	e := ConfigJsonBody{}
	e.Id = uuid.New().String()
	e.DriverString = base64Captcha.NewDriverString(46, 140, 2, 2, 4, "234567890abcdefghjkmnpqrstuvwxyz", &color.RGBA{240, 240, 246, 246}, []string{"wqy-microhei.ttc"})
	driver := e.DriverString.ConvertFonts()
	cap := base64Captcha.NewCaptcha(driver, store)
	return cap.Generate()
}

func NewDriverDigit(height int, width int, length int, maxSkew float64, dotCount int) *DriverDigit {
	return &DriverDigit{Height: height, Width: width, Length: length, MaxSkew: maxSkew, DotCount: dotCount}
}

func DriverDigitFunc() (id, b64s string, err error) {
	e := ConfigJsonBody{}
	e.Id = uuid.New().String()
	//DefaultDriverDigit is a default driver of digit

	e.DriverDigit = (*base64Captcha.DriverDigit)(PrinceDriverDigit)
	driver := e.DriverDigit
	cap := base64Captcha.NewCaptcha(driver, store)
	return cap.Generate()
}
