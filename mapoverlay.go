package main

import (
	"errors"
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"log"
	"os"
	""
)

type loc struct {
	lat float64
	lng float64
}


// TODO: 1) Get Coords 					[X]
// TODO: 2) Get Map    					[ ] //kinda done, just call loc2map
// TODO: 3) Superimpose Map over image 	[ ]
// TODO:  a) pick package				[ ]
// TODO:  b) use package 				[ ]
func main() {
	goodFile := "DJI_0289.JPG"
	badFile := "doit.jpg"
	loc1, err1 := getCoords(goodFile)
	fmt.Printf("File 1 results - Location: %v \n Error: %v", loc1, err1)
	loc2, err2 := getCoords(badFile)
	fmt.Printf("File 2 results - Location: %v \n Error: %v", loc2, err2)
}

func buildCompositeImage(fileName string) error {
	loc, err := getCoords(fileName)
	if err != nil {
		if loc == nil {
			// location is nil, we can not recover
			log.Fatal(err)
		}
		// location is not nil, might be recoverable
		log.Print(err)
	}

	// the following TODOs assume use of https://github.com/disintegration/imaging
	//TODO: 1: call github.com/forrest321/loc2map method with coords to get map
	//TODO: 2: create new file (do not modify incoming, original file) with original's dimensions
	//TODO: 3: paste files in, creating new file

	return nil
}

func getCoords(fileName string) (*loc, error) {
	if fileName == "" {
		return nil, errors.New("fileName required")
	}
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	x, err := exif.Decode(f)
	if err != nil {
		return nil, err
	}
	lat, long, err := x.LatLong()
	if err != nil {
		return nil, err
	}
	return &loc{lat: lat,lng: long}, f.Close()
}
