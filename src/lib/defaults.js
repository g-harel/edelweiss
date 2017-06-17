'use strict';

const user = {
    role: 'admin',
};

const domainConfig = {
    roles: [{
        name: 'admin',
        permission: 0,
    }],
    data: {},
};

const domainConfigLevel = {
    permissions: {
        read: 8,
        write: 0,
    },
    children: null,
};

module.exports = {user, domainConfig, domainConfigLevel};
