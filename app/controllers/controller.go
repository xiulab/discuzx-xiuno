package controllers

// Controller Controller
type Controller interface {
	ToConvert() (err error) // 转换
}
