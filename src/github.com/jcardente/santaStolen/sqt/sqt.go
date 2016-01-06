package sqt

import (
	"fmt"
	"os"
	"math"
//      "github.com/jcardente/santaStolen/types/location"
//      "github.com/jcardente/santaStolen/types/gift"
//      "github.com/jcardente/santaStolen/types/trip"
)


var _sqrt3 float64;

func init() {
	_sqrt3 = math.Sqrt(3)
}

// SQT ------------------------------------------------------------

// NB: This implements a very rudimentary spherical quadtree. It
//     assumes all nodes are added before the parent face triangles
//     are subdivided. It also does not support a general lookup
//     up mechanism. Instead, it only groups entries into corresponding
//     triangles and provides those groups on request.
//
//     The only thing that makes this resemble a spherical quadtree
//     is the mapping from lat/lon to triangle XY coordinates based
//     on XXXX


type SphereQuadTree struct {
	Triangles []*Triangle
	MaxLen    float64
	Splitted  bool
}


func NewSQT() *SphereQuadTree {	
	faces  := make([]*Triangle, 8);
        MaxLen := 1024.0
	
	for i:=0; i < 4; i++ {
		faces[i*2]   = NewTriangle(0,0, MaxLen, true)
		faces[i*2+1] = NewTriangle(0,0, MaxLen, false)
	}
	
	return &SphereQuadTree{faces, MaxLen, false}
}



func (s *SphereQuadTree) AddNode(id int, w float64, lat float64, lon float64) {

	if s.Splitted {
		fmt.Println("Error: adding node after splitting")
		os.Exit(1)
	}

	faceIdx := whichFace(lat, lon)

	// Determine relative lat/lon
        dlat := deg2rad(math.Abs(lat))
	dlon := deg2rad(math.Mod(lon, 90.0)) // - float64((int(math.Abs(lon)) % 90) * 90)

	// Determine XY coordinates for lat/long in RADIANS
	x, y := triangleProject(dlat, dlon, s.MaxLen)
	
	// Create node
	n := NewNode(x,y,id,w)
	
	// Add it appropriate face
	s.Triangles[faceIdx].AddNode(n)	
}


func deg2rad(d float64) float64 {
	return (d/360)*2*math.Pi
}

func whichFace(lat float64, lon float64) int {

	// Determine face
	faceIdx := 0
	if lat < 0 {
		faceIdx += 4
	}
	faceIdx += int(math.Mod(lon+360.0, 360.0)/90) 

	return faceIdx
}


func triangleProject(lat float64, lon float64, maxlen float64) (float64, float64) {
	x := maxlen/math.Pi * (lat + 2*lon*(1-2/math.Pi*lat));
	y := maxlen * _sqrt3 /math.Pi * lat;
    
	return x,y;
}


func (s *SphereQuadTree) Split(cb func(t *Triangle) bool) {

	s.Splitted = true
	
	// Recurisevly split until there's nothing left to do.
	anySplit := true
	for anySplit == true {
		newTris  := []*Triangle{}	
		anySplit = false
		for _, tri := range s.Triangles {			
			if cb(tri) {
				st := tri.Split()
				newTris = append(newTris, st...)
				anySplit = true
			} else {
				newTris = append(newTris, tri)
			}

		}
		s.Triangles = newTris			
	}
}


// NODE ------------------------------------------------------------


type Node struct {
	X float64
	Y float64
	Id int
	Weight float64 // caching here to avoid lookups
}

func NewNode(x float64, y float64, id int, w float64) *Node {
	return &Node{x,y,id, w}
}


// TRIANGLE ------------------------------------------------------------

type Triangle struct {
	X       float64
	Y       float64
	Len     float64
	Upward  bool
	Weight  float64
	Nodes   []*Node
}

func NewTriangle(x float64, y float64, len float64, upward bool) *Triangle {
	return &Triangle{x,y,len, upward, 0.0, nil}
}


func (t *Triangle) AddNode(n *Node) {
	t.Nodes   = append(t.Nodes, n)
	t.Weight += n.Weight;
}


func (t *Triangle) NumNodes() int {
	return len(t.Nodes)
}

func (t *Triangle) WeightNodes() float64 {
  return t.Weight
}


func (t *Triangle) Split() []*Triangle {
	dircoef := 1.0
	if !t.Upward {
		dircoef = -1.0
	}

	
	halfLen := t.Len/2
	Subs    := make([]*Triangle, 4);
	Subs[0] = NewTriangle(t.X, t.Y, halfLen, t.Upward)
	Subs[1] = NewTriangle(t.X + halfLen,   t.Y, halfLen, t.Upward)
	Subs[2] = NewTriangle(t.X + halfLen/2, t.Y + _sqrt3/2*halfLen*dircoef, halfLen, t.Upward)
	Subs[3] = NewTriangle(t.X + halfLen/2, t.Y + _sqrt3/2*halfLen*dircoef, halfLen, !t.Upward)

        for _, n := range t.Nodes {
		dx := math.Abs(n.X - t.X)
		dy := math.Abs(n.Y - t.Y)

		if (dy > _sqrt3*halfLen/2) {
			Subs[2].AddNode(n)
		} else if (dy < (_sqrt3*(t.Len - 2*dx)/2)) {

			Subs[0].AddNode(n)
		} else if (dy < (_sqrt3*(2*dx - t.Len)/2)) {
			Subs[1].AddNode(n)
		} else {
			Subs[3].AddNode(n)
		}
	}

  return Subs
}
