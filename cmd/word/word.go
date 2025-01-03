package word

import (
	"github.com/google/uuid"
	"github.com/lukasjarosch/go-docx"
	"github.com/rs/zerolog/log"
)

func Generate(m docx.PlaceholderMap) (string, error){
        // replaceMap is a key-value map whereas the keys
	// represent the placeholders without the delimiters
        // read and parse the template docx
	doc, err := docx.Open("static/docx/pattern.docx")
	if err != nil {
		log.Err(err).Msg("Error opening docx template file")
		return "", err
	}

    // replace the keys with values from replaceMap
	err = doc.ReplaceAll(m)
	if err != nil {
		log.Err(err).Msg("Error replacing placeholders in docx template file")
		return "",err
	}

	name := uuid.New().String() + ".docx"
    // write out a new file
	err = doc.WriteToFile("static/docx/" + name)
	if err != nil {
		log.Err(err).Msg("Error writing docx file")
		return "",err
	}
	return name, nil
}