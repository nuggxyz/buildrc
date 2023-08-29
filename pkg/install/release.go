package install

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/walteh/buildrc/pkg/file"
	"golang.org/x/oauth2"
)

func InstallLatestGithubRelease(ctx context.Context, fls afero.Fs, org string, name string, token string) error {

	var err error

	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/repos/"+org+"/"+name+"/releases/latest", nil)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/vnd.github.v3+json")

	var client *http.Client

	if token != "" {
		client = oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
	} else {
		client = &http.Client{}
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zerolog.Ctx(ctx).Debug().Err(err).Msg("error reading body")
		return err
	}

	if resp.StatusCode != 200 {
		zerolog.Ctx(ctx).Debug().Err(err).RawJSON("response_body", body).Msg("bad status")
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	var release struct {
		Assets []struct {
			BrowserDownloadURL string `json:"browser_download_url"`
			Name               string `json:"name"`
		} `json:"assets"`
		URL string `json:"url"`
	}

	if err := json.Unmarshal(body, &release); err != nil {
		zerolog.Ctx(ctx).Debug().Err(err).RawJSON("response_body", body).Msg("error unmarshaling body")
		return err
	}

	zerolog.Ctx(ctx).Debug().Interface("respdata", release).Msg("got respdata")

	targetPlat := runtime.GOOS + "-" + runtime.GOARCH

	if os.Getenv("GOARM") != "" {
		targetPlat += "-" + os.Getenv("GOARM")
	}

	dl := ""

	for _, asset := range release.Assets {
		if strings.HasSuffix(asset.Name, targetPlat+".tar.gz") {
			dl = asset.BrowserDownloadURL
			break
		}
	}

	if dl == "" {
		return fmt.Errorf("no release found for %s", targetPlat)
	}

	fle, err := downloadFile(ctx, client, fls, dl)
	if err != nil {
		return err
	}

	defer fle.Close()

	// untar the release
	out, err := file.Untargz(ctx, fls, fle.Name())
	if err != nil {
		return err
	}

	// install the release
	err = InstallAs(ctx, fls, out.Name(), name)
	if err != nil {
		return err
	}

	return nil

}

func downloadFile(ctx context.Context, client *http.Client, fls afero.Fs, str string) (fle afero.File, err error) {

	base := filepath.Base(str)

	// Create the file
	out, err := afero.TempDir(fls, "", "")
	if err != nil {
		return nil, err
	}

	fle, err = fls.Create(filepath.Join(out, base))
	if err != nil {
		return nil, err
	}

	// Get the data
	resp, err := client.Get(str)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		zerolog.Ctx(ctx).Debug().Err(err).Str("file_name", str).Msg("bad status for GET to download file")
		return nil, fmt.Errorf("bad status for GET to download file: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(fle, resp.Body)
	if err != nil {
		return nil, err
	}

	return fle, nil
}
