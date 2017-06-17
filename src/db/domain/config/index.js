'use strict';

const {assert, checkType} = require('../../../lib');

const config = () => {
    const store = {};

    const add = (callback, domainConfig) => {
        domainConfig = checkType('domainConfig', domainConfig);
        assert(!store[domainConfig.name], 'domain name must be unique');
        store[domain.name] = domainConfig;
        callback(null, true);
        console.log('> added domain\n', store);
    };

    const find = (callback, name) => {
        if (store[name]) {
            callback(null, store[name]);
        } else {
            callback(null, null);
        }
        console.log('> found\n', name);
    };

    const remove = (callback, name) => {
        delete store[name];
        callback(null, true);
        console.log('> removed domain\n', store);
    };

    return {add, find, remove};
};

module.exports = config;
