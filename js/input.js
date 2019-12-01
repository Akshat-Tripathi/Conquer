
var source = "";
var dest = "";
var troopers = 0;

class action {
    constructor() {
        this.Player = ""
        this.Src = ""
        this.Dest = ""
        this.Numsrc = 0
        this.Numdest = 0
        this.Movetype = 0
    }
}
 

function attack() {
    if (source != "") {
        if(dest != "") {
            a = new action(player, source, dest, 0, 0, 0)
            send(a)
        }
    }
}

function donate() {
    if (troopers != 0) {
        a = new action(player, source, dest, -troopers, troopers, 1)
        send(a)
    }
}

function move() {
    if (troopers != 0) {
        a = new action(player, source, dest, -troopers, troopers, 2)
        send(a)
    }
}

function add(troops){
    troopers += troops
}

function click(country) {
    if (source == "") {
        source = country
    }else {
        if (dest == "") {
            dest = country
        } else{
            source = ""
            dest = ""
        }
    }
}
