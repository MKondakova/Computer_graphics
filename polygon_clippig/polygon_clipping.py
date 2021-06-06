import time 
import pyglet
from pyglet.gl import *
from pyglet.libs.x11.xlib import False_


class Point:
    def __init__(self, x, y):
        self.x = x
        self.y = y

    def __str__(self):
        return "(" + str(self.x) + ", " + str(self.y) + ")"

    def __eq__(self, other):
        if isinstance(other, Point):
            return self.x == other.x and self.y == other.y
        return False


class Line:
    def __init__(self, p1, p2):
        self.p1 = p1
        self.p2 = p2
        self.x = self.p2.x - self.p1.x
        self.y = self.p2.y - self.p1.y

    def __str__(self):
        return "(" + str(self.p1) + ", " + str(self.p2) + ")"
    
    def scalar_product(self, other):
        "скалярное произведение"
        if isinstance(other, Line):
            x1 = self.p2.x - self.p1.x
            y1 = self.p2.y - self.p1.y
            x2 = other.p2.x - other.p1.x
            y2 = other.p2.y - other.p1.y
            return float(x1*x2 + y1*y2)
        return 0

    def cross_product(self, other):
        "векторное произведение"
        if isinstance(other, Line):
            return float(other.x*self.y-self.x*other.y)
        return 0
    def get_normal(self):
        "возвращает нормаль"
        if self.x == 0 :
            return Line(Point(0, 0), Point(0, 1))
        return Line(Point(0, 0), Point(-1*(self.y/self.x), 1))
    
    def get_midpoint(self):
        return Point((self.p2.x + self.p1.x)/2, (self.p2.y + self.p1.y)/2)


    def reverse(self):
        temp = self.p1
        self.p1 = self.p2
        self.p2 = temp
        


sizeX = 1200
sizeY = 900


points = []
clipper_points = []
segments = []
normals = []
clipper_segments = []
additional_points = []
additional_clippers = []
mouse = Point(0, 0)

DRAWING_LINES = 0
DRAWING_POLYGON = 1
CLIPPING = 2

state = DRAWING_LINES


window = pyglet.window.Window(sizeX, sizeY, resizable=True)


@window.event
def on_mouse_press(x, y, button, modifiers):
    global state
    if button == pyglet.window.mouse.LEFT:
        if state == DRAWING_LINES:
            points.append(Point(x, y))
        if state == DRAWING_POLYGON:
            clipper_points.append(Point(x, y))

    if button == pyglet.window.mouse.RIGHT:
        if state == DRAWING_LINES and len(points)<2:
            return
        if state == DRAWING_POLYGON and len(clipper_points)<3:
            return
        state = (state + 1) % (CLIPPING + 1)
        print(state)


@window.event
def on_mouse_motion(x, y, dx, dy):
    global mouse
    mouse = Point(x, y)

@window.event
def on_key_press(symbol, modifiers):
    if symbol == pyglet.window.key.R:
        reset()


@window.event
def on_resize(width, height):
    global sizeX, sizeY
    sizeX = width
    sizeY = height
    reset()

def reset():
    global points, segments, state, clipper_points, clipper_segments
    points = []
    segments = []
    state = DRAWING_LINES
    clipper_points = []
    clipper_segments = []

def points_to_segments(points):
    segments = []
    if len(points) % 2 == 1:
        points.pop()
    for i in range(0, len(points), 2):
        segments.append(Line(points[i], points[i+1]))
    return segments

def points_to_polygon(points):
    segments = []
    for i in range(len(points)):
        segments.append(Line(points[i], points[(i+1)%len(points)]))  
    return segments


def is_convex(faces):
    if len(faces) < 3:
        return False
    product_positive = faces[len(faces) - 1].cross_product(faces[0]) > 0
    for i in range(len(faces) - 1):
        if (faces[i].cross_product(faces[i+1]) > 0) != product_positive:
            return False
    return True

def find_normals(faces):
    normals = []
    for i in range(len(faces)):
        normal = faces[i].get_normal()
        if normal.scalar_product(Line(faces[i].p1, faces[(i+1)%len(faces)].p2)) < 0:
            normal.reverse()
        normals.append(normal)
    return normals

def cyrus_beck(segment, faces, normals, is_inner):
    t_start = 0
    t_end = 1
    for i in range(len(faces)):
        d = segment.scalar_product(normals[i])
        w = normals[i].scalar_product(Line(faces[i].p1, segment.p1))
        if d == 0:
            if w < 0: #параллельно грани и при этом снаружи
                if is_inner:
                    return []
                return [segment]
            continue
        t = -1*w/d
        if d > 0:
            t_start = max(t_start, t)
        if d < 0:
            t_end = min(t_end, t)
    if t_start <= t_end:
        if is_inner:
            return [Line(Point(segment.p1.x+segment.x*t_start, segment.p1.y+segment.y*t_start),
                Point(segment.p1.x+segment.x*t_end, segment.p1.y+segment.y*t_end))]

        return [Line(Point(segment.p1.x, segment.p1.y),
                Point(segment.p1.x+segment.x*t_start, segment.p1.y+segment.y*t_start)),
                Line(Point(segment.p1.x+segment.x*t_end, segment.p1.y+segment.y*t_end),
                Point(segment.p2.x, segment.p2.y))]
    if is_inner:
        return []
    return [segment]

        

def draw():
    glColor3d(1, 1, 1)
    if state < CLIPPING:
        glBegin(GL_LINES)
        for point in points:
            glVertex2f(point.x, point.y)
        if state == DRAWING_LINES:
            glVertex2f(mouse.x, mouse.y)
        glEnd()

    glBegin(GL_LINE_LOOP)
    for point in clipper_points:
        glVertex2f(point.x, point.y)
    if state == DRAWING_POLYGON:
        glVertex2f(mouse.x, mouse.y)
    glEnd()
            
    if state == CLIPPING:
        #Если нужно показать дополнения
        #glBegin(GL_LINES)
        #for clipper in additional_clippers:
        #    for segment in clipper:
        #        glColor3f(1, 0, 0)
        #        glVertex2f(segment.p1.x, segment.p1.y)
        #        glColor3f(0, 0, 1)
        #        glVertex2f(segment.p2.x, segment.p2.y)
        #glEnd()
        glColor3d(1, 1, 1)
        glBegin(GL_LINES)
        for s in segments:
            glVertex2f(s.p1.x, s.p1.y)
            glVertex2f(s.p2.x, s.p2.y)
        glEnd()


def clipping(clipper_segments, segments, is_inner):
    normals = find_normals(clipper_segments)
    result = []
    for s in segments:
        visible = cyrus_beck(s, clipper_segments, normals, is_inner)
        if len(visible)> 0:
            result.extend(visible)
    return result

def complete_polygon(points):
    global additional_points
    additional_points = []
    temp = []
    if len(points) < 3:
        return 
    while True:
        changes = 0
        previous = points[len(points) - 1]
        i = 0
        while i < len(points):
            current = points[i]
            next = points[(i+1)%len(points)]
            if (Line(previous, current).cross_product(Line(current, next)) > 0):
                changes += 1
                temp.append(previous)
                points.pop(i)
                i -= 1
            elif len(temp) > 0:
                temp.append(previous)
                temp.append(current)
                additional_points.append(temp)
                temp = []
            previous = current
            i += 1
        if len(temp) > 0:
            temp.append(previous)
            temp.append(points[0])
            additional_points.append(temp)
            temp = []
        if changes == 0:
            break
    return points


def display():
    global segments, clipper_segments, clipper_points, additional_clippers
    if state == CLIPPING :
        segments = points_to_segments(points)
        clipper_segments = points_to_polygon(clipper_points)
        convex = is_convex(clipper_segments)
        print(convex)
        if convex:
            segments = clipping(clipper_segments, segments, False)
        else:
            local_points = clipper_points.copy()
            outter_segments = clipping(points_to_polygon(complete_polygon(local_points)), segments, False)
            additional_clippers = []
            for clipper in additional_points:
                additional_clippers.append(points_to_polygon(clipper))
                outter_segments.extend( clipping(points_to_polygon(clipper), segments, True))
            segments = outter_segments

    draw()
    glFlush()



@window.event
def on_draw():
    window.clear()
    display()


pyglet.app.run()
glClearColor(0, 0, 0, 1)
