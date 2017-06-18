'use strict';

const express = require('express');
const bodyParser = require('body-parser');

const db = {
    user: require('../db/user')(),
    domain: {
        config: require('../db/domain/config')(),
        data: require('../db/domain/data')(),
    },
};

const app = express();

app.use(bodyParser.json());
app.use(bodyParser.urlencoded({extended: true}));

app.post('/api/user/add', (req, res) => {
    db.user.add((err, user) => {
        if (err) {
            res.send({status: err});
            return;
        }
        res.send({status: 'success'});
    }, req.body);
});

app.listen(1234, (err) => {
    console.log('> 1234 <');
});

