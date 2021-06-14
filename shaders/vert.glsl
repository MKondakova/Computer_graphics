#version 110

varying vec4 color;
varying vec2 texCoord;
varying vec3 normal;
varying vec3 fragPos;

void main() {
    texCoord = gl_MultiTexCoord0.xy;
    gl_Position = gl_ProjectionMatrix * gl_ModelViewMatrix * gl_Vertex;
    color = gl_Color;
    vec4 temp = gl_ModelViewMatrix * vec4(gl_Normal, 0.0);
    normal = temp.xyz * -1.0;
    vec4 position = gl_ModelViewMatrix * gl_Vertex;
    fragPos = position.xyz;
}