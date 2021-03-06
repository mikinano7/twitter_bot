package service

import (
	"./google"
	"./itunes"
	"fmt"
	"github.com/mikinano7/xvideos4go"
	"github.com/mikinano7/dropbox4go"
	"github.com/spf13/viper"
	"math/rand"
	"time"
	"net/http"
	"path"
	"strings"
	"errors"
"github.com/mikinano7/soundcloud"
)

func DropboxUpload(url string) string {
	pos := strings.LastIndex(url, "/")
	fileName := url[pos + 1:]
	ext := path.Ext(fileName)

	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif":
		if (!strings.Contains(url, "http")) {
			return onError(errors.New("incorrect resource."))
		}
	default:
		return onError(errors.New("incorrect extension."))
	}

	token := viper.GetString("dropbox.access_token")

	httpClient := http.DefaultClient
	resp, err := httpClient.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return onError(err)
	}

	svc := dropbox4go.New(httpClient, token)
	req := dropbox4go.Request{
		File: resp.Body,
		Parameters: dropbox4go.Parameters{
			Path: "/twitter/" + fileName,
			Mode: "overwrite",
			AutoRename: false,
			ClientModified: time.Now().UTC().Format(time.RFC3339),
			Mute: true,
		},
	}

	result, err := svc.Upload(req)

	if err != nil {
		return onError(err)
	} else {
		return fmt.Sprintf("file %s has uploaded. (size: %d bytes)", fileName, result.Size)
	}
}

func ITunes(query []string) string {
	if res, err := itunes.Search(query); err != nil {
		return err.Error()
	} else {
		if len(res) > 0 {
			return fmt.Sprintf(
				"%s / %s - %s",
				res[0].TrackName,
				res[0].ArtistName,
				res[0].PreviewUrl,
			)
		} else {
			return fmt.Sprintf(
				"検索結果が0件でした。 [%s]",
				time.Now().String(),
			)
		}
	}
}

func YouTube(query []string) string {
	if res, err := google.YouTube(query); err != nil {
		return err.Error()
	} else {
		if len(res) > 0 {
			return fmt.Sprintf(
				"%s - %s%s",
				res[0].Snippet.Title,
				"https://www.youtube.com/watch?v=",
				res[0].Id.VideoId,
			)
		} else {
			return fmt.Sprintf(
				"検索結果が0件でした。 [%s]",
				time.Now().String(),
			)
		}
	}
}

func Xvideos(query []string) string {
	res := xvideos4go.Search(query)

	if len(res) > 0 {
		rand.Seed(time.Now().UnixNano())
		rand.Intn(len(res) - 1)

		return fmt.Sprintf(
			"%s%s - %s",
			res[0].Title,
			res[0].Duration,
			res[0].Url,
		)
	} else {
		return fmt.Sprintf(
			"検索結果が0件でした。 [%s]",
			time.Now().String(),
		)
	}
}

func Soundcloud(query []string) string {
	clientId := viper.GetString("soundcloud.client_id")
	clientSecret := viper.GetString("soundcloud.client_secret")
	api, _ := soundcloud.NewClient(clientId, clientSecret)
	res := api.GetTracks(strings.Join(query, "+"))

	if len(res) > 0 {
		rand.Seed(time.Now().UnixNano())
		rand.Intn(len(res) - 1)

		return fmt.Sprintf(
			"%s - %s",
			res[0].Title,
			res[0].PermalinkUrl,
		)
	} else {
		return fmt.Sprintf(
			"検索結果が0件でした。 [%s]",
			time.Now().String(),
		)
	}
}

func onError(err error) string {
	return fmt.Sprintf(
		"%s [%s]",
		err.Error(),
		time.Now().In(time.FixedZone("Asia/Tokyo", 9 * 60 * 60)).Format(time.RFC3339),
	)
}
