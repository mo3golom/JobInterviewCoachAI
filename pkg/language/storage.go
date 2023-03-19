package language

type DefaultStorage struct {
	texts map[TextKey]string
}

func NewStorage(texts map[TextKey]string) DefaultStorage {
	return DefaultStorage{texts: texts}
}

func (s DefaultStorage) GetText(key TextKey) string {
	return s.texts[key]
}
