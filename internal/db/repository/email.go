package repository

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"sync"

	"github.com/rs/zerolog"
)

type EmailRepository struct {
	sender    string
	templates map[string]string
	mutex     *sync.RWMutex
}

func NewEmailRepository(logger *zerolog.Logger) (*EmailRepository, error) {
	templateFiles := make(map[string]string)

	err := filepath.Walk("./internal/email/templates/build/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		defer func() {
			if err := file.Close(); err != nil {
				logger.Panic().Err(err).Msg("unable to close file")
			}
		}()

		info, err2 := file.Stat()

		if err2 != nil {
			return err2
		}

		fileByteContents := make([]byte, info.Size())
		_, err3 := file.Read(fileByteContents)

		if err3 != nil {
			return err3
		}

		fileName := filepath.Base(path)
		templateFiles[fileName[:len(fileName)-len(filepath.Ext(fileName))]] = string(fileByteContents[:])

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &EmailRepository{templates: templateFiles}, nil
}

func (emailRepo *EmailRepository) GetTemplates() *map[string]string {
	return &emailRepo.templates
}

func (emailRepo *EmailRepository) GetTemplate(name string) string {
	if len(emailRepo.templates) == 0 {
		return ""
	}

	emailRepo.mutex.RLock()
	templateString, ok := emailRepo.templates[name]
	emailRepo.mutex.RUnlock()

	if !ok {
		return ""
	}

	return templateString
}

func (emailRepo *EmailRepository) SetTemplate(name string, content string) {
	emailRepo.mutex.Lock()
	emailRepo.templates[name] = content
	emailRepo.mutex.Unlock()
}

func (emailRepo *EmailRepository) GetSender() string {
	return emailRepo.sender
}

func (emailRepo *EmailRepository) GetTemplateCount() int {
	return len(emailRepo.templates)
}

func (emailRepo *EmailRepository) SetSender(sender string) {
	emailRepo.sender = sender
}

func (emailRepo *EmailRepository) RenderTemplate(name string, variables interface{}) string {
	if len(emailRepo.templates) == 0 {
		return ""
	}

	rawTemplate, ok := emailRepo.templates[name]

	if !ok {
		return ""
	}

	templateInstance, err := template.New(name).Parse(rawTemplate)
	if err != nil {
		return ""
	}

	buffer := new(bytes.Buffer)

	if err := templateInstance.Execute(buffer, variables); err != nil {
		return ""
	}

	return buffer.String()
}
