let index = {
    init: function() {
        // Init
        asticode.loader.init();
        asticode.modaler.init();
        asticode.notifier.init();

        // This will wait for the astilectron namespace to be ready
        document.addEventListener('astilectron-ready', function() {
            // This will listen to messages sent by GO
            astilectron.onMessage(function(message) {
                // Process message
                if (message.name === "route") {
                    return {payload: message.payload + " world"};
                }
            });
        })

        // This will wait for the astilectron namespace to be ready
        document.addEventListener('astilectron-ready', function() {
            // This will send a message to GO
            astilectron.sendMessage({name: "event.name", payload: "hello"}, function(message) {
                console.log("received " + message.payload)
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
    var sx = document.getElementById('sx').value;
    var sy = document.getElementById('sy').value;
    var ex = document.getElementById('ex').value;
    var ey = document.getElementById('ey').value;
    var optimizer = document.querySelector('input[name="optimizer"]:checked').value;
    var iter = document.getElementById('iterations').value;
    var method = document.querySelector('input[name="input"]:checked').value;
    var order = document.getElementById('order').value;
    var ordersfile = document.getElementById('ordersfile').value;
    var outputfile = document.getElementById('outputfile').value;
    // This will send a message to GO
    astilectron.sendMessage({name: "input", payload: [sx, sy, ex, ey, optimizer, iter, method, order, ordersfile, outputfile]}, function(message) {
        console.log("received " + message.payload)
    });
}