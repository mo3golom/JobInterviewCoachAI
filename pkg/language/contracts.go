package language

type (
	WordStorage interface {
		GetText(key TextKey) string
	}

	Storage interface {
		GetText(lang Language, key TextKey) string
	}
)
