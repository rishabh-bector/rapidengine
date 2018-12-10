#version 410

layout (triangles) in;
layout (triangle_strip, max_vertices = 9) out;

in VS_OUT {
    vec3 FragPos;
    vec3 TexCoords;
    vec3 Normal;
    mat3 TBN;

    float Visibility;
} gs_in[];

out vec3 FragPos;
out vec3 TexCoords;
out vec3 Normal;
out mat3 TBN;
out float Visibility;

void main() {
    gl_Position = gl_in[0].gl_Position;

    FragPos = gs_in[0].FragPos;
    TexCoords = gs_in[0].TexCoords;
    Normal = gs_in[0].Normal;
    TBN = gs_in[0].TBN;
    Visibility = gs_in[0].Visibility;

    EmitVertex();

    gl_Position = gl_in[1].gl_Position;

    FragPos = gs_in[1].FragPos;
    TexCoords = gs_in[1].TexCoords;
    Normal = gs_in[1].Normal;
    TBN = gs_in[1].TBN;
    Visibility = gs_in[1].Visibility;

    EmitVertex();

    gl_Position = gl_in[2].gl_Position;

    FragPos = gs_in[2].FragPos;
    TexCoords = gs_in[2].TexCoords;
    Normal = gs_in[2].Normal;
    TBN = gs_in[2].TBN;
    Visibility = gs_in[2].Visibility;

    EmitVertex();

    EndPrimitive();
}