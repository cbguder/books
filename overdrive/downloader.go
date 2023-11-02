package overdrive

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"

	"github.com/cbguder/books/api_client"
)

const numWorkers = 8

type Downloader struct {
	httpClient *http.Client
	apiClient  *api_client.ApiClient
}

func NewDownloader() *Downloader {
	jar, _ := cookiejar.New(nil)

	httpClient := &http.Client{
		Jar: jar,
	}

	apiClient := api_client.NewApiClientWithHttpClient(httpClient)

	return &Downloader{
		httpClient: httpClient,
		apiClient:  apiClient,
	}
}

func (d *Downloader) Download(ctx context.Context, loan *OpenLoanResponse, destFolder string) error {
	rosters, err := d.getRosters(ctx, loan.Urls.Rosters, loan.Message)
	if err != nil {
		return err
	}

	contentRoster, err := findContentRoster(rosters)
	if err != nil {
		return err
	}

	urls := []string{
		loan.Urls.Openbook,
	}

	for _, entry := range contentRoster.Entries {
		urls = append(urls, entry.Url)
	}

	fmt.Printf("Downloading %d files to \"%s\"...\n", len(urls), destFolder)

	err = os.MkdirAll(destFolder, 0755)
	if err != nil {
		return err
	}

	return d.downloadAll(ctx, destFolder, urls)
}

func (d *Downloader) getRosters(ctx context.Context, rostersUrl, message string) ([]Roster, error) {
	rostersUrl = rostersUrl + "?" + message

	req, err := http.NewRequestWithContext(ctx, "GET", rostersUrl, nil)
	if err != nil {
		return nil, err
	}

	var rosters []Roster
	err = d.apiClient.Do(req, &rosters)
	return rosters, err
}

func (d *Downloader) downloadAll(ctx context.Context, destFolder string, urls []string) error {
	jobs := make(chan string, len(urls))
	results := make(chan error, len(urls))

	for i := 0; i < numWorkers; i++ {
		go d.downloadWorker(ctx, destFolder, jobs, results)
	}

	for _, srcUrl := range urls {
		jobs <- srcUrl
	}

	close(jobs)

	var errors int

	for i := 0; i < len(urls); i++ {
		err := <-results
		if err != nil {
			errors += 1
		}
	}

	if errors > 0 {
		return fmt.Errorf("failed to download %d files", errors)
	}

	return nil
}

func (d *Downloader) downloadWorker(ctx context.Context, destFolder string, urls <-chan string, results chan<- error) {
	for srcUrl := range urls {
		results <- d.downloadToFile(ctx, destFolder, srcUrl)
	}
}

func (d *Downloader) downloadToFile(ctx context.Context, destFolder, srcUrl string) error {
	fname, err := filenameFromUrl(srcUrl)
	if err != nil {
		return err
	}

	fpath := filepath.Join(destFolder, fname)

	parentDir := filepath.Dir(fpath)
	err = os.MkdirAll(parentDir, 0755)
	if err != nil {
		return err
	}

	f, err := os.Create(fpath)
	if err != nil {
		return err
	}

	defer f.Close()

	resp, err := d.download(ctx, srcUrl)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}

func (d *Downloader) download(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp, err
}

func findContentRoster(rosters []Roster) (*Roster, error) {
	for _, roster := range rosters {
		if roster.Group == "title-content" {
			return &roster, nil
		}
	}

	return nil, fmt.Errorf("failed to find content roster")
}

func filenameFromUrl(entryUrl string) (string, error) {
	u, err := url.Parse(entryUrl)
	if err != nil {
		return "", err
	}

	return u.Path, nil
}
