const express = require('express');
const app = express();

app.get('/', (req, res) => {
    res.send('Hello World!')
});

app.get('/archimedes', (req, res) => {
    res.json({ coming: "soon" })
})

app.listen(3000, () => {
    console.log('example app listening on port 3000')
});

console.log('init')