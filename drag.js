var element = document.getElementById("IN");
var left = 0
var t = 0
var h = 15
var w = 15
var reso = 1

function click_country(id) {
    element = document.getElementById(id);
    left = parseInt(element.style.left, 10);
    t = parseInt(element.style.top, 10);
    if (isNaN(left)) {
        t = 0
        left = 0
    }
}

document.onkeydown = checkKey;

function checkKey(e) {

    e = e || window.event;

    if (e.keyCode == '38') {
        t -= reso 
    }
    else if (e.keyCode == '40') {
        t += reso
    }
    else if (e.keyCode == '37') {
       left -= reso
    }
    else if (e.keyCode == '39') {
       left += reso
    }
    else if (e.keyCode == '87') {
        h += reso
    }
    else if (e.keyCode == '83') {
        h -= reso
    }
    else if (e.keyCode == '65') {
        w -= reso
    }
    else if (e.keyCode == '68') {
        w += reso
    }
    element.style.left = left.toString()
    element.style.top =  t.toString()
    element.style.height = h.toString()+"%"
    element.style.width = w.toString()+"%"
}

function over(id) {

}

function leave(id) {

}