package language

import (
	"job-interviewer/pkg/language"
)

type DefaultService struct {
	storages                 map[language.Language]language.Storage
	userLanguageStorage      language.Storage
	interviewLanguageStorage language.Storage
}

func (s *DefaultService) GetTextFromAllLanguages(key language.TextKey) []string {
	result := make([]string, 0, len(s.storages))
	for _, storage := range s.storages {
		text := storage.GetText(key)
		if text == "" {
			continue
		}

		result = append(result, text)
	}

	return result
}

func NewService(dictionaries map[language.Language]Dictionary) *DefaultService {
	storages := make(map[language.Language]language.Storage, len(dictionaries))
	for lang, dict := range dictionaries {
		storages[lang] = language.NewStorage(dict.GetTexts())
	}

	return &DefaultService{storages: storages}
}

func (s *DefaultService) GetText(lang language.Language, key language.TextKey) string {
	storage, ok := s.storages[lang]
	if !ok {
		return ""
	}

	return storage.GetText(key)
}
