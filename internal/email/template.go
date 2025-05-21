package email

import (
	"VEDA95/open_board/api/internal/log"
	"bytes"
	"errors"
	"html/template"
	"os"
	"path/filepath"
)

var MailTemplateStore *TemplateStore

type TemplateStore struct {
	templates map[string]string
}

func InitializeEmailTemplateStore() error {
	if MailClient == nil {
		return errors.New("mail client not initialized. skipping template store initialization")
	}

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
				log.Logger.Panic().Err(err).Msg("unable to close file")
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
		return err
	}
	
	MailTemplateStore = &TemplateStore{templates: templateFiles}

	return nil
}

func (templateStore *TemplateStore) GetTemplate(name string) string {
	if len(templateStore.templates) == 0 {
		return ""
	}

	return templateStore.templates[name]
}

func (templateStore *TemplateStore) SetTemplate(name string, content string) {
	templateStore.templates[name] = content
}

func (templateStore *TemplateStore) RenderTemplate(name string, variables interface{}) string {
	if len(templateStore.templates) == 0 {
		return ""
	}

	rawTemplate, ok := templateStore.templates[name]

	if !ok {
		return ""
	}

	templateInstance, err := template.New(name).Parse(rawTemplate)

	if err != nil {
		return ""
	}

	buffer := new(bytes.Buffer)
	err2 := templateInstance.Execute(buffer, variables)

	if err2 != nil {
		return ""
	}

	return buffer.String()
}
