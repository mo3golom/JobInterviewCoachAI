package language

type Storage interface {
	GetText(key TextKey) string
}
