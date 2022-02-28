package geometry

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)
type Point struct {
	X, Y  float64
	Label string
}

func (p *Point) Stringer() string {
	return fmt.Sprintf("%s: %.3f,%.3f", p.Label, p.X, p.Y)
}

func (p *Point) Dist(q Point) float64{
	return math.Sqrt((p.X-q.X)*(p.X-q.X)+(p.Y-q.Y)*(p.Y-q.Y))
}

type Path struct {
	Waypoints []Point
	Label string
}

func (p *Path) AziStart() float64 {
	p0 := p.Waypoints[0]
	p1 := p.Waypoints[1]
	// corr:=math.Cos(p0.y)
	return 90.0 - math.Atan2((p1.Y-p0.Y), p1.X-p0.X)*180/math.Pi
}

func (p *Path) Length() float64{
	d:=0.0
	for i:=0;i<len(p.Waypoints)-1;i++{
		d+= p.Waypoints[i].Dist(p.Waypoints[i+1])
	}
	return d
}

type ScreenTransform struct {
	Minx, Maxx, Miny, Maxy float64
	Scale, Xc, Yc, W, H   float64
}

func (t *ScreenTransform) Stringer() string {
	return fmt.Sprintf("Transform: scale=%.1f centre(%.3f,%.3f) for screen %.1f x %.1f", t.Scale, t.Xc, t.Yc, t.W, t.H)
}

func (t *ScreenTransform) NewWindowSize(w, h float64) {
	t.W = w
	t.H = h
}

func (t *ScreenTransform) ToWorld(x, y float64) (a, b float64) {
	a = x - t.W/2
	a /= t.Scale
	a += t.Xc
	b = y - t.H/2
	b /= (-t.Scale)
	b += t.Yc
	// println(fmt.Sprintf("%+v",t))
	return a, b
}
func (t *ScreenTransform) ToScreen(x, y float64) (a, b float64) {
	a = x - t.Xc
	a *= t.Scale
	a += t.W / 2
	b = y - t.Yc
	b *= (-t.Scale)
	b += t.H / 2
	// println(fmt.Sprintf("%+v",t))
	return a, b
}

// converts a DMS P1/90-style DMS string to decimal
func decFromDMS(text string) (float64, error) {
	var mult float64 = 1
	ns := text[len(text)-1]
	if strings.ToLower(string(ns)) == "s" {
		mult = -1
	}
	if strings.ToLower(string(ns)) == "w" {
		mult = -1
	}
	s, err := strconv.ParseFloat(text[len(text)-6:len(text)-1], 64)
	if err != nil {
		return 0, errors.New("Bad seconds value")
	}
	m, err := strconv.Atoi(text[len(text)-8 : len(text)-6])
	if err != nil {
		return 0, errors.New("Bad minutes value")
	}
	d, err := strconv.Atoi(text[0 : len(text)-8])
	if err != nil {
		return 0, errors.New("Bad degrees value")
	}
	s /= 60.0
	s += float64(m)
	s /= 60.0
	s += float64(d)
	return mult * s, nil
}

func dmFromDecDegree(decdeg float64) string {
	// sign:=decdeg<0
	a := math.Abs(decdeg)
	d := math.Trunc(a)
	m := a - d
	m *= 60
	return fmt.Sprintf("%.0f° %.4f'", d, m)
}

func SphericalDistance(p1,p2 Point,r float64) float64{
	a:=math.Sin((p2.Y/180*math.Pi-p1.Y/180*math.Pi)/2.0)
	a*=a
	b:=math.Sin(p2.X/180*math.Pi-p1.X/180*math.Pi)+a
	c:=math.Cos(p1.Y/180*math.Pi)*math.Cos(p2.Y/180*math.Pi)*b*b
	c=math.Sqrt(c)
	c=math.Asin(c)
	c*=2*r
	return c
}