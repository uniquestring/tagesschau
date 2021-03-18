package main

const (
	// Dedicated playlist for the videos we are looking for
	PLAYLIST_ID = "PL4A2F331EE86DCC22"

	// Playlist of all uploads; fallback if they weren't added to the first
	// playlist
	UPLOAD_PLAYLIST_ID = "UU5NOEUbkLheQcaaRldYW5GA"

	// Limit how many videos should be taken from the playlist
	//
	// This is to prevent long loading and paring times, because there might be
	// thousands of videos on this channel/playlist.
	// Usually we'll only be looking for videos that are 0-3 days old. This look
	// around should be sufficient.
	LOOK_AROUND_AMOUNT = 10
)

var (
	// Additional parameters to pass to mpv
	MPV_PARAMETERS = [...]string{
		// Start at chapter 2, usually videos are tagged and chapter 1 is the intro.
		// If they are not yet tagged, mpv will ignore this and just start at the
		// beginning
		"--start=#2",

		// Watch at 2x speed with acceptable sound
		"--speed=2",
		"--af=scaletempo2",
	}
)
