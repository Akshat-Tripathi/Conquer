//Setups up the game

const links = {};
var colour = {};

function init() {
    //position labels
    var countries = document.getElementById("map").querySelectorAll(".country");
    var x = {
        "AL": 20, 
        "CA": -65,
        "SF": -15,
        "NY": -15,
        "ME": -28,
        "GR": -10,
        "VE": -15,
        "BR": -10,
        "PE": -15,
        "AG": -20,
        "UK": 0,
        "RO": -5,
        "SC": 2,
        "PR": -12,
        "RU": -10,
        "NR": -10,
        "PL": -10,
        "UR": -20,
        "SI": -20,
        "AF": -5,
        "MI": -20,
        "IN": -17,
        "SE": -27,
        "JP": -7,
        "CH": -20,
        "PA": -12,
        "NA": -15,
        "EG": -10,
        "WA": -20,
        "CN": -18,
        "SA": -13,
        "MA": -15,
        "CO": 10,
        "WO": -10,
        "EO": -10,
        "NG": -10
    }
    var y = {
        "GR": -23,
        "AL": -15,
        "CA": -15,
        "SF": -15,
        "NY": -20,
        "ME": -30,
        "VE": -23,
        "BR": -23,
        "PE": -15,
        "AG": -30,
        "UK": -5,
        "RO": -20,
        "PR": -20,
        "SC": -25,
        "RU": -15,
        "SI": -3, 
        "PL": -20,
        "NR": -15,
        "UR": -22,
        "NA": -20,
        "WA": -30,
        "EG": -15,
        "CN": -15,
        "MA": -15,
        "SA": -15,
        "MI": -20,
        "AF": -25,
        "PA": -20,
        "IN": -15,
        "SE": -32,
        "CH": -20,
        "JP": -10,
        "NG": -15,
        "CO": -10,
        "WO": -10,
        "EO": -10
    }
    for (i=0; i<countries.length; i++) {
        var ID = countries[i].id;
        var original = document.getElementById(ID);
        var lab = document.getElementsByName(ID)[0];
        var w = 0;
        var h = 0;
        var dimensions = original.getBoundingClientRect();
        var width = dimensions["width"] 
        var height = dimensions["height"]
        var left = dimensions["left"]
        var top = dimensions["top"]
        w = x[ID]
        h = y[ID]
        lab.style.position = "absolute"
        lab.style.left = left + width/2 + w -100;
        lab.style.top = top + height/2 + h;
        lab.style.color = "#e8ecf1";
        lab.style.fontFamily = "sans-serif";
        lab.style.fontSize = 20;
        lab.style.fontWeight = 200;
    }
}

init();

function initMap() {
    links["PO"] = [];
    var raw = document.getElementById("countries");
    var text = raw.contentWindow.document.body.childNodes[0].innerText;
    var lines = text.split("\n").slice(0, 36);
    for (var i=0; i<36; i++) {
        var temp = lines[i].split(";");
        var ele = document.getElementById(temp[0]).children[0].children[0];
        ele.style.fill = temp[2];
        if (temp[2] == "#d2d7d3") {
            document.getElementsByName(temp[0])[0].style.color = "#2e3131"
        }
        colour[temp[0]] = temp[2];
        var lab = document.getElementsByName(temp[0])[0];
        lab.innerHTML = temp[1];
        if (temp[2].toString() == colId) {
            links["PO"] = links["PO"].concat(temp[0])
        }
        links[temp[0]] = temp.slice(3, temp.length);
    }
}

var cook = document.cookie.split("; ")
var player = cook[0].replace("id=", "") 
var colId = cook[1].replace("col=", "")
console.log(player, colId);

function colours() {
    document.getElementById("moves").style="background-color: "+colId+";"
    document.getElementById("execute").style="background-color: "+colId+";"
}

colours();

let socket = new WebSocket(window.location.href.replace("http://", "ws://") + "ws/"+player);
//let socket = new WebSocket("ws://146.169.207.63:8080/ws/"+player);
socket.onopen = () => {
    console.log("Successfully Connected");
};

socket.onclose = event => {
    console.log("Socket Closed Connection: ", event);
    socket.send("Client Closed!")
    disable()
};

socket.onerror = error => {
    console.log("Socket Error: ", error);
};

socket.onmessage = (msg) => {
    var action = JSON.parse(msg.data);
    if (action.Player == "admin") {
        alert("invalid move")
        return
    }
    console.log(msg);
    if (action.Dest == "PO") {
        var src = document.getElementsByName(action.Src)[0];
        src.innerHTML = "Available Troops: "+(action.Numdest+1).toString();
    } else {
        console.log(action)
        execute_action(action);
    };
}

function send(data) {
    console.log(JSON.stringify(data));
    console.log(data);
    socket.send(JSON.stringify(data));
    source = "";
    dest = "";
    updateUI()
}

function execute_action(action) {
    console.log(action.Player, colId)
    var src = document.getElementsByName(action.Src)[0]
    var dest = document.getElementsByName(action.Dest)[0]
    if (action.Src == "PO") {
        if (action.Player == colId) {
            src.innerHTML = "Available Troops: " + (parseInt(src.innerHTML.slice(18), 10)+action.Numsrc).toString();    
        }
    } else {
        src.innerHTML = (parseInt(src.innerHTML, 10)+action.Numsrc).toString();
    }
    dest.innerHTML = (parseInt(dest.innerHTML, 10)+action.Numdest).toString();
    console.log(action.MoveType)
    if (action.MoveType == 3) {
        var col = action.Player
        console.log(col, "hi", action.Dest)
        render(col, action.Dest)
        colour[action.Dest] = col
        links["PO"] = links["PO"].concat(action.Dest)
    }
    console.log("done")
};

function render(col, country) {
    var ele = document.getElementById(country).children[0].children[0];
    ele.style.fill = col;
    var lab = document.getElementsByName(country)[0].style.color = "#e8ecf1";
}

var source = "";
var dest = "";
var troopers = 0;
var actionTaken = ""
var locations = {
    "AL": "Alaska",
    "CA": "Canada",
    "GR": "Greenland",
    "SF": "the Western USA",
    "NY": "the Eastern USA",
    "ME": "Mexico",
    "VE": "Venezuela",
    "PE": "Peru",
    "BR": "Brazil",
    "AG": "Argentina",
    "UK": "the United Kingdom",
    "RO": "Rome 2.0",
    "PR": "Prussia",
    "RU": "Russia",
    "SC": "Scandinavia",
    "SI": "Siberia",
    "NR": "Northern Russia",
    "PL": "Place",
    "UR": "Kamchatka",
    "MI": "the Middle East",
    "AF": "Afganistan",
    "PA": "Pablo",
    "IN": "India",
    "SE": "South East Asia",
    "CH": "China",
    "JP": "Japan",
    "NA": "Northern Africa",
    "WA": "Western Africa",
    "SA": "South Africa",
    "CN": "Central Africa",
    "MA": "Madagascar",
    "EG": "Egypt",
    "EO": "Eastern Australia",
    "WO": "Western Australia",
    "NG": "New Guinea",
    "CO": "Central Australia",
    "PO": "Base"
}

function updateUI() {
    if (actionTaken != "Attacking") {
        document.getElementById("moves").innerHTML = actionTaken + " " + troopers.toString() + " troops from " + locations[source] + " to " + locations[dest]
    } else {
        document.getElementById("moves").innerHTML = "Attacking " + locations[dest] + " from " + locations[source]
    }

}

class action {
    constructor(Player, Src, Dest, Numsrc, Numdest, MoveType) {
        this.Player = Player;
        this.Src = Src;
        this.Dest = Dest;
        this.Numsrc = Numsrc;
        this.Numdest = Numdest;
        this.MoveType = MoveType;
    }
}
var act = new action()

function process() {
    if (act.Src != null) {
        send(act)
        troopers=0;
        updateUI()
    }
}

function attack() {
    if (source != "") {
        if (source == "PO") {
            alert("Cannot attack from base")
        }
        if(dest != "") {
            a = new action(player, source, dest, 0, 0, 0)
            actionTaken = "Attacking"
            updateUI()
            act = a
        } else {
            alert("No destination selected")
        }
    } else {
        alert("No source selected")
    }
}

function donate() {
    if (troopers != 0 & source != "" & dest != "") {
        a = new action(player, source, dest, -troopers, troopers, 1)
        actionTaken="Donating"
        updateUI()
        act = a
    } else {
        if (troopers == 0) {
            alert("Moving 0 troops")
        }
        if (source == "") {
            alert("No source selected")
        }
        if (dest == "") {
            alert("No destination selected")
        }
    }
}

function move() {
    if (troopers != 0 & source != "" & dest != "") {
        a = new action(player, source, dest, -troopers, troopers, 2)
        actionTaken="Moving"
        updateUI()
        act = a
    } else {
        if (troopers == 0) {
            alert("Moving 0 troops")
        }
        if (source == "") {
            alert("No source selected")
        }
        if (dest == "") {
            alert("No destination selected")
        }
    }
}

function add(troops){
    troopers += troops
    actionTaken = "Moving"
    updateUI()
}

var highlighted = [];
var disabled = [];

function reset_countries() {
    while (disabled.length != 0) {
        var id = disabled.pop();
        var ele = document.getElementById(id).children[0].children[0];
        ele.style = "pointer-events: fill;";
        ele.style.fill = colour[id]
    }
    while (highlighted.length != 0) {
        var id = highlighted.pop();
        var ele = document.getElementById(id);
        ele.style = "";
    }
}

function click_country(code) {
    if (source=="PO") {
        dest=code
        updateUI()
        return
    }
    if (highlighted.length != 0) {
        if (highlighted[0] == code) {
            source = "";
        } else {
            dest = code;
        }
        updateUI()
        reset_countries();
        return
    }
    source = code
    dest = ""
    updateUI()
    var temp = links[code];
    for (var i=0; i<temp.length; i++) {
        var ele = document.getElementById(temp[i]);
        ele.style="filter: brightness(130%);";
    }
    highlighted = [code].concat(temp);
    var countries = Object.keys(colour);
    for (var i=0; i<36; i++) {
        if (!highlighted.includes(countries[i])) {
            var ele = document.getElementById(countries[i]).children[0].children[0];
            ele.style = "pointer-events:none;";
            ele.style.fill = colour[countries[i]]
            disabled = disabled.concat(countries[i]);
        }
    }
}

function disable() {
    var countries = Object.keys(colour);
    for (var i=0; i<36; i++) {
        var country = document.getElementById(countries[i])
        country.children[0].children[0].style = "pointer-events: none;"
    }
}

function over(code) {
    document.getElementById(code).style="filter: brightness(85%);";
}

function leave(code) {
    if (highlighted.slice(1).includes(code)) {
        document.getElementById(code).style="filter: brightness(130%)"
    } else {
        document.getElementById(code).style = "filter: brightness(100%)"
    }
}

initMap();