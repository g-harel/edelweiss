'use strict';

const {assert, checkType} = require('../../lib');

const user = () => {
    let store = [];

    const add = (callback, user) => {
        user = checkType('user', {id: Math.random()*1000000000000000000}, user);
        assert(!store.find((storedUser) => storedUser.id === user.id), 'user id must be unique');
        assert(!store.find((storedUser) => (storedUser.email === user.email && storedUser.domain === user.domain)), 'user email must be unique for the same domain');
        store.push(user);
        callback(null, true);
        console.log('> added user\n', store);
    };

    const edit = (callback, id, changes) => {
        store = store
            .map((user) => {
                if (user.id === id) {
                    user = checkType('user', user, changes);
                }
                return user;
            });
        callback(null, true);
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
        callback(null, true);
        console.log('> removed user\n', store);
    };

    const auth = (callback, domain, email, password) => {
        callback(null, store.find((user) => {
            return user.domain === domain && user.email === email && user.password === password;
        }));
    };

    return {add, edit, remove, auth};
};

module.exports = user;
