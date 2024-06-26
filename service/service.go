package service

import (
	"fmt"
	"time"
	"../dto"
	"../API"
	"../repo"
)

// Get récupère la liste des artistes avec leurs détails depuis le repository
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

// GetArtistById récupère un artiste spécifique par son ID avec ses détails depuis le repository
func GetArtistById(id int) (*dto.Artist, error) {
	artist, err := repository.GetArtistById(id)
	if err != nil {
		return nil, err
	}
	dtoArtist, err := createDto(*artist)
	if err != nil {
		return nil, err
	}
	return dtoArtist, nil
}

// createDto crée un DTO (Data Transfer Object) à partir d'un artiste API avec ses détails
func createDto(artist API.Artist) (*dto.Artist, error) {
	var err error
	dtoArtist := &dto.Artist{
		Id:           artist.Id,
		CreationDate: artist.CreationDate,
		FirstAlbum:   artist.FirstAlbum,
		Image:        artist.Image,
		Members:      artist.Members,
		Name:         artist.Name,
	}

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

// parallel récupère de manière asynchrone les détails de localisation, dates de concert et relations pour un artiste
func parallel(a API.Artist, chanArt chan<- dto.Artist) {
	chanLoc := make(chan API.Location)
	chanDate := make(chan API.Date)
	chanRel := make(chan API.Relation)
	dtoArtist := dto.Artist{
		Id:           a.Id,
		CreationDate: a.CreationDate,
		FirstAlbum:   a.FirstAlbum,
		Image:        a.Image,
		Members:      a.Members,
		Name:         a.Name,
	}
	go repository.GetLocationsFromArtistAsync(a.LocationsUrl, chanLoc)
	go repository.GetConcertDatesFromArtistAsync(a.ConcertDatesUrl, chanDate)
	go repository.GetRelationsFromArtistAsync(a.RelationsUrl, chanRel)
	n := 0
	for n != 3 {
		select {
		case loc := <-chanLoc:
			dtoArtist.Location = loc
			n++
		case rel := <-chanRel:
			dtoArtist.Relations = rel
			n++
		case date := <-chanDate:
			dtoArtist.ConcertDates = date
			n++
		}
	}

	chanArt <- dtoArtist
}

// createDtos crée des DTO pour tous les artistes donnés
func createDtos(artists []API.Artist) ([]dto.Artist, error) {
	var dtoArtists []dto.Artist
	start := time.Now()
	chanArt := make(chan dto.Artist)
	for _, a := range artists {
		go parallel(a, chanArt)
	}
	for len(dtoArtists) != len(artists) {
		select {
		case elem := <-chanArt:
			dtoArtists = append(dtoArtists, elem)
		}
	}
	elapsed := time.Since(start)
	fmt.Printf("took %s \n", elapsed)
	return dtoArtists, nil
}
