import * as Sequelize from 'sequelize';

import * as domains from './tables/domains';
import * as users from './tables/users';

const sequelize = new Sequelize('postgres://postgres@localhost:5432/edelweiss');

const init = async (force: boolean = false) => {
  console.log('0');
  const Domains = await domains.init(sequelize);
  console.log('1');
  const Users = await users.init(sequelize);
  Users.belongsTo(Domains, {onDelete: 'CASCADE', hooks: true});
  await Domains.sync({force});
  await Users.sync({force});

  Domains.create({
    name: 'test',
    data: '{"test": true}',
  });

  Users.create({
    domainId: 1,
    email: 'test@example.com',
    password: 'password123',
  });

  return {
    domains: domains.query(Domains),
    users: users.query(Users),
  };
};

export {init};
