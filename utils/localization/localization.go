package localization

import (
	"9DotTechnology/trading_platform/constants/common_constants"
	"9DotTechnology/trading_platform/utils/logwrapper"
	"encoding/json"
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// Bundle - to log output on console
var Bundle *i18n.Bundle

var Languages = []string{"en"}

// LoadBundle local locales file
func LoadBundle() *i18n.Bundle {
	// Initialize i18n
	Bundle = i18n.NewBundle(language.Make(common_constants.DEFAULT_LANGUAGE))
	Bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	for _, lang := range Languages {
		logwrapper.Logger.Debugln("loading file ", fmt.Sprintf("locales/%v.json", lang))
		Bundle.MustLoadMessageFile(fmt.Sprintf("locales/%v.json", lang))
	}

	return Bundle
}

// GetMessage get message from local file
func GetMessage(lang interface{}, id string, templateData interface{}) string {
	language := common_constants.DEFAULT_LANGUAGE
	if lang != nil {
		language = lang.(string)
	}
	if pos(Languages, language) == -1 {
		language = common_constants.DEFAULT_LANGUAGE
	}
	localizer := i18n.NewLocalizer(Bundle, language)

	message, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: id,
		},
		TemplateData: templateData,
	})

	if err != nil || message == "" {
		message = id
	}

	return message
}

func pos(slice []string, value string) int {
	for p, v := range slice {
		if v == value {
			return p
		}
	}
	return -1
}
