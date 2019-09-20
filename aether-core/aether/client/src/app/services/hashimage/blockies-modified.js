(function() {
    // The random number is a js implementation of the Xorshift PRNG
    var randseed = new Array(4); // Xorshift: [x, y, z, w] 32 bit values

    function seedrand(seed) {
        for (var i = 0; i < randseed.length; i++) {
            randseed[i] = 0;
        }
        for (var i = 0; i < seed.length; i++) {
            randseed[i % 4] = ((randseed[i % 4] << 5) - randseed[i % 4]) + seed.charCodeAt(i);
        }
    }

    function rand() {
        // based on Java's String.hashCode(), expanded to 4 32bit values
        var t = randseed[0] ^ (randseed[0] << 11);

        randseed[0] = randseed[1];
        randseed[1] = randseed[2];
        randseed[2] = randseed[3];
        randseed[3] = (randseed[3] ^ (randseed[3] >> 19) ^ t ^ (t >> 8));

        return (randseed[3] >>> 0) / ((1 << 31) >>> 0);
    }

    // function createColorOld() {
    //     //saturation is the whole color spectrum
    //     var h = Math.floor(rand() * 360);
    //     //saturation goes from 40 to 100, it avoids greyish colors
    //     var s = ((rand() * 60) + 40) + '%';
    //     //lightness can be anything from 0 to 100, but probabilities are a bell curve around 50%
    //     var l = ((rand() + rand() + rand() + rand()) * 25) + '%';

    //     var color = 'hsl(' + h + ',' + s + ',' + l + ')';
    //     return color;
    // }
    let used = []

    function isUsed(colour) {
        for (var i = 0; i < used.length; i++) {
            return (function(i) {
                if (colour === used[i]) {
                    return true
                }
            })(i)
        }
        return false
    }

    function createColor() {
        var hueNum = Math.floor(rand() * 6)
        let red = '#EC5f67'
        let orange = '#F99157'
        let yellow = '#FAC863'
        // let green = '#99C794'
        let turquoise = '#5FB3B3'
        let cerulean = '#6699CC'
        let purple = '#C594C5'
        // let earth = '#AB7967'
        let colours = [red, orange, yellow, turquoise, cerulean, purple]
        while (isUsed(colours[hueNum])) {
            hueNum = Math.floor(rand() * 6) // Mind that rand() is deterministic given a specific seed. So this is actually not generating a random number, and the secondary colour is actually not randomly chosen - it will be exactly the same at every run.
        }
        used.push(colours[hueNum])
        return colours[hueNum]
    }

    function createImageData(size) {
        var width = size; // Only support square icons for now
        var height = size;

        var dataWidth = Math.ceil(width / 2);
        var mirrorWidth = width - dataWidth;

        var data = [];
        for (var y = 0; y < height; y++) {
            var row = [];
            for (var x = 0; x < dataWidth; x++) {
                // this makes foreground and background color to have a 43% (1/2.3) probability
                // spot color has 13% chance
                row[x] = Math.floor(rand() * 2.3);
            }
            var r = row.slice(0, mirrorWidth);
            r.reverse();
            row = row.concat(r);

            for (var i = 0; i < row.length; i++) {
                data.push(row[i]);
            }
        }

        return data;
    }

    function buildOpts(opts) {
        var newOpts = {};

        newOpts.seed = opts.seed || Math.floor((Math.random() * Math.pow(10, 16))).toString(16);

        seedrand(newOpts.seed);

        newOpts.size = opts.size || 8;
        newOpts.scale = opts.scale || 4;
        // Reset used colours first
        used = []
        newOpts.color = opts.color || createColor();
        newOpts.bgcolor = opts.bgcolor || createColor();
        newOpts.spotcolor = opts.spotcolor || createColor();

        return newOpts;
    }

    function renderIcon(opts, canvas) {
        opts = buildOpts(opts || {});
        var imageData = createImageData(opts.size);
        var width = Math.sqrt(imageData.length);

        canvas.width = canvas.height = opts.size * opts.scale;

        var cc = canvas.getContext('2d');
        cc.fillStyle = opts.bgcolor;
        cc.fillRect(0, 0, canvas.width, canvas.height);
        cc.fillStyle = opts.color;

        for (var i = 0; i < imageData.length; i++) {

            // if data is 0, leave the background
            if (imageData[i]) {
                var row = Math.floor(i / width);
                var col = i % width;

                // if data is 2, choose spot color, if 1 choose foreground
                cc.fillStyle = (imageData[i] == 1) ? opts.color : opts.spotcolor;

                cc.fillRect(col * opts.scale, row * opts.scale, opts.scale, opts.scale);
            }
        }
        return canvas;
    }

    function createIcon(opts) {
        var canvas = document.createElement('canvas');

        renderIcon(opts, canvas);

        return canvas;
    }

    var api = {
        create: createIcon,
        render: renderIcon
    };

    if (typeof module !== "undefined") {
        module.exports = api;
    }
    if (typeof window !== "undefined") {
        window.blockies = api;
    }

})();