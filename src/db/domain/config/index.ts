'use strict';

const config = () => {
    const store = {};

    const add = (callback, domainConfig) => {
        store[domainConfig.name] = domainConfig;
        callback(null, domainConfig);
    };

    const find = (callback, name) => {
        if (store[name]) {
            callback(null, store[name]);
        } else {
            callback('could not find domain', null);
        }
    };

    const remove = (callback, name) => {
        const domainConfig = store[name];
        if (domainConfig) {
            delete store[name];
            callback(null, domainConfig);
        } else {
            callback('could not delete domain', null);
        }
    };

    return {add, find, remove};
};

export = config;
