const express = require('express');
const fs = require('fs');
const path = require('path');
const app = express();

const imageDir = process.argv.slice(2, 3)[0];
console.log(`serve folder: ${imageDir}`);

const port = 5555;
const imageSuffixes = ['jpg', 'jpeg', 'webp', 'png', 'gif', 'bmp'];

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
            } else {
                const filename = e.endsWith('/') ? e.substring(0, e.length - 2) : e;
                for (const suffix of imageSuffixes) {
                    if (filename.endsWith(suffix) || filename.endsWith(suffix.toUpperCase())) {
                        content.push(`<img src='${filePath}' style="width:100%;"/>`);
                        return;
                    }
                }
                // it's not a image
                content.push(`<p><a href='${filePath}'>${filePath}</a></p>`);
                return;
            }
        });

        // put link at top
        const _content = content.sort((a, b) => {
            if (a.indexOf(`<a`) != -1) return -1;
            if (b.indexOf(`<a`) != -1) return 1;
            return a.localeCompare(b);
        })

        res.send(pre + _content.join('') + suf);
    }
});

app.listen(port, function () {
    console.log(`listening on ${port}`);
});

