package extractors

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/dop251/goja"
	"github.com/minna-sync/neptune/config"
)

var (
	ScriptRegex = regexp.MustCompile(`(?s)(eval)(\(f.*?)(</script>)`)
	SourceRegex = regexp.MustCompile(`https.*?m3u8`)
)

var (
	ErrFailedToExtractMediaURL = errors.New("failed to extract media url")
)

type KwikReferrerRoundTripper struct {
	http.RoundTripper
}

func (t *KwikReferrerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	clone := req.Clone(req.Context())

	headers := map[string]string{
		"Referer": "https://" + config.C.ProviderURLs.Animepahe,
	}

	for k, v := range headers {
		clone.Header.Set(k, v)
	}

	return t.RoundTripper.RoundTrip(clone)
}

type KwikExtractor struct {
	client *http.Client
	vm     *goja.Runtime
}

func NewKwikExtractor() *KwikExtractor {
	vm := goja.New()
	client := &http.Client{
		Transport: &KwikReferrerRoundTripper{http.DefaultTransport},
	}

	return &KwikExtractor{client: client, vm: vm}
}

func (k *KwikExtractor) unpack(content io.ReadCloser) (*string, error) {
	body, err := io.ReadAll(content)
	if err != nil {
		return nil, err
	}

	html := string(body)
	matches := ScriptRegex.FindStringSubmatch(html)
	if len(matches) < 3 {
		return nil, ErrFailedToExtractMediaURL
	}

	function := matches[2]
	function = strings.Replace(function, "eval", "", 1)

	unpacked, err := k.vm.RunString(function)
	if err != nil {
		return nil, err
	}

	src := SourceRegex.FindString(unpacked.String())
	if src == "" {
		return nil, ErrFailedToExtractMediaURL
	}

	return &src, nil
}

// Extract will exact the media url.
func (k *KwikExtractor) Extract(ctx context.Context, urlStr string) (*string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}

	resp, err := k.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		return nil, errors.New("content type is not text/html")
	}

	mediaSrc, err := k.unpack(resp.Body)
	if err != nil {
		return nil, err
	}

	return mediaSrc, nil
}
