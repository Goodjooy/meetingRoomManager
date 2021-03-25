package manage
func Contains(target interface{},array []interface{}) bool {
	var contain=false;
	for _, v := range array {
		if v==target{
			contain=true
			break
		}
	}
	return contain;
}