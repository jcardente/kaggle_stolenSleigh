// --------------------------------------------------------------------------
// icotest
//
// Simple prototype of a spherical quad tree for the Kaggle Santa Stolen
// Sleigh challenge. 
//
// J. Cardente
// 2016

import gifAnimation.*;

GifMaker gifExport;
int frames = 0;
int totalFrames = 120;
int duration    = 5; // seconds
int fRate       = Math.round(totalFrames / duration);
int fDelayms    = duration/totalFrames*1000;

PVector origin;
int height = 200;
int width  = 600;

SQT sqt1,sqt2,sqt3;

void setup() {
  Node n;
  float radius = min(height, width/3)/2-30;

  size(600, 200, P3D);
  ortho();
  frameRate(fRate);

  gifExport = new GifMaker(this, "export.gif");
  gifExport.setRepeat(0);
  //gifExport.setTransparent(0,0,0);

  origin = new PVector(width/2, height/2, 0);

  sqt1 = new SQT(radius);
  sqt2 = new SQT(radius);
  sqt3 = new SQT(radius);

  n = new Node(35.0, 45.0);
  sqt1.addNode(n);
  sqt2.addNode(n);
  sqt3.addNode(n);
  
  // NB - all nodes must be added before splitting.
  sqt2.split(1);
  sqt3.split(2);
}


void draw() {
  int i;
  float rot;
  
  rot = 2*PI*frames/totalFrames;
 
  rotateX(-1*radians(5));
  background(#eff0f6);
  noFill();
  stroke(#8999BD);  
  strokeWeight(1);

  pushMatrix();
  translate(width/6, origin.y, origin.z);
  rotateY(rot);
  sqt1.display();
  popMatrix();
  
  pushMatrix();
  translate(3*width/6, origin.y, origin.z);
  rotateY(rot);
  sqt2.display();
  popMatrix();

  pushMatrix();
  translate(5*width/6, origin.y, origin.z);
  rotateY(rot);
  sqt3.display();
  popMatrix();

  export();
}


void export() {
 if (frames < totalFrames) {

   //gifExport.setDelay(fDelayms);
   gifExport.addFrame(g.get());
   frames++;
 } else {
   gifExport.finish();
   frames++;
   println("Export finished");
   exit();
 }
  
}


// UTILITIES -----------------------------------------------------------------------------------

PVector sphere2cart(float lat, float lon, float radius) {
  // Assumes latitudes between +-90 degrees
  float x, y, z;
  float _lat, _lon;

  _lat = PI/2+lat;
  _lon = lon;
  x = radius * sin(_lat) * cos(_lon);
  y = -1*radius * sin(_lat) * sin(_lon);
  z = radius * cos(_lat);
  return new PVector(x, z, y);
}



float roundsmall(float f) {
  return (abs(f) < 0.0001) ? 0.0 : f;
}


// SQT CLASS --------------------------------------------------------------------------------------

class SQT {
  Face[] faces = {};
  float radius;

  SQT(float _radius) {
    radius = _radius;
    int i;
    
    // NB - Looping over nothern and southern hemisphere's
    //      separately to make the "which face" calculation
    //      easier.
    for (i=0; i < 4; i++) {
      faces = (Face[]) append(faces, new Face(0, i * 90.0, true));
    }
    
    for (i=0; i < 4; i++) {
      faces = (Face[]) append(faces, new Face(0, i * 90.0, false));
    }

  }

  int whichFace(float lat, float lon) {
    int  faceIdx = 0;
    if (lat < 0.0) {
      faceIdx += 1;
    }
    faceIdx += 2*floor((lon+360.0) % 360.0)/90; 

    return faceIdx;
  }

  void addNode(Node n) {
    int faceIdx;
    faceIdx = whichFace(n.lat, n.lon);
    faces[faceIdx].addNode(n);
  }

  void split(int count) {
    int i, j;
    for (j=0; j < count; j++) {
      for (i=0; i < faces.length; i++) {
        faces[i].split();
      }
    }
  }

  void display() {
    int i;
    for (i=0; i < faces.length; i++) {
      faces[i].display(radius);
    }
  }
}

// LATLONG CLASS --------------------------------------------------------------------------------------

class Geo {
  float lat, lon;

  Geo(float _lat, float _lon) {
    lat = _lat;
    lon = _lon;
  }
}


// NODE ---------------------------------------------------------------------------

class Node {
  float lat, lon;
  Geo g;
  PVector p;

  Node(float _lat, float _lon) {
    g = new Geo(0, 0);
    g.lat = _lat;
    g.lon = _lon;
    p = new PVector(0, 0);
  }

  void printNode() {
    print("N: " + g.lat + ":" + g.lon + "  " + p.x + ":" + p.y + "\n");
  }
}


// FACE ------------------------------------------------------------------------------------

class Face {
  float lat, lon;
  boolean upward;
  Triangle[] tris = {};
  float len = 1;

  Face(float _lat, float _lon, boolean _upward) {
    lat = _lat;
    lon = _lon;
    upward = _upward;

    tris = (Triangle[]) append(tris, new Triangle(0.0, 0.0, 1.0, true));
  }

  void addNode(Node n) {
    // Project onto triangle
    n.p = triangleProject(abs(n.g.lat-lat), abs(n.g.lon-lon));
    tris[0].addNode(n);
  }


  // lat lon relative to face's start    
  PVector triangleProject(float _lat, float _lon) {
    PVector p = new PVector(0, 0);
    float rlat, rlon;

    rlat = radians(_lat);
    rlon = radians(_lon);
    
    p.x = len/PI * (rlat + 2*rlon*(1-2/PI*rlat));
    p.y = len * sqrt(3) /PI * rlat;

    return p;
  }


  void split() {
    int i, j;
    Triangle[] newtris = {};
    Triangle[] splits;

    for (i = 0; i < tris.length; i++) {
      splits = tris[i].split();
      for (j=0; j < splits.length; j++) {
        newtris = (Triangle[]) append(newtris, splits[j]);
      }
    }

    tris = newtris;
  }


  void printTris() {
    int i;
    for (i = 0; i < tris.length; i++) {
      print(i + ": " + tris[i].x + " " + tris[i].y + " " + tris[i].len + " " + tris[i].upward + "\n");
    }
  }

  void printVerts() {
    int i;
    for (i = 0; i < tris.length; i++) {
      print(i + ": ");
      tris[i].printVerts();
    }
  }

  void printSphereVerts() {
    int i;
    for (i = 0; i < tris.length; i++) {
      print(i + ": ");
      tris[i].printSphereVerts();
    }
  }


  void displayFlat(float size) {
    int i;
    for (i = 0; i < tris.length; i++) {
      tris[i].displayFlat(size);
    }
  }

  void display(float radius) {
    int i;
    PVector n1, n2, n3; 
    Geo[] sverts = {null, null, null};    
    float dircoef = (upward) ? 1.0 : -1.0;
    float rlat, rlon;

    rlat = radians(lat);
    rlon = radians(lon);
    for (i = 0; i < tris.length; i++) {

      // Get spherical coordinates        
      sverts = tris[i].sphereVerts();

      // Convert to cartesian

      n1 = sphere2cart(rlat + sverts[0].lat * dircoef, rlon + sverts[0].lon, radius);
      n2 = sphere2cart(rlat + sverts[1].lat * dircoef, rlon + sverts[1].lon, radius);
      n3 = sphere2cart(rlat + sverts[2].lat * dircoef, rlon + sverts[2].lon, radius);

      if (tris[i].nodes.length > 0) {
        fill(#294080, 192);
      } else {
        noFill();
      }

      // Draw shape
      beginShape();
      vertex(n1.x, n1.y, n1.z);
      vertex(n3.x, n3.y, n3.z);   
      vertex(n2.x, n2.y, n2.z);
      endShape(CLOSE);
    }
  }
}


// TRIANGLE ----------------------------------------------------------------------------------

class Triangle {
  float   x, y, len;
  boolean upward; 
  Node[] nodes = {};

  Triangle(float _x, float _y, float _len, boolean _upward) {
    x   = roundsmall(_x);
    y   = roundsmall(_y);
    len = _len;
    upward = _upward;
  }


  void addNode(Node n) {
    nodes = (Node[]) append(nodes, n);
  }

  PVector[] vertices() {
    PVector[] verts = {new PVector(0, 0), 
      new PVector(0, 0), 
      new PVector(0, 0)};
    float dircoef = upward ?  1.0 : -1.0;

    verts[0].x = x;
    verts[0].y = y;

    verts[1].x = x + len;
    verts[1].y = y;

    verts[2].x = x + len/2;
    verts[2].y = roundsmall(y + sin(PI/3) * len * dircoef);

    return verts;
  }

  void printVerts() {
    PVector[] verts = {null, null, null};

    verts = vertices();
    print("Verts: " + 
      verts[0].x + ":" + verts[0].y +  " " + 
      verts[1].x + ":" + verts[1].y +  " " + 
      verts[2].x + ":" + verts[2].y + "\n");
  }

  Geo sphereProject(PVector p, float facelen) {
    Geo g = new Geo(0, 0); 
    float tmp;   

    g.lat = (PI * p.y) / (sqrt(3)*facelen);

    tmp = sqrt(3)*p.x-p.y;
    g.lon = (abs(tmp) < 0.005) ? 0 : PI/(2*facelen) * (tmp / (sqrt(3) - 2*p.y/facelen));

    return g;
  }


  Geo[] sphereVerts() {
    PVector[] verts = vertices();

    // FIXME: Send face size to sphereProject
    Geo[] sverts = {sphereProject(verts[0], 1), 
      sphereProject(verts[1], 1), 
      sphereProject(verts[2], 1)};    
    return sverts;
  }


  void printSphereVerts() {
    PVector[] verts = vertices();
    Geo[] sverts = {sphereProject(verts[0], 1), 
      sphereProject(verts[1], 1), 
      sphereProject(verts[2], 1)};    


    print("SVerts: " + 
      degrees(sverts[0].lat) + ":" + degrees(sverts[0].lon) + " " +
      degrees(sverts[1].lat) + ":" + degrees(sverts[1].lon) + " " +
      degrees(sverts[2].lat) + ":" + degrees(sverts[2].lon) + "\n");
  }


  void displayFlat(float size) {
    PVector[] verts = {null, null, null};

    verts = vertices();

    beginShape();
    vertex(verts[0].x*size, verts[0].y*size, 0);
    vertex(verts[1].x*size, verts[1].y*size, 0);
    vertex(verts[2].x*size, verts[2].y*size, 0);
    endShape(CLOSE);
  }

  Triangle[] split() {
    Triangle[] subs = {null, null, null, null};
    float dircoef = (upward) ? 1.0 : -1.0;
    int i;
    float dx, dy;

    subs[0] = new Triangle(x, y, len/2, upward);
    subs[1] = new Triangle(x+len/2, y, len/2, upward);
    subs[2] = new Triangle(x+len/4, y+sqrt(3)/2*len/2*dircoef, len/2, upward);
    subs[3] = new Triangle(x+len/4, y+sqrt(3)/2*len/2*dircoef, len/2, !upward);

    for (i = 0; i < nodes.length; i++) {

      dx = abs(nodes[i].p.x - x);
      dy = abs(nodes[i].p.y - y);

      if (dy > sqrt(3)*len/4) {
        subs[2].addNode(nodes[i]);
      } else if (dy < (sqrt(3)*(len-2*dx)/2)) {
        subs[0].addNode(nodes[i]);
      } else if (dy < (sqrt(3)*(2*dx-len)/2)) {
        subs[1].addNode(nodes[i]);
      } else {
        subs[3].addNode(nodes[i]);
      }
    }


    return subs;
  }
}