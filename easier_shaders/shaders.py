import ctypes
import math
import random


import pyglet
from pyglet.gl import *


window = pyglet.window.Window(900, 900)

vert_shader_source = b'''
    #version 110

    varying vec4 color;

    void main() {
        gl_Position = gl_ProjectionMatrix * gl_ModelViewMatrix * gl_Vertex;
        color = gl_Color;
    }
 '''

frag_shader_source = b'''
    #version 110

    varying vec4 color;

    void main() { 
        gl_FragColor = color;
    }

 '''

def compile_shader(shader_type, shader_source):
    shader_name = glCreateShader(shader_type)
    src_buffer = ctypes.create_string_buffer(shader_source)
    buf_pointer = ctypes.cast(ctypes.pointer(ctypes.pointer(src_buffer)), ctypes.POINTER(ctypes.POINTER(ctypes.c_char)))
    length = ctypes.c_int(len(shader_source) + 1)
    glShaderSource(shader_name, 1, buf_pointer, ctypes.byref(length))
    glCompileShader(shader_name)
    return shader_name

frag_shader = compile_shader(GL_FRAGMENT_SHADER, frag_shader_source)
vert_shader = compile_shader(GL_VERTEX_SHADER, vert_shader_source)

program = glCreateProgram()
glAttachShader(program, vert_shader)
glAttachShader(program, frag_shader)
glLinkProgram(program)
glUseProgram(program)

# Конец части с инициализацией шейдеров
glClearColor(0.7, 0.2, 0.2, 1)

glEnable(GL_CULL_FACE)
glFrontFace(GL_CW)

alpha = (1 + math.sqrt(5)) / 2

pos = [0, 0, -1.5]
rot = 0

vertices = [None] * 12
for i in range(4):
    vertices[i] = [0, 1 - 2 * (i & 2) / 2, alpha - 2 * alpha * (i % 2)]
for i in range(4, 8):
    vertices[i] = [1 - 2 * (i & 2) / 2, alpha - 2 * alpha * (i % 2), 0]
for i in range(8, 12):
    vertices[i] = [alpha - 2 * alpha * (i % 2), 0, 1 - 2 * (i & 2) / 2]

triangles = [
    [0, 2, 8], [0, 8, 4], [0, 4, 6], [0, 6, 9], [0, 9, 2],
    [2, 7, 5], [2, 5, 8], [2, 9, 7], [8, 5, 10], [8, 10, 4],
    [10, 5, 3], [10, 3, 1], [10, 1, 4], [1, 6, 4], [1, 3, 11],
    [1, 11, 6], [6, 11, 9], [11, 3, 7], [11, 7, 9], [3, 5, 7]
]
colors = []
for i in range(20):
    colors.append([random.random(), random.random(), random.random()])

# Каст в сишные вектора
c_vdata = []
for v in vertices:
    c_vdata.append(ctypes.pointer((ctypes.c_float * len(v))(*v)))

def cross(a, b):
    a = [b[1] * a[2] - b[2] * a[1],
         b[2] * a[0] - b[0] * a[2],
         b[0] * a[1] - b[1] * a[0]]

    return a

def drawFigure():
    glBegin(GL_TRIANGLES)
    for i in range(20):
        p1 = vertices[triangles[i][0]]
        p2 = vertices[triangles[i][1]]
        p3 = vertices[triangles[i][2]]
        div = 16
        product = cross([p2[0]-p1[0], p2[1]-p1[1], p2[2]-p1[2]],
                        [p3[0]-p1[0], p3[1]-p1[1], p3[2]-p1[2]])

        glNormal3f(product[0]/div, product[1]/div, product[2]/div)
        glColor3d(colors[i][0], colors[i][1], colors[i][2])
        glVertex3fv(c_vdata[triangles[i][0]][0])
        glVertex3fv(c_vdata[triangles[i][1]][0])
        glVertex3fv(c_vdata[triangles[i][2]][0])
    glEnd()

def draw():

    # Лево низ
    glViewport(0, 0, 450, 450)
    glMatrixMode(GL_MODELVIEW)
    glLoadIdentity()
    glTranslatef(*pos)
    glScalef(0.3, .3, .3)
    glRotatef(rot, 0, 90, 0)
    glMatrixMode(GL_PROJECTION)
    glLoadIdentity()
    glOrtho(-1, 1, -1, 1, 0, 60)
    drawFigure()

    #Право низ
    glViewport(450, 0, 450, 450)
    glMatrixMode(GL_PROJECTION)
    glLoadIdentity()
    gluPerspective(90, 1, 0, 59)
    glMatrixMode(GL_MODELVIEW)
    glLoadIdentity()
    glTranslatef(*pos)
    glScalef(0.4, .4, .4)
    glRotatef(rot, 0, 1, 0)
    drawFigure()


    # Лево верх
    glViewport(0, 450, 450, 450)

    glMatrixMode(GL_MODELVIEW)
    glLoadIdentity()
    glTranslatef(*pos)
    glScalef(0.3, .3, .3)
    glRotatef(rot, 0, 1, 0)

    glMatrixMode(GL_PROJECTION)
    glLoadIdentity()
    glOrtho(-1, 1, -1, 1, 0, 60)
    drawFigure()

    # Право верх
    glViewport(450, 450, 450, 450)

    glMatrixMode(GL_MODELVIEW)
    glLoadIdentity()
    glTranslatef(*pos)
    glScalef(0.3, .3, .3)
    glRotatef(rot, 0, 0, 1)

    glMatrixMode(GL_PROJECTION)
    glLoadIdentity()
    glOrtho(-1, 1, -1, 1, 0, 60)

    drawFigure()


def display():
    glMatrixMode(GL_MODELVIEW)
    glLoadIdentity()
    glTranslatef(*pos)
    draw()
    glFlush()


@window.event
def on_key_press(s, m):
    global rot
    if s == pyglet.window.key.W:
        pos[2] -= .5
    if s == pyglet.window.key.S:
        pos[2] += .5
    if s == pyglet.window.key.G:
        pos[1] -= .5
    if s == pyglet.window.key.H:
        pos[1] += .5
    if s == pyglet.window.key.B:
        pos[0] -= .5
    if s == pyglet.window.key.N:
        pos[0] += .5
    if s == pyglet.window.key.A:
        rot += 15
    if s == pyglet.window.key.D:
        rot -= 15
    if s == pyglet.window.key.X:
        glDisable(GL_CULL_FACE)
        glFrontFace(GL_CW)
        glPolygonMode(GL_FRONT_AND_BACK, GL_LINE)
    if s == pyglet.window.key.Z:
        glEnable(GL_CULL_FACE)
        glFrontFace(GL_CW)
        glPolygonMode(GL_FRONT_AND_BACK, GL_FILL)


@window.event
def on_draw():
    window.clear()
    display()



pyglet.app.run()
