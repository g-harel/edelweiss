'use strict';

const user = {
    id: 'number',
    email: 'string',
    password: 'string',
    role: 'string',
    domain: 'string',
};

const domainConfig = {
    name: 'string',
    roles: [{
        name: 'string',
        permission: 'number',
    }],
    data: null, // domainConfigLevel
};

const domainConfigLevel = {
    permissions: {
        read: 'number',
        write: 'number',
    },
    children: null, //  object OR array
};

module.exports = {user, domainConfig, domainConfigLevel};
