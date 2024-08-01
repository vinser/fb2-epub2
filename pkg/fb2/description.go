package fb2

import (
	"encoding/xml"
	"strings"

	"github.com/vinser/fb2-epub2/pkg/epub2"
)

func (p *FB2Parser) parseDescription(e *epub2.EPUB) error {
	for {
		token, err := p.Token()
		if err != nil {
			return err
		}

		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "title-info":
				if err = p.fillTitleInfo(e); err != nil {
					return err
				}
			}

		case xml.EndElement:
			if t.Name.Local == "description" {
				return nil
			}
		}
	}
}

func (p *FB2Parser) fillTitleInfo(e *epub2.EPUB) error {
	for {
		token, err := p.Token()
		if err != nil {
			return err
		}

		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "genre":
				genre, err := p.getText()
				if err != nil {
					return err
				}
				e.AddMetadataSubject(genre)

			case "author":
				firstName, middleName, lastname, err := p.getNames("author")
				if err != nil {
					return err
				}
				e.AddMetadataAuthor(firstName, middleName, lastname)

			case "book-title":
				title, err := p.getText()
				if err != nil {
					return err
				}
				e.AddMetadataTitle(title)

			case "lang":
				lang, err := p.getText()
				if err != nil {
					return err
				}
				e.AddMetadataLanguage(lang)

			case "coverpage":
				if err = p.addCoverPage(e); err != nil {
					return err
				}
			}

		case xml.EndElement:
			if t.Name.Local == "title-info" {
				return nil
			}
		}
	}
}

func (p *FB2Parser) getNames(tag string) (fn, mn, ln string, err error) {
	var token xml.Token

	for {
		if token, err = p.Token(); err != nil {
			return
		}

		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "first-name":
				if fn, err = p.getText(); err != nil {
					return
				}

			case "middle-name":
				if mn, err = p.getText(); err != nil {
					return
				}

			case "last-name":
				if ln, err = p.getText(); err != nil {
					return
				}
			}

		case xml.EndElement:
			if t.Name.Local == tag {
				return
			}
		}
	}
}

func (p *FB2Parser) addCoverPage(e *epub2.EPUB) error {
	for {
		token, err := p.Token()
		if err != nil {
			return err
		}

		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local == "image" {
				for _, a := range t.Attr {
					if a.Name.Local == "href" {
						imageName := a.Value
						if imageName[0] == '#' {
							imageName = strings.TrimLeft(imageName, "#")
						}
						if err = e.AddItem("cover", "cover", `<div class="cover"><img class="coverimage" alt="Cover" src="`+imageName+`" /></div>`); err != nil {
							return err
						}
						e.AddMetadataCover(imageName)
					}
				}
			}

		case xml.EndElement:
			if t.Name.Local == "coverpage" {
				return nil
			}
		}
	}
}
