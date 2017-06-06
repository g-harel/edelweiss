'use strict';

const merge = require('lodash.merge');

const secret = require('./secret');

const config = {};

module.exports = merge(config, secret);
