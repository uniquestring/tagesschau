package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"syscall"
	"time"
)

type Video struct {
	Id    string
	Date  time.Time
	Title string
}

type video struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type Playlist struct {
	Entries []video `json:"entries"`
}

const (
	TIME_FORMAT          = "02.01.2006"
	TIME_FORMAT_NO_ZEROS = "2.1.2006"
	URL_PRFIX_VIDEO      = "https://www.youtube.com/watch?v="
	URL_PRFIX_PLAYLIST   = "https://www.youtube.com/playlist?list="
	VIDEO_TITLE_PATTERN  = "tagesschau 20:00 Uhr, (?P<Date>\\d{2}\\.\\d{2}.\\d{4})"
)

func usage() {
	fmt.Println("Usage: tagesschau [options] [url]")
	flag.PrintDefaults()
	fmt.Println("  <url>")
	fmt.Printf("\tIf an url is given, the linked video will be played. Any date options will be ignored. \n")
}

func main() {
	flag.Usage = usage
	sTargetDate := flag.String("d",
		time.Now().Format(TIME_FORMAT),
		"Target date of the video")
	offset := flag.Int("o", 0, "Date offset in days")

	flag.Parse()

	var url string
	var title string

	if flag.Arg(0) != "" {
		title = "<title unknown>"
		url = flag.Arg(0)
	} else {
		targetDate, err := getTargetDate(*sTargetDate, *offset)

		if err != nil {
			fmt.Print(err.Error())
			return
		}

		video, err := getVideoByDate(targetDate)

		if err != nil {
			fmt.Print(err.Error())
			return
		}

		title = video.Title
		url = URL_PRFIX_VIDEO + video.Id
	}

	fmt.Printf("Playing: %s - %s\n", title, url)

	if err := playVideo(url); err != nil {
		fmt.Printf("Error starting mpv:\n\t%s\n", err.Error())
	}
}

func getTargetDate(dateString string, offset int) (time.Time, error) {
	var err error
	var targetDate time.Time

	for _, format := range []string{TIME_FORMAT, TIME_FORMAT_NO_ZEROS} {
		targetDate, err = time.Parse(format, dateString)

		if err == nil {
			break
		}
	}

	if err != nil {
		return time.Now(), errors.New(
			fmt.Sprintf("Could not pare date \"%s\".", dateString))
	}

	targetDate = targetDate.AddDate(0, 0, offset)

	return targetDate, nil
}

func getVideoByDate(targetDate time.Time) (*Video, error) {
	// Try find the video in any of the playlists
	for _, playlist_id := range []string{PLAYLIST_ID, UPLOAD_PLAYLIST_ID} {
		videos, err := getPlaylistVideos(playlist_id)

		if err != nil || len(videos) == 0 {
			continue
		}

		for _, video := range videos {
			if video.Date.Equal(targetDate) {
				return &video, nil
			}
		}
	}

	return nil, errors.New(
		fmt.Sprintf("Could not find any matching video for %s\n",
			targetDate.Format(TIME_FORMAT),
		))
}

func getPlaylistVideos(playlistId string) ([]Video, error) {
	cmd := exec.Command("youtube-dl",
		"--playlist-start", "1",
		"--playlist-end", strconv.Itoa(LOOK_AROUND_AMOUNT),
		"--flat-playlist",
		"--dump-single-json",
		URL_PRFIX_PLAYLIST+playlistId,
	)

	out, err := cmd.Output()

	if err != nil {
		return []Video{}, err
	}

	var playlist Playlist
	err = json.Unmarshal(out, &playlist)

	if err != nil {
		return []Video{}, err
	}

	var videos []Video

	for _, vid := range playlist.Entries {
		vidParsed, err := parseVideo(vid)

		if err == nil && vidParsed != nil {
			videos = append(videos, *vidParsed)
		}
	}

	return videos, nil
}

func parseVideo(video video) (*Video, error) {
	matches := regexp.
		MustCompile(VIDEO_TITLE_PATTERN).
		FindStringSubmatch(video.Title)

	if len(matches) != 2 {
		return nil, errors.New("Could parse video date from title.")
	}

	date, err := time.Parse(TIME_FORMAT, matches[1])

	if err != nil {
		return nil, err
	}

	return &Video{
		Id:    video.Id,
		Title: video.Title,
		Date:  date}, nil
}

func playVideo(videoUrl string) error {
	mpv, err := exec.LookPath("mpv")

	if err != nil {
		return err
	}

	parameters := []string{"mpv"}
	parameters = append(parameters, MPV_PARAMETERS[:]...)
	parameters = append(parameters, videoUrl)

	return syscall.Exec(mpv,
		parameters,
		os.Environ())
}
