'use strict';

const ence = require('ence');
const diip = require('diip');

const types = require('./types');
const defaults = require('./defaults');

const assert = (assertion, message) => {
    if (!assertion) {
        throw new Error(message);
    }
};

const isObject = (obj) => {
    return !!obj || (obj.constructor === Object);
};

const checkType = (type, ...objects) => {
    assert(types[type] && defaults[type], 'unrecognized type');
    [...objects].forEach((obj) => assert(isObject(obj)));
    const obj = Object.assign({}, defaults[type], ...objects);
    assert(!diip(types[type], ence(obj)), 'object doesn\'t match type');
    return obj;
};

module.exports = {assert, isObject, checkType};
