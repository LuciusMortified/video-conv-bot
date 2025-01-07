package ent

import "text/template"

var ConvertTemplates map[ConvertStatus]*template.Template

func init() {
	ConvertTemplates = make(map[ConvertStatus]*template.Template)
	var err error

	ConvertTemplates[ConvertUnsupported], err = template.
		New("unsupported").
		Parse(`Неподдерживаемый файл`)
	check(err)

	ConvertTemplates[ConvertDownloading], err = template.
		New("downloading").
		Parse(`Скачиваю файл`)
	check(err)

	ConvertTemplates[ConvertConverting], err = template.
		New("converting").
		Parse(`Конвертирую файл в mp4`)
	check(err)

	ConvertTemplates[ConvertDone], err = template.
		New("done").
		Parse(`Загружаю файл`)
	check(err)

	ConvertTemplates[ConvertError], err = template.
		New("error").
		Parse(`Ошибка! {{.Error}}`)
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
