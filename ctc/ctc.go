// Package CTC : coordinate to country
package CTC

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"gonum.org/v1/gonum/spatial/kdtree"
)
const (
	geoCodeURL = "https://download.geonames.org/export/dump/cities1000.zip"
	geoCodePath = "./ctc/data/geocode.csv"
)
var NotFoundCountry = errors.New("can't find country for coordinates")

type Geo struct {
	locations map[int]string
	coordinates [][]float64
	KDTree *kdtree.Tree
}

func New() *Geo {
	var geo Geo
	geo.loadCoordinate()
	geo.makeKDTree()
	return &geo
}

func (g *Geo) GetCountryByCoordinate(lat, lng float64) (country string, err error){
	a, _ := g.KDTree.Nearest(kdtree.Point{lat, lng})
	if a == nil {
		return "", NotFoundCountry
	}
	for idx, point := range g.coordinates {
		p, _ := a.(kdtree.Point)
		if point[0] == p[0] || point[1] == p[1]{
			return g.locations[idx], nil
		}
	}
	return "", NotFoundCountry
}

func (g *Geo) loadCoordinate() {
	f, err := os.Open(geoCodePath)
	if err != nil {
		fmt.Println(err)
		if err := g.downloadGeoCode(); err != nil {
			fmt.Println(err)
		}
	}
	r := csv.NewReader(f)

	location := make(map[int]string)
	idx := 0
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			break
		}
		x, _ := strconv.ParseFloat(row[0], 64)
		y, _ := strconv.ParseFloat(row[1], 64)
		g.coordinates = append(g.coordinates, []float64{x, y})
		location[idx] = row[2]
		idx ++
	}

	g.locations = location
}

func (g *Geo) downloadGeoCode() error {
	// TODO : If the file does not exist, download geocode file 
	return nil
}

func (g *Geo) makeKDTree() {
	var points []kdtree.Point
	for _, coord := range g.coordinates {
		points = append(points, kdtree.Point{coord[0], coord[1]})
	}

	p := make(kdtree.Points, len(points))
	p = points
	g.KDTree = kdtree.New(p, false)
}