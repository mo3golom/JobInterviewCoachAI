package language

import (
	"errors"
	"job-interviewer/pkg/language"
)

type DefaultService struct {
	dictionaries             map[language.Language]Dictionary
	userLanguageStorage      language.Storage
	interviewLanguageStorage language.Storage
}

func NewService(dictionaries map[language.Language]Dictionary) *DefaultService {
	return &DefaultService{dictionaries: dictionaries}
}

func (s *DefaultService) InitUserLanguage(lang language.Language) error {
	storage, err := s.init(lang)
	if err != nil {
		return err
	}

	s.userLanguageStorage = storage
	return nil
}

func (s *DefaultService) GetUserLanguageText(key language.TextKey) string {
	return s.userLanguageStorage.GetText(key)
}

func (s *DefaultService) InitInterviewLanguage(lang language.Language) error {
	storage, err := s.init(lang)
	if err != nil {
		return err
	}

	s.interviewLanguageStorage = storage
	return nil
}

func (s *DefaultService) GetInterviewLanguageText(key language.TextKey) string {
	return s.interviewLanguageStorage.GetText(key)
}

func (s *DefaultService) init(lang language.Language) (language.Storage, error) {
	dict, ok := s.dictionaries[lang]
	if !ok {
		return nil, errors.New("not found language dictionary")
	}

	return language.NewStorage(dict.GetTexts()), nil
}
