package structs

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type ScrapObject struct {
	Url   *url.URL
	name  string
	Sizes []string
}

func (s ScrapObject) String() string {
	return fmt.Sprintf("%s in %s at %s", s.name, s.Sizes, s.Url)
}

var sizesRegex = regexp.MustCompile("(X*S|M|X*L|[[:digit:]]+)")

func NewScrapObjectByCSVLine(line []string) (*ScrapObject, error) {
	if len(line) < 2 {
		return nil, fmt.Errorf("too less arguments to create object")
	}
	url, err := url.ParseRequestURI(line[0])
	if err != nil {
		return nil, fmt.Errorf("cannot parse URL, due: %w", err)
	}
	sizes := strings.Split(line[1], ",")
	res := &ScrapObject{
		Url:   url,
		Sizes: make([]string, 0, len(sizes)),
	}
	for _, s := range sizes {
		upper := strings.ToUpper(s)
		if sizesRegex.MatchString(upper) {
			res.Sizes = append(res.Sizes, upper)
		}
	}
	if len(res.Sizes) < 1 {
		return nil, fmt.Errorf("no size was supplied for object")
	}
	return res, nil
}

func (s ScrapObject) GetSite() *url.URL {
	return s.Url
}

func (s ScrapObject) GetSizes() []string {
	return s.Sizes
}
