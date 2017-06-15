'use strict';

const ence = require('ence');
const diip = require('diip');

const userType = {
    id: 'number',
    email: 'string',
    password: 'string',
    role: 'string',
    domain: 'string',
};

const defaultUser = {
    role: 'admin',
};

const assert = (assertion, message) => {
    if (!assertion) {
        throw new Error(message);
    }
};

const isObject = (obj) => {
    return !!obj || (obj.constructor === Object);
};

const check = (user) => {
    return !diip(userType, ence(user));
};

const users = () => {
    let store = [];

    const add = (user) => {
        assert(isObject(user), 'user is not an object');
        user = Object.assign({id: Math.random()*1000000000000000000}, defaultUser, user);
        assert(check(user), 'user format is incorrect');
        assert(!store.find((storedUser) => storedUser.id === user.id), 'user id must be unique');
        assert(!store.find((storedUser) => (storedUser.email === user.email && storedUser.domain === user.domain)), 'user email must be unique for the same domain');
        store.push(user);
        console.log('> added user\n', store);
    };

    const edit = (id, changes) => {
        store = store
            .map((user) => {
                if (user.id === id) {
                    assert(isObject(changes), 'changes is not an object');
                    let newUser = Object.assign({}, user, changes);
                    assert(check(newUser), 'user format is incorrect');
                }
                return user;
            });
        console.log('> edited user\n', store);
    };

    const remove = (id) => {
        store.find((user, index) => {
            if (user.id === id) {
                store.splice(index, 1);
                return true;
            }
            return false;
        });
        console.log('> removed user\n', store);
    };

    const auth = (domain, email, password) => {
        return !!store.find((user) => {
            return user.domain === domain && user.email === email && user.password === password;
        });
    };

    return {add, edit, remove, auth};
};

module.exports = users;
