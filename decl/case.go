package decl

func RuneToCase(r rune) string {
	caseText := ""
	switch r {
	case '1':
		caseText = "主格"
	case '2':
		caseText = "属格"
	case '3':
		caseText = "与格"
	case '4':
		caseText = "対格"
	case 'a':
		caseText = "奪格"
	case 'v':
		caseText = "呼格"
	case 'l':
		caseText = "地格"
	default:
		panic("unknown case")
	}
	return caseText
}
