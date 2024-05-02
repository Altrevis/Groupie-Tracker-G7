package service

import (
	"fmt"
	"../dto"
	"../API"
	"../repo"
	"time"
)

func Get() ([]dto.Artist, error) {
	artists, err := repository.GetArtists()
	if err != nil {
		return nil, err
	}
	dtoArtists, err := createDtos(artists)
	if err != nil {
		return nil, err
	}
	return dtoArtists, nil
}

func GetArtistById(id int) (*dto.Artist, error) {
	artist, err := repository.GetArtistById(id)
	if err != nil {
		return nil, err
	}
	fmt.Println(artist)
	dtoArtist, err := createDto(*artist)
	if err != nil {
		return nil, err
	}
	return dtoArtist, nil
}

func createDto(artist API.Artist) (*dto.Artist, error) {
	var err error
	dtoArtist := &dto.Artist{}
	dtoArtist.Id = artist.Id
	dtoArtist.CreationDate = artist.CreationDate
	dtoArtist.FirstAlbum = artist.FirstAlbum
	dtoArtist.Image = artist.Image
	dtoArtist.Members = artist.Members
	dtoArtist.Name = artist.Name
	dtoArtist.Relations, err = repository.GetRelationsFromArtist(artist.RelationsUrl)
	if err != nil {
		return nil, err
	}
	dtoArtist.Location, err = repository.GetLocationsFromArtist(artist.LocationsUrl)
	if err != nil {
		return nil, err
	}
	dtoArtist.ConcertDates, err = repository.GetConcertDatesFromArtist(artist.ConcertDatesUrl)
	if err != nil {
		return nil, err
	}

	return dtoArtist, nil
}

func parallel(a API.Artist, chanArt chan<- dto.Artist) {
	chanLoc := make(chan API.Location)
	chanDate := make(chan API.Date)
	chanRel := make(chan API.Relation)
	dtoArtist := dto.Artist{}
	dtoArtist.Id = a.Id
	dtoArtist.CreationDate = a.CreationDate
	dtoArtist.FirstAlbum = a.FirstAlbum
	dtoArtist.Image = a.Image
	dtoArtist.Members = a.Members
	dtoArtist.Name = a.Name
	go repository.GetLocationsFromArtistAsync(a.LocationsUrl, chanLoc)
	go repository.GetConcertDatesFromArtistAsync(a.ConcertDatesUrl, chanDate)
	go repository.GetRelationsFromArtistAsync(a.RelationsUrl, chanRel)
	n := 0
	for n != 3 {
		select {
		case loc := <-chanLoc:
			{
				dtoArtist.Location = loc
				n++
			}
		case rel := <-chanRel:
			{
				dtoArtist.Relations = rel
				n++
			}
		case date := <-chanDate:
			{
				dtoArtist.ConcertDates = date
				n++
			}
		}
	}

	chanArt <- dtoArtist
}

func createDtos(artists []API.Artist) ([]dto.Artist, error) {
	var dtoArtists []dto.Artist
	start := time.Now()
	chanArt := make(chan dto.Artist)
	for _, a := range artists {
		go parallel(a, chanArt)
	}
	for len(dtoArtists) != len(artists) {
		select {
		case elem := <-chanArt :
			{
				dtoArtists = append(dtoArtists, elem)
			}
		}
	}
	elapsed := time.Since(start)

	fmt.Printf("took %s \n", elapsed)

	return dtoArtists, nil
}