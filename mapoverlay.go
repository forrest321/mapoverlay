package main

import (
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"github.com/forrest321/loc2map"
)

type loc struct {
	lat float64
	lng float64
}


// TODO: 1) Get Coords 					[X]
// TODO: 2) Get Map    					[X] //kinda done, just call loc2map
// TODO: 3) Copy Original				[X]
// TODO: 4) Superimpose Map over image 	[X]
// TODO: 5) Check work from 3 & 4		[ ]

func main() {
	goodFile := "DJI_0289.JPG"
	badFile := "doit.jpg"
	loc1, err1 := getCoords(goodFile)
	fmt.Printf("File 1 results - Location: %v \n Error: %v", loc1, err1)
	loc2, err2 := getCoords(badFile)
	fmt.Printf("File 2 results - Location: %v \n Error: %v", loc2, err2)
}

func buildCompositeImage(fileName string) error {
	sourceFileStat, err := os.Stat(fileName)
	if err != nil {
		return err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", fileName)
	}

	loc, err := getCoords(fileName)
	if err != nil {
		if loc == nil {
			// location is nil, we can not recover
			log.Fatal(err)
		}
		// location is not nil, might be recoverable
		// if we got here, there was prob a failure in the file.Close() which can prob recover
		log.Print(err)
	}

	src, err := imaging.Open(fileName)
	if err != nil {
		return err
	}

	// the following TODOs assume use of https://github.com/disintegration/imaging
	//TODO: [X] 1: call github.com/forrest321/loc2map method with coords to get map
	mapImage, err := loc2map.Loc2ByteArrayOfMapImage(loc.lat, loc.lng)
	if err != nil {
		return err
	}

	// Write the image data to file
	if err := ioutil.WriteFile("map-tmp.png", mapImage, 0644); err != nil {
		log.Fatal(err)
	}

	//TODO: [X] 2: create new file (do not modify incoming, original file) with original's dimensions

	mapImg, err := imaging.Open("map-tmp.png")
	if err != nil {
		return err
	}
	mapImg = imaging.Resize(mapImg, src.Bounds().Max.X / 2, src.Bounds().Max.Y / 2, imaging.Lanczos)

	//TODO: [X] 3: create new image, paste files in
	newImg := imaging.New(src.Bounds().Max.X, src.Bounds().Max.Y, color.NRGBA{})
	newImg = imaging.Paste(newImg, src, image.Pt(0,0))
	newImg = imaging.Paste(newImg, mapImg, image.Pt(0,0))
	return imaging.Save(newImg, "output.png", imaging.PNGCompressionLevel(1))
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
	_ = f.Close() //close explicitly instead of waiting for defer. could also log possible error

	lat, long, err := x.LatLong()
	if err != nil {
		return nil, err
	}
	return &loc{lat: lat,lng: long}, nil
}
