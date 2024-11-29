package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"
	"time"
)

type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

var tracks = []*Track{
	{"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
	{"Go", "Moby", "Moby", 1992, length("3m37s")},
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, length("4m36s")},
	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, length("4m24s")},
}

func length(input string) time.Duration {
	d, err := time.ParseDuration(input)
	if err != nil {
		log.Fatal(err)
	}
	return d
}

func printTracks(tracks []*Track) {
	const format = "%v\t%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Title", "Artist", "Album", "Year", "Length")
	fmt.Fprintf(tw, format, "-----", "-----", "-----", "-----", "-----")
	for _, track := range tracks {
		fmt.Fprintf(tw, format, track.Title, track.Album, track.Album, track.Year, track.Length)
	}
	tw.Flush()
	fmt.Println()
}

type byArtist []*Track

func (t byArtist) Len() int {
	return len(t)
}

func (t byArtist) Less(i, j int) bool {
	return t[i].Artist < t[j].Artist
}

func (t byArtist) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

type customSort struct {
	t    []*Track
	less func(i, j *Track) bool
}

func (s *customSort) Len() int {
  return len(s.t)
}

func (s *customSort) Less(i, j int) bool {
  return s.less(s.t[i],s.t[j])
}

func (s *customSort) Swap(i, j int) {
  s.t[i], s.t[j] = s.t[j], s.t[i]
}

func main() {
	printTracks(tracks)
	sort.Sort(sort.Reverse(byArtist(tracks)))
	printTracks(tracks)
	sort.Sort(&customSort{
    t: tracks,
    less: func (i, j *Track) bool {
      if i.Title !=  j.Title {
        return i.Title < j.Title
      }
      if i.Year != i.Year {
        return i.Year < j.Year
      }
      return false
    },
  })
  printTracks(tracks)
}
