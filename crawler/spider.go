package crawler

import (
	"bytes"
	"regexp"
	"sync"

	"github.com/gocolly/colly"
	"github.com/pnptcn/datura"
	"github.com/theapemachine/wrkspc/spd"
	"github.com/wrk-grp/errnie"
)

/*
Spider is a wrapper around colly.
*/
type Spider struct {
	url       string
	collector *colly.Collector
	visited   []string
	links     []string
	output    *spd.Datagram
	wg        *sync.WaitGroup
	store     datura.Store
}

/*
NewSpider creates a new Spider starting at the given URL.
*/
func NewSpider(url string) *Spider {
	return &Spider{
		url,
		colly.NewCollector(),
		[]string{},
		[]string{},
		nil,
		&sync.WaitGroup{},
		datura.NewStore(&datura.DataLake{}),
	}
}

/*
Configure the spider.
*/
func (s *Spider) Configure() *Spider {
	// setting a valid User-Agent header
	s.collector.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"
	s.collector.IgnoreRobotsTxt = true
	s.collector.DisallowedURLFilters = append(s.collector.DisallowedURLFilters, regexp.MustCompile(`m\.`))
	s.collector.DisallowedURLFilters = append(s.collector.DisallowedURLFilters, regexp.MustCompile(`#`))

	s.collector.OnRequest(func(r *colly.Request) {
		if s.output != nil {
			s.store.Write(s.output.Encode())
		}

		s.output = spd.New(
			spd.TXTHTML, spd.ARTIFACT, spd.SCRAPE, bytes.Join(
				[][]byte{
					[]byte(r.URL.Host),
					[]byte(r.URL.Path),
				}, []byte("/"),
			),
		)

		errnie.Debugs("visit", r.URL.String())
	})

	s.collector.OnError(func(_ *colly.Response, err error) {
		errnie.Handles(err)
	})

	s.collector.OnResponse(func(r *colly.Response) {
	})

	s.collector.OnHTML("*", func(e *colly.HTMLElement) {
		// Remove multiple spaces and newlines from the text.
		re := regexp.MustCompile(`\s+`)
		e.Text = re.ReplaceAllString(e.Text, " ")

		// Guard against empty text, or text with only spaces.
		if len(e.Text) == 0 || e.Text == " " {
			return
		}

		// Append the text to the output slice.
		s.output.Write([]byte(e.Text))
	})

	s.collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		// Guard against anchors.
		if link[0] == '#' {
			return
		}

		if normalizedLink, ok := s.normalizeLink(link); ok {
			// Check if the link is already in the visited slice.
			if !s.contains(s.visited, normalizedLink) {
				s.visited = append(s.visited, normalizedLink)
				e.Request.Visit(normalizedLink)
			}
		}
	})

	return s
}

/*
Run the spider.
*/
func (s *Spider) Run() error {
	errnie.Trace()
	s.collector.Visit(s.url)
	return nil
}
