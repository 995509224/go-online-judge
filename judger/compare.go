package judger

func compare(a, b string) bool {
	u := []byte(a)
	v := []byte(b)
	var i, j int = 0, 0
	for {
		if i == len(a) {
			if j == len(b) {
				return true
			}
			for j < len(b) {
				if b[j] != ' ' && b[j] != '\n' && b[j] != '\r' {
					return false
				}
				j++
			}
			return true
		}
		if j == len(b) {
			if i == len(a) {
				return true
			}
			for i < len(a) {
				if a[i] != ' ' && a[i] != '\n' && a[i] != '\r' {
					return false
				}
				i++
			}
			return true
		}
		if u[i] == ' ' || u[i] == '\n' || u[i] == '\r' {
			i++
			continue
		}
		if v[j] == ' ' || v[i] == '\n' || v[i] == '\r' {
			j++
			continue
		}
		if u[i] != v[i] {
			return false
		}
		i++
		j++
	}
}
