package sqt

import (
	"testing"
	"math"
)

func TestAlloc(t *testing.T) {

	s := NewSQT()
	if len(s.Triangles) != 8 {
		t.Fatal("Incorrect number of faces")	
	}
}



func round(f float64, places int) float64 {
	pow := math.Pow(10, float64(places))
	return math.Floor(f * pow + 0.5) / pow
}


func TestTriProject(t *testing.T) {

	tests := [][]float64{
		{45,45,0.5,0.433},
		{0,90,1,0},
		{90,0,0.5,0.866},
		{0,0,0,0},
		{56.60,32.00,0.446,0.545},
		{55.80,64.50,0.582,0.537},
		{46.40,70.20,0.636,0.446},
		{31.30,82.50,0.772,0.301},
		{45.50,69.70,0.636,0.438},
		{23.90,48.00,0.524,0.23},
		{26.00,17.40,0.282,0.25},
		{64.50,48.90,0.512,0.621},
		{74.00,22.40,0.455,0.712},
		{48.40,30.80,0.427,0.466},
		{13.40,35.60,0.411,0.129},
		{38.20,48.60,0.523,0.368},
		{19.40,48.90,0.534,0.187},
		{63.60,65.40,0.566,0.612}}

	for _, test  := range tests {
		x,y := triangleProject(deg2rad(test[0]), deg2rad(test[1]),1.0)
		if (round(x,3) != test[2]) || (round(y,3) != test[3]) {
			t.Fatal("Failed test case, ", test, " -- ", x, ",",y)
		}
	}
}


func TestTriSplit(t *testing.T) {
	tri := NewTriangle(0,0,1,true)

	s := tri.Split()

	if len(s) != 4 {
		t.Fatal("Incorrect number of sub triangles, " , len(s))
	}
	
        if (s[0].X != 0) || (s[0].Y !=0) || (s[0].Len != 0.5) || (!s[0].Upward) {
		t.Fatal("S0 incorrect vertices, ", s[0])
	}

        if (s[1].X != 0.5) || (s[1].Y !=0) || (s[1].Len != 0.5) || (!s[1].Upward) {
		t.Fatal("S1 incorrect vertices, ", s[1])
	}

        if (s[2].X != 0.25) || (s[2].Y !=math.Sqrt(3)/4) || (s[2].Len != 0.5) || (!s[2].Upward) {
		t.Fatal("S2 incorrect vertices, ", s[2])
	}
        if (s[3].X != 0.25) || (s[3].Y !=math.Sqrt(3)/4) || (s[3].Len != 0.5) || (s[3].Upward) {
		t.Fatal("S3 incorrect vertices, ", s[3])
	}

	
	s2 := s[3].Split()
	if len(s2) != 4 {
		t.Fatal("Incorrect number of sub triangles, " , len(s))
	}
	
        if (s2[0].X != 0.25) || (s2[0].Y !=math.Sqrt(3)/4) || (s2[0].Len != 0.25) || (s2[0].Upward) {
		t.Fatal("S0 incorrect vertices, ", s2[0])
	}

        if (s2[1].X != 0.5) || (s2[1].Y !=math.Sqrt(3)/4) || (s2[1].Len != 0.25) || (s2[1].Upward) {
		t.Fatal("S1 incorrect vertices, ", s2[1])
	}

        if (s2[2].X != 0.375) || (s2[2].Y !=math.Sqrt(3)/8) || (s2[2].Len != 0.25) || (s2[2].Upward) {
		t.Fatal("S2 incorrect vertices, ", s2[2])
	}
        if (s2[3].X != 0.375) || (s2[3].Y !=math.Sqrt(3)/8) || (s2[3].Len != 0.25) || (!s2[3].Upward) {
		t.Fatal("S3 incorrect vertices, ", s2[3])
	}	
}


func TestTriAddNode(t *testing.T) {

	tri := NewTriangle(0,0,1,true)

	if tri.NumNodes() != 0 {
		t.Fatal("Incorrect number of nodes after alloc")
	}

	n := NewNode(0,0,1,1.0)

	tri.AddNode(n)

	if tri.NumNodes() != 1 {
		t.Fatal("Incorrect number of nodes after addition")
	}
	
}

func TestWhichFace(t *testing.T) {

	for i:= 0; i < 8; i++ {
		lat:= 45.0
		if (i >= 4) {
			lat = lat * -1
		}
		lon := 45.0 + float64(i%4) * 90.0;
		
		fidx := whichFace(lat, lon)
		if fidx != i {
			t.Fatal("Incorrect face for ", lat , ":", lon , "  i:", i, " fidx:", fidx)
		}
	}

	for i:= 0; i < 8; i++ {
		lat:= 45.0
		if (i >= 4) {
			lat = lat * -1
		}
		lon := -1.0 * (45.0 + float64(i%4) * 90.0);
		
		fidx := whichFace(lat, lon)
		expect := 3 + (int(i/4)*8) - i
		if fidx != expect  {
			t.Fatal("Incorrect face for ", lat , ":", lon , "  i:", i, " fidx:", fidx, " exp:", expect)
		}
	}	
}


func TestSphereAddNode(t *testing.T) {
	s := NewSQT()
	for i:=0; i <4; i++ {
		s.AddNode(i,     1.0, 45.0, 45.0+float64(i)*90.0,)
		s.AddNode(1+i*4, 1.0, -45.0, 45.0+float64(i)*90.0)
	}

	for i,ft := range s.Triangles {
		if (ft.NumNodes() ==0 ) {
			t.Fatal("Face didn't get a node ", i)
		}
	}

	
}


func TestTriNodeSplit(t *testing.T) {

	h := math.Sqrt(3)/2
	
	coords := [][]float64{
		{0.25,h/4},
		{0.5, h/4},
		{0.75, h/4},
		{0.5, 3*h/4}}

	
	tri := NewTriangle(0,0,1,true)

	for i, xy := range coords {
		n := NewNode(xy[0],xy[1],i,1.0)
		tri.AddNode(n)
	}

	if tri.NumNodes() != 4 {
		t.Fatal("Incorrect number of nodes before split")
	}

	newTries := tri.Split()

	if len(newTries) != 4 {
		t.Fatal("Incorrect number of triangles after split")
	}
	
	for _, nt := range newTries {
		if nt.NumNodes() != 1 {
			t.Fatal("Incorrect number of nodes after split")
		}		
	}
}


func TestSphereNodeSplit(t *testing.T) {

	coords := [][]float64{
		{22.5, 22.5},
		{22.5, 45.0},
		{22.5, 67.5},
		{67.5, 45.0}}


	s := NewSQT()
	for q := 0; q < 4; q++ {
		for i, ll := range coords {
			s.AddNode(q*len(coords)+i*2,1.0, ll[0], ll[1] + float64(q)*90.0)
			s.AddNode(q*len(coords)+i*2+1,1.0, ll[0] * -1.0, ll[1] + float64(q)*90.0)			
		}
	}
	
	for i,ft := range s.Triangles {
		if ft.NumNodes() !=4  {
			t.Fatal("Face didn't get enough nodes", i)
		}
	}

	// Split and check each has one
	s.Split(func (tri *Triangle) bool {
		retval := false
		if len(tri.Nodes) > 1 {
			retval = true
		}
		return retval
	})

	if len(s.Triangles) != (8*4) {
		t.Fatal("Incorrect number of triangles after split ", len(s.Triangles))
	}

	for _, tri := range s.Triangles {
		if (len(tri.Nodes) != 1) {
			t.Fatal("Incorrect number of nodes after split")
		}
	}

}


func TestSphereSplitUneven(t *testing.T) {

	coords := [][]float64{
		{22.5, 22.5},
		{22.5, 45.0},
		{22.5, 67.5},
		{67.5, 45.0}}


	s := NewSQT()
	for q := 0; q < 1; q++ {
		for i, ll := range coords {
			s.AddNode(q*len(coords)+i*2,1.0, ll[0], ll[1] + float64(q)*90.0)
			s.AddNode(q*len(coords)+i*2+1,1.0, ll[0] * -1.0, ll[1] + float64(q)*90.0)			
		}
	}
	
	for i,ft := range s.Triangles {
		if (i==0) || (i==4) {
			if ft.NumNodes() !=4  {
				t.Fatal("Face didn't get enough nodes", i)
			}				
		} else {
			if ft.NumNodes() !=0  {
				t.Fatal("Face didn't get enough nodes", i)
			}				
		}
	}

	// Split and check each has one
	s.Split(func (tri *Triangle) bool {
		retval := false
		if len(tri.Nodes) > 1 {
			retval = true
		}
		return retval
	})

	if len(s.Triangles) != (2*4+6) {
		t.Fatal("Incorrect number of triangles after split ", len(s.Triangles))
	}

	for _, tri := range s.Triangles {
		if (len(tri.Nodes) != 1) && (len(tri.Nodes) != 0) {
			t.Fatal("Incorrect number of nodes after split")
		}
	}
	
}

