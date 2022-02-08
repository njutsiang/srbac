package srbac

// 年月日时分秒常量
var TimeYmdhis = "2006-01-02 15:04:05"

// 判断是否有错误，有则记录错误日志，并抛出 panic
func CheckError(err interface{}) {
	if err == nil {
		return
	}
	panic(err)
}