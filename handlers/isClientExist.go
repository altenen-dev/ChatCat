package handlers

func isClientExist(name string) bool {
	for _, v := range Clients {
		if v.Name == name {
			return true
		}
	}
	return false
}
