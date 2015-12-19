const fs = require('fs');
const glob = require('glob');
const spawn = require('child_process').spawn;

var handle = function(err, files) {
    if (err) return;

    files.forEach(function(file) {
        fs.watch(file, function(e) {
            if (e === 'change') {
                spawn('lessc', ['stylesheets/main.less', 'stylesheets/main.css']);
                console.log('regenerating')
            }
        });
    });
};

glob('stylesheets/*.less', null, handle);
