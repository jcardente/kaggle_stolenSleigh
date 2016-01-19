

PVector origin;
int height = 500;
int width  = 500;
int radius = 100;
final float circumference = radius * 2 * PI;

Face[] faces = {};

void setup() {
  int i;
  Node n;
  
  size(500, 500, P3D);
  textFont(createFont("dialog", 16));
  
  origin = new PVector(width/2,height/2,0);

  //faces = (Face[]) append(faces, new Face(0.0,0.0,true));
  for (i=0; i < 4; i++) {
   faces = (Face[]) append(faces, new Face(0, i * PI/2, true));
   faces = (Face[]) append(faces, new Face(0, i * PI/2, false));
  }

  n = new Node(new Geo(radians(-85),radians(45)),
               new PVector(0,0));
  
  faces[1].addNode(n);

  
  for (i=0; i < faces.length; i++) {
  faces[i].split(); 
  faces[i].split(); 
      faces[i].split(); 
  }

  
  ////faces[0].split();
  //faces[0].printTris();
  //faces[0].printVerts();
  //faces[0].printSphereVerts();
}


void draw() {
  int i;
    
  background(16);
  noFill();
  stroke(255,255,255);  
  lights();
    
  pushMatrix();
  translate(2*width/3, origin.y, origin.z);
  
  rotateY(frameCount*PI/300);
  for (i=0; i < faces.length; i++) {
    //faces[i].displayFlat(radius);
    faces[i].display();
  }  
  popMatrix();
   
}



// UTILITIES -----------------------------------------------------------------------------------

PVector sphere2cart(float lat, float lon) {
   // Assumes latitudes between +-90 degrees
   float x, y, z;
   float _lat, _lon;
   
   _lat = PI/2+lat; //PI-radians(90-lat);
   _lon = lon; //radians(lon);
   x = radius * sin(_lat) * cos(_lon);
   y = -1*radius * sin(_lat) * sin(_lon);
   z = radius * cos(_lat);
   return new PVector(x,z,y);   
}



float roundsmall(float f) {
   return (abs(f) < 0.0001) ? 0.0 : f; 
}


// LATLONG CLASS --------------------------------------------------------------------------------------

class Geo {
  float lat, lon;
   
  Geo(float _lat, float _lon) {
    lat = _lat;
    lon = _lon;
  }
}


Geo sphereProject(PVector p, float facelen) {
   Geo g = new Geo(0,0); 
   float tmp;   
   
   g.lat = (PI * p.y) / (sqrt(3)*facelen);
   
   tmp = sqrt(3)*p.x-p.y;
   g.lon = (abs(tmp) < 0.005) ? 0 : PI/(2*facelen) * (tmp / (sqrt(3) - 2*p.y/facelen));
  
  return g;
}


// NODE ---------------------------------------------------------------------------

class Node {
   Geo g;
   PVector p;
   
   Node(Geo _g, PVector _p) {
     g = _g;
     p = _p;
   }
   
   void display() {
     PVector p = sphere2cart(g.lat, g.lon);
     
     pushMatrix();
     fill(255,0,0);
     sphere(10);
     popMatrix();
     
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
    
    tris = (Triangle[]) append(tris, new Triangle(0.0,0.0,1.0, true));        
  }

  void addNode(Node n) {
    // Project onto triangle
    n.p = triangleProject(abs(n.g.lat-lat), abs(n.g.lon-lon));
    tris[0].addNode(n); 
  }


  // lat lon relative to face's start    
  PVector triangleProject(float _lat, float _lon) {
     PVector p = new PVector(0,0);
   
     p.x = len/PI * (_lat + 2*_lon*(1-2/PI*_lat));
     p.y = len * sqrt(3) /PI * _lat;
    
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
    
  void display() {
    int i;
    PVector n1, n2, n3; 
    Geo[] sverts = {null, null, null};    
    float dircoef = (upward) ? 1.0 : -1.0;
    
    for (i = 0; i < tris.length; i++) {
      
      // Get spherical coordinates        
      sverts = tris[i].sphereVerts();

     // Convert to cartesian
     n1 = sphere2cart(lat + sverts[0].lat * dircoef, lon + sverts[0].lon);
     n2 = sphere2cart(lat + sverts[1].lat * dircoef, lon + sverts[1].lon);
     n3 = sphere2cart(lat + sverts[2].lat * dircoef, lon + sverts[2].lon);

     if (tris[i].nodes.length > 0) {
      fill(255,0,0); 
     } else {
      noFill(); 
     }
   
     // Draw shape
     beginShape();
     vertex(n1.x,n1.y,n1.z);
     vertex(n3.x,n3.y,n3.z);   
     vertex(n2.x,n2.y,n2.z);
     endShape(CLOSE);
    }
  }
}


// TRIANGLE ----------------------------------------------------------------------------------

class Triangle {
 float   x,y,len;
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
  PVector[] verts = {new PVector(0,0),
                    new PVector(0,0),
                    new PVector(0,0)};
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
   vertex(verts[0].x*size, verts[0].y*size,0);
   vertex(verts[1].x*size, verts[1].y*size,0);
   vertex(verts[2].x*size, verts[2].y*size,0);
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
     
     //print("S: " + dx + ":" + dy + "\n");
     //print("Sub2: dy < " + sqrt(3)*len/4 + "\n");
     //print("Sub0: dy < " + (sqrt(3)*(len-2*dx)/2)+ "\n");
     //print("Sub1: dy < " + (sqrt(3)*(2*dx-len)/2) + "\n");
     
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