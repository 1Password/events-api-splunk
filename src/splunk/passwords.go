package splunk

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type PasswordsResponse struct {
	XMLName    xml.Name `xml:"feed"`
	Text       string   `xml:",chardata"`
	Xmlns      string   `xml:"xmlns,attr"`
	S          string   `xml:"s,attr"`
	Opensearch string   `xml:"opensearch,attr"`
	Title      string   `xml:"title"`
	ID         string   `xml:"id"`
	Updated    string   `xml:"updated"`
	Generator  struct {
		Text    string `xml:",chardata"`
		Build   string `xml:"build,attr"`
		Version string `xml:"version,attr"`
	} `xml:"generator"`
	Author struct {
		Text string `xml:",chardata"`
		Name string `xml:"name"`
	} `xml:"author"`
	Link []struct {
		Text string `xml:",chardata"`
		Href string `xml:"href,attr"`
		Rel  string `xml:"rel,attr"`
	} `xml:"link"`
	TotalResults string `xml:"totalResults"`
	ItemsPerPage string `xml:"itemsPerPage"`
	StartIndex   string `xml:"startIndex"`
	Messages     string `xml:"messages"`
	Entry        struct {
		Text    string `xml:",chardata"`
		Title   string `xml:"title"`
		ID      string `xml:"id"`
		Updated string `xml:"updated"`
		Link    []struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
		} `xml:"link"`
		Author struct {
			Text string `xml:",chardata"`
			Name string `xml:"name"`
		} `xml:"author"`
		Content struct {
			Text string `xml:",chardata"`
			Type string `xml:"type,attr"`
			Dict struct {
				Text string `xml:",chardata"`
				Key  []struct {
					Text string `xml:",chardata"`
					Name string `xml:"name,attr"`
					Dict struct {
						Text string `xml:",chardata"`
						Key  []struct {
							Text string `xml:",chardata"`
							Name string `xml:"name,attr"`
							Dict struct {
								Text string `xml:",chardata"`
								Key  []struct {
									Text string `xml:",chardata"`
									Name string `xml:"name,attr"`
									List struct {
										Text string   `xml:",chardata"`
										Item []string `xml:"item"`
									} `xml:"list"`
								} `xml:"key"`
							} `xml:"dict"`
							List struct {
								Text string `xml:",chardata"`
								Item string `xml:"item"`
							} `xml:"list"`
						} `xml:"key"`
					} `xml:"dict"`
				} `xml:"key"`
			} `xml:"dict"`
		} `xml:"content"`
	} `xml:"entry"`
}

func (s *SplunkAPI) GetPasswords(ctx context.Context, passwordKey, passwordRealm string) (*PasswordsResponse, error) {
	url := fmt.Sprintf("/servicesNS/nobody/onepassword_events_api/storage/passwords/%s:%s:", passwordRealm, passwordKey)
	res, err := s.request(ctx, "GET", url, nil)
	if err != nil {
		err := fmt.Errorf("could not make SplunkAPIRequest: %w", err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		err := fmt.Errorf("received a non 200 response: %d", res.StatusCode)
		return nil, err
	}

	passwordsResponse := &PasswordsResponse{}
	decoder := xml.NewDecoder(res.Body)
	err = decoder.Decode(passwordsResponse)
	if err != nil {
		err := fmt.Errorf("could not decode response: %w", err)
		return nil, err
	}

	return passwordsResponse, nil
}

func (s *SplunkAPI) CreatePassword(ctx context.Context, name, password, realm string) error {
	endpoint := "/servicesNS/nobody/onepassword_events_api/storage/passwords"
	data := url.Values{}
	data.Set("name", name)
	data.Set("password", password)
	data.Set("realm", realm)
	res, err := s.request(ctx, "POST", endpoint, data)
	if err != nil {
		err := fmt.Errorf("could not make SplunkAPIRequest: %w", err)
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusConflict {
		// A StatusConflict 409 if the password was already created
		b, _ := ioutil.ReadAll(res.Body)
		log.Println(string(b))
		err := fmt.Errorf("received a non 201 response: %d", res.StatusCode)
		return err
	}

	return nil
}
