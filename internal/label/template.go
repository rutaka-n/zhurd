package label

import (
	"bytes"
	"errors"
	"unicode"
	"unicode/utf8"
)

var (
	MissingPlaceholderError = errors.New("missing placeholder")
	DecodingError           = errors.New("template body decoding error")
)

const (
	separator = "^_"
)

type Template struct {
	ID      int64 `json:"id"`
	labelID int64
	Type    string `json:"type"`
	Body    []byte `json:"body"`
}

func NewTemplate(labelID int64, pType string, body []byte) (Template, error) {
	escapedBody, err := escapeBody(body)
	if err != nil {
		return Template{}, err
	}
	return Template{
		labelID: labelID,
		Type:    pType,
		Body:    escapedBody,
	}, nil
}

func (t Template) Print(placeholders map[string]string) ([]byte, error) {
	output := bytes.NewBuffer([]byte{})
	parts := bytes.Split(t.Body, []byte(separator))
	for _, part := range parts {
		if isPlaceholder(part) {
			val, ok := placeholders[string(part)]
			if !ok {
				return []byte{}, MissingPlaceholderError
			}
			if _, err := output.WriteString(val); err != nil {
				return []byte{}, err
			}
			continue
		}
		if _, err := output.Write(part); err != nil {
			return []byte{}, err
		}
	}
	return output.Bytes(), nil
}

func isAllowedSymbol(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func isPlaceholder(bs []byte) bool {
	return bytes.HasPrefix(bs, []byte("_")) &&
		bytes.HasSuffix(bs, []byte("_")) &&
		bytes.ContainsFunc(bs, isAllowedSymbol)
}

func escapeSeq(s []rune) ([]byte, error) {
	res := bytes.NewBuffer([]byte{})
	if _, err := res.WriteString(separator); err != nil {
		return []byte{}, err
	}
	for i := range s {
		if _, err := res.WriteRune(s[i]); err != nil {
			return []byte{}, err
		}
	}
	if _, err := res.WriteString(separator); err != nil {
		return []byte{}, err
	}
	return res.Bytes(), nil
}

func escapeBody(body []byte) ([]byte, error) {
	escapedBody := bytes.NewBuffer([]byte{})
	ch := []byte{}
	placeholder := []rune{}
	for _, b := range body {
		ch = append(ch, b)
		if !utf8.FullRune(ch) {
			continue
		}
		r, rlen := utf8.DecodeRune(ch)
		if rlen != len(ch) {
			return []byte{}, DecodingError
		}
		// placeholder started
		if len(placeholder) == 0 && r == '_' {
			placeholder = append(placeholder, r)
			ch = []byte{}
			continue
		}
		// in placeholder, add allowed symbols
		if len(placeholder) > 0 && (isAllowedSymbol(r)) {
			placeholder = append(placeholder, r)
			ch = []byte{}
			continue
		}
		// placeholder ended
		if len(placeholder) > 1 && placeholder[len(placeholder)-1] == '_' {
			bs, err := escapeSeq(placeholder)
			if err != nil {
				return []byte{}, err
			}
			if _, err := escapedBody.Write(bs); err != nil {
				return []byte{}, err
			}
		} else {
			// is not a valid placeholder, treat it as a plain content
			for i := range placeholder {
				if _, err := escapedBody.WriteRune(placeholder[i]); err != nil {
					return []byte{}, err
				}
			}
		}
		// reset placeholder
		placeholder = []rune{}
		// reset character
		ch = []byte{}

		if _, err := escapedBody.WriteRune(r); err != nil {
			return []byte{}, err
		}
	}

	return escapedBody.Bytes(), nil
}
