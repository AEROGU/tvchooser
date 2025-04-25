package tvclang

type Texts struct {
	Cancel       string
	Accept       string
	Modfied      string
	Size         string
	AccessDenied string
	ThisPC       string
	HomeDir      string
	Devices      string
	Favorites    string
}

var translations *Texts

func GetTranslations() Texts {
	if translations == nil {
		defaultLang := LangEnglish()
		translations = &defaultLang
	}
	return *translations
}

func SetTranslations(texts Texts) {
	translations = &texts
}
