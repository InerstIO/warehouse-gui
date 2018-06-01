var batchresult;
var reOrderMapping;
var reverseMap;

let index = {
    init: function() {
        // Init
        asticode.loader.init();
        asticode.modaler.init();
        asticode.notifier.init();

        // This will wait for the astilectron namespace to be ready
        document.addEventListener('astilectron-ready', function() {
            // This will send a message to GO
            astilectron.sendMessage({name: "batchresult"}, function(message) {
                batchresult = JSON.parse(message.payload);
            });
            astilectron.sendMessage({name: "reOrderMapping"}, function(message) {
                reOrderMapping = JSON.parse(message.payload);
            });
            astilectron.sendMessage({name: "reverseMap"}, function(message) {
                reverseMap = JSON.parse(message.payload);
            });
        })
    },

    listen: function() {
        astilectron.onMessage(function(message) {
            switch (message.name) {
                case "about":
                    index.about(message.payload);
                    return {payload: "payload"};
                    break;
                case "check.out.menu":
                    asticode.notifier.info(message.payload);
                    break;
            }
        });
    }
};

var sendInput = function() {
    clearMap();
    var sx = document.getElementById('sx').value;
    var sy = document.getElementById('sy').value;
    var ex = document.getElementById('ex').value;
    var ey = document.getElementById('ey').value;
    var weight = document.getElementById('weight').value;
    var order = document.getElementById('order').value;
    // This will send a message to GO
    astilectron.sendMessage({name: "input", payload: [sx, sy, ex, ey, weight, order]}, function(message) {
        result = JSON.parse(message.payload);
        c1 = showRouteItem(result, true);
        c1.inc();
    });
};

function showRouteItem(result, manual) {
    var i = -1;
    return {
        inc: function() {
            i += 1;
            clearMap();
            showRoute(i, result.Paths);
            showItem(i, result.Products, manual);
            showRouteString(i, result.Orders, result.Start, result.End);
        },
        dec: function() {
            i -= 1;
            clearMap();
            showRoute(i, result.Paths);
            showItem(i, result.Products, manual);
            showRouteString(i, result.Orders, result.Start, result.End);
        }
    }
}

function showRouteString(i, orders, start, end) {
    // This will send a message to GO
    astilectron.sendMessage({name: "route", payload: {Order:orders[i], Start:start, End:end}}, function(message) {
        routeString = JSON.parse(message.payload);
        document.getElementById('instruction').innerHTML = routeString;
    });
    astilectron.sendMessage({name: "effort", payload: {Order:orders[i], Start:start, End:end}}, function(message) {
        effort = JSON.parse(message.payload);
        document.getElementById('effort').innerHTML = "Effort: " + effort;
    });
}

function showBatchOrder() {
    clearMap();
    var orderid = document.getElementById("orderid").value;
    c1 = showRouteItem(batchresult[reOrderMapping[parseInt(orderid)]], false);
    c1.inc();
}

function clearMap() {
    for (i=0; i<39; i++) {
        for (j=0; j<23; j++) {
            if (i%2*j%2) {
                var locid = j*39+i;
                var pos = document.getElementsByClassName('item' + locid.toString())[0];
                pos.style.backgroundColor = "rgba(251, 219, 121, 1)";
                pos.innerHTML = ""
            } else {
                var locid = j*39+i;
                var pos = document.getElementsByClassName('item' + locid.toString())[0];
                pos.style.backgroundColor = "rgba(255, 255, 255, 1)";
            }
        }
    }
}

function showRoute(orderid, result) {
    for (var i in result[orderid].slice(1)) {
        var cur = result[orderid][i];
        var next = result[orderid][parseInt(i)+1];
        if (next.X == cur.X) {
            var range = Array.from(new Array(Math.abs(cur.Y-next.Y)+1), (x,i) => i + Math.min(cur.Y,next.Y));
            for (var j of range) {
                var locid = j*39+cur.X;
                var pos = document.getElementsByClassName('item' + locid.toString())[0];
                pos.style.backgroundColor = "lightblue";
            }
        } else if (next.Y == cur.Y) {
            var range = Array.from(new Array(Math.abs(cur.X-next.X)+1), (x,i) => i + Math.min(cur.X,next.X));
            for (var j of range) {
                var locid = cur.Y*39+j;
                var pos = document.getElementsByClassName('item' + locid.toString())[0];
                pos.style.backgroundColor = "lightblue";
            }
        }
    }
}

function showItem(orderid, items, manual) {
    for (var i of items[orderid]) {
        var locid = i.Pos.Y*39 + i.Pos.X;
        var pos = document.getElementsByClassName('item' + locid.toString())[0];
        pos.style.backgroundColor = "Fuchsia";
        if (manual) {
            pos.innerHTML = orderid;
        } else {
            pos.innerHTML = i.OrderID;
        }
    }
}