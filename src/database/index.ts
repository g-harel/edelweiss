import * as Sequelize from 'sequelize';

import * as domains from './tables/domains';
import * as users from './tables/users';

const sequelize = new Sequelize('postgres://postgres@localhost:5432/edelweiss');

const Domains = sequelize
  .define(domains.name, domains.attributes, domains.options);

const Users = sequelize
  .define(users.name, users.attributes, users.options);

const init = async (force: boolean = false) => {
  Users.belongsTo(Domains, {onDelete: 'CASCADE', hooks: true});
  await Domains.sync({force});
  await Users.sync({force});
  return {Domains, Users};
};

export {init};
