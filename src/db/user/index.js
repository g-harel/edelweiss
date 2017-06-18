'use strict';

const user = () => {
    let store = [];

    const add = (callback, user) => {
        store.push(user);
        callback(null, user);
    };

    const edit = (callback, id, changes) => {
        let editedUser = null;
        store = store.map((user) => {
            if (user.id === id) {
                user = Object.assign({}, user, changes);
                editedUser = user;
            }
            return user;
        });
        callback(editedUser?null:'could not edit user', editedUser||null);
    };

    const remove = (id) => {
        let removedUser = null;
        store.find((user, index) => {
            if (user.id === id) {
                store.splice(index, 1);
                removedUser = user;
                return true;
            }
            return false;
        });
        callback(removedUser?null:'could not remove user', removedUser||null);
    };

    const auth = (callback, domain, email, password) => {
        let authUser = store.find((user) => {
            return user.domain === domain && user.email === email && user.password === password;
        });
        callback(authUser?null:'could not authenticate user', authUser||null);
    };

    return {add, edit, remove, auth};
};

module.exports = user;
