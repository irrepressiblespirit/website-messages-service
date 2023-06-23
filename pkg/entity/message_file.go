package entity

import (
	"strings"
)

type MessageFile struct {
	Message  `bson:",inline"`
	Addition MessageFileAddition `json:"addition" bson:"addition"`
}

type MessageFileAddition struct {
	Size     uint64 `json:"size" bson:"size"`
	FileName string `json:"filename" bson:"filename"`
	S3Path   string `json:"s3path" bson:"s3path"`
	MimeType string `json:"mimetype" bson:"mimetype"`
}

func (m *MessageFile) ParseAddition(addition map[string]interface{}) error {
	if val, ok := addition["size"]; ok {
		m.Addition.Size = uint64(val.(float64))
	}
	if val, ok := addition["filename"]; ok {
		m.Addition.FileName = val.(string)
	}
	if val, ok := addition["s3path"]; ok {
		m.Addition.S3Path = val.(string)
	}
	if val, ok := addition["mimetype"]; ok {
		m.Addition.MimeType = val.(string)
	}
	return nil
}

func (m *MessageFile) GetAddition() map[string]interface{} {
	addition := make(map[string]interface{})
	addition["size"] = m.Addition.Size
	addition["filename"] = m.Addition.FileName
	addition["s3path"] = m.Addition.S3Path
	addition["mimetype"] = m.Addition.MimeType
	return addition
}

func (m *MessageFile) IsCorrect() error {
	fieldError := make(map[string]string)
	if m.Body != "" {
		if len(strings.TrimSpace(m.Body)) == 0 {
			fieldError["body"] = "errors.messages.body.is.empty"
		}
	}
	if m.Addition.Size == 0 {
		fieldError["addition.size"] = "errors.messages.size.is.not.found"
	}
	if len(m.Addition.FileName) < 3 {
		fieldError["addition.fileName"] = "errors.messages.file.name.is.not.found"
	}
	if len(m.Addition.S3Path) < 10 {
		fieldError["addition.s3path"] = "errors.messages.s3path.is.not.found"
	}
	if !IsAvaliableMimeType(m.Addition.MimeType) {
		fieldError["addition.mimeType"] = "errors.messages.mime.type.is.not.found"
	}
	if len(fieldError) > 0 {
		return MessageError{
			StatusCode:  3,
			Message:     "errors.messages.validation",
			WrongFields: fieldError,
		}
	}
	return nil
}

func IsAvaliableMimeType(mime string) bool {
	avaliableMimeType := []string{
		"audio/mp4",
		"audio/aac",
		"audio/mpeg",
		"audio/ogg",
		"audio/vnd.wave",
		"audio/x-aac",
		"audio/x-flac",
		"audio/speex",
		"audio/vorbis",
		"image/gif",
		"image/jpeg",
		"image/pjpeg",
		"image/png",
		"image/svg+xml",
		"image/tiff",
		"image/vnd.microsoft.icon",
		"image/webp",
		"text/csv",
		"video/mpeg",
		"video/mp4",
		"video/ogg",
		"video/quicktime",
		"video/webm",
		"video/x-ms-wmv",
		"video/x-flv",
		"video/x-msvideo",
		"video/3gpp",
		"video/3gpp2",
		"application/vnd.oasis.opendocument.text",
		"application/vnd.oasis.opendocument.spreadsheet",
		"application/vnd.oasis.opendocument.presentation",
		"application/vnd.oasis.opendocument.graphics",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"application/vnd.ms-powerpoint",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/pdf",
		"application/zip",
		"application/xml",
		"application/msword",
	}
	for _, val := range avaliableMimeType {
		if val == mime {
			return true
		}
	}
	return false
}
