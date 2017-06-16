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
        priority: 'number',
    }],
    data: null,
};

module.exports = {user, domainConfig};
