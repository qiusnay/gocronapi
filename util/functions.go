package util

func GetMyConf() {

}

func InArray(need interface{}, needArr []string) bool {
	for _, v := range needArr {
		if need == v {
			return true
		}
	}
	return false
}
