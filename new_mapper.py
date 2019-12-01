# -*- coding: utf-8 -*-
"""
Created on Wed Aug 14 20:38:49 2019

@author: akshat
"""

from os import listdir

fill = " fill=#2c82c9 "
end ='preserveAspectRatio="xMidYMid meet" class="country" id="%s">'
a = '<a onclick="click_country(this.parentElement.id)" href="javascript:void(0);" onmouseover="over(this.parentElement.id);" onmouseleave="leave(this.parentElement.id)">\n\t\t'
txt = "\n\t</a>\n\t<text name='%s' x='%d' y='%d' fill='#FFFFFF' text-anchor='middle' alignment-baseline='middle'>0</text>"

def htm(path, name):
    with open(path) as file:
        text = file.read()
    p = text[text.find("<path"):text.find("/>")+2]
    p = p.replace(" ", fill, 1)
    svg = text[text.find("<svg"):text.find("<g")]
    tags = svg.split(" ")
    viewBox = "".join(i+" " for i in tags[5:9])
    width = tags[3]
    height = tags[4]
    h = int(height.replace('height="', "").replace('px"', ""))//2
    w = int(width.replace('width="', "").replace('px"', ""))//2

    element = "\n<svg " + width + " " + height + " " + viewBox + (end % name)+ "\n\t" + a + p + (txt % (name, w, h)) + "\n</svg>"

    return element

text = ""

for folder in listdir("new_map"):
    for file in listdir("new_map/"+folder+"/"):
        text += htm("new_map/"+folder+"/"+file, file.replace(".svg", ""))
