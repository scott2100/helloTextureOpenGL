#version 330 core
out vec4 FragColor;

in vec3 ourColor;
in vec2 TexCoord;

// texture sampler
uniform sampler2D texture1;
uniform sampler2D texture2;
uniform float mixValue;

void main()
{
	//FragColor = mix( texture(texture1, TexCoord), texture(texture2, TexCoord), texture(texture2, TexCoord).a * 1);
	//FragColor = mix( texture(texture1, TexCoord), texture(texture2, vec2(TexCoord.x, TexCoord.y)), 0.2, texture(texture2, vec2(TexCoord.x, TexCoord.y)).a * 1);
	//FragColor = mix( texture(texture1, TexCoord), texture(texture2, vec2(TexCoord.x, TexCoord.y)), (texture(texture2, vec2(TexCoord.x, TexCoord.y)).a * 1) * 0.5);
	FragColor = mix( texture(texture1, TexCoord), texture(texture2, vec2(TexCoord.x, TexCoord.y)), (texture(texture2, vec2(TexCoord.x, TexCoord.y)).a * 1) * mixValue);
}
    \x00