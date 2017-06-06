'use strict';

const port = process.env.PORT || 3000;

const express = require('express');
const bodyParser = require('body-parser');

const config = require('./config');

const app = express();

app.use(bodyParser.urlencoded({
    extended: false,
}));

app.get('/', (req, res) => {
    res.send({status: 'success'});
});

app.listen(port, (err) => {
    if (err) {
        throw new Error(err);
    }
    console.log(`listing on port ${port}`);
});
