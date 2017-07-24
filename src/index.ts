import * as bcrypt from 'bcrypt';
import * as express from 'express';
import * as Sequelize from 'sequelize';

import IDomain, * as domains from './tables/domains';
import IUser, * as users from './tables/users';

const sequelize = new Sequelize('postgres://postgres@localhost:5432/edelweiss');

const Domains = sequelize
  .define(domains.name, domains.attributes, domains.options);

const Users = sequelize
  .define(users.name, users.attributes, users.options);

const init = async (force?: boolean) => {
  await Domains.sync({force});
  await Users.sync({force});
};

init(true).then(async () => {
  console.log(await Users.findAll({raw: true})); console.log(await Domains.findAll({raw: true}));
  await Domains.create({
    name: 'test',
    data: '{}',
  });
  await Domains.create({
    name: 'test2',
    data: '{}',
  });
  await Users.create({
    domain: 2,
    email: 'test@test.test',
    password: 'test',
  });
  console.log(await Users.findAll({raw: true})); console.log(await Domains.findAll({raw: true}));
});

setTimeout(() => process.exit(), 4000);
