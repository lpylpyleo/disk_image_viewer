const express = require('express');
const fs = require('fs');
const path = require('path');
const app = express();

const imageDir = process.argv.slice(2, 3)[0];
console.log(`serve folder: ${imageDir}`);

const port = 5555;
const imageSuffixes = ['jpg', 'jpeg', 'webp', 'png', 'gif', 'bmp'];
const videoSuffixes = ['mp4', 'webm', 'ts', 'wmv', 'mkv', 'avi', 'mts'];

function isKindOf(name, suffixes) {
    const _filename = name.endsWith('/') ? name.substring(0, e.length - 2) : name;
    for (const suffix of suffixes) {
        if (_filename.endsWith(suffix) || _filename.endsWith(suffix.toUpperCase())) {
            return true;
        }
    }
    return false;
}

const isImage = (name) => isKindOf(name, imageSuffixes);

const isVideo = (name) => isKindOf(name, videoSuffixes);

app.get('/v/:vid', function (req, res) {
    const reqPath = decodeURIComponent(req.path);
    const absolutePath = path.join(imageDir, reqPath);
    res.send(`<video src='/${req.params.vid}' controls autoplay='true'></video>`);
});

app.get('*', function (req, res) {
    const reqPath = decodeURIComponent(req.path);
    const absolutePath = path.join(imageDir, reqPath);

    if (fs.statSync(absolutePath).isFile()) {
        res.sendFile(absolutePath);
    } else {
        const pre = `<html><body>`;
        const suf = `</body></html>`;
        const files = fs.readdirSync(absolutePath);
        let content = [];
        files.forEach(e => {
            // ignore dotfile
            if (e.startsWith('.')) return;

            const stat = fs.statSync(path.join(absolutePath, e));
            const filePath = (reqPath === '/' ? '' : reqPath) + '/' + e;
            if (stat.isDirectory()) {
                content.push(`<p><a href='${filePath}'>${filePath}</a></p>`);
                return;
            } else if (isImage(e)) {
                content.push(`<img src='${filePath}' style="width:100%;"/>`);
            } else {
                content.push(`<p><a href='${isVideo(filePath) && 'v'}${filePath}'>${filePath}</a></p>`);
            }
        });

        // put link at top
        const _content = content.sort((a, b) => {
            if (a.indexOf(`<a`) != -1) return -1;
            if (b.indexOf(`<a`) != -1) return 1;
            return a.localeCompare(b);
        });

        res.send(pre + _content.join('') + suf);
    }
});

app.listen(port, function () {
    console.log(`listening on ${port}`);
});

