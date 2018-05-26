var batchresult;

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
    var optimizer = document.querySelector('input[name="optimizer"]:checked').value;
    var order = document.getElementById('order').value;
    // This will send a message to GO
    astilectron.sendMessage({name: "input", payload: [sx, sy, ex, ey, optimizer, order]}, function(message) {
        result = JSON.parse(message.payload);
        showOrder(0, result);
    });
};

function showBatchOrder() {
    clearMap();
    var orderid = document.getElementById("orderid").value;
    showOrder(orderid, batchresult);
}

function clearMap() {
    for (i=0; i<39; i++) {
        for (j=0; j<23; j++) {
            if (i%2*j%2) {
                var locid = j*39+i;
                var pos = document.getElementsByClassName('item' + locid.toString())[0];
                pos.style.backgroundColor = "rgba(251, 219, 121, 1)";
            } else {
                var locid = j*39+i;
                var pos = document.getElementsByClassName('item' + locid.toString())[0];
                pos.style.backgroundColor = "rgba(255, 255, 255, 1)";
            }
        }
    }
}

function showOrder(orderid, result) {
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