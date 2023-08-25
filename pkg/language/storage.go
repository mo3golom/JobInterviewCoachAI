package language

type DefaultWordStorage struct {
	texts map[TextKey]string
}

func (s *DefaultWordStorage) GetText(key TextKey) string {
	return s.texts[key]
}

func NewWordStorage(texts map[TextKey]string) *DefaultWordStorage {
	return &DefaultWordStorage{texts: texts}
}

type DefaultLangStorage struct {
	languages map[Language]WordStorage
}

func NewLangStorage(languages map[Language]WordStorage) *DefaultLangStorage {
	return &DefaultLangStorage{languages: languages}
}

func (d *DefaultLangStorage) GetText(lang Language, key TextKey) string {
	storage, ok := d.languages[lang]
	if !ok {
		return ""
	}

	return storage.GetText(key)
}
