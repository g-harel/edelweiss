import * as Sequelize from 'sequelize';

import {IConfig} from './tables';
import IDomain, {DomainsTable} from './tables/domains';
import IUser, {UsersTable} from './tables/users';

const sequelize = new Sequelize('postgres://postgres@localhost:5432/edelweiss');

const init = async (force: boolean = false) => {
  const Domains = new DomainsTable();
  const Users = new UsersTable();

  // init all tables
  const tablesArray = [Domains, Users];
  await Promise.all(
    tablesArray.map((table) => {
      return table.init(sequelize);
    }),
  );

  // map tablesArray into a hashtable
  const tables = {};
  tablesArray.forEach((table) => {
    tables[table.getName()] = table;
  });

  // handle belongsTo relationships
  tablesArray.forEach((table) => {
    const config = table.getConfig() as IConfig;
    if ('belongsTo' in config) {
      const owner = tables[config.belongsTo.name];
      if (owner !== undefined) {
        table.model.belongsTo(owner.model, config.belongsTo.options);
      }
    }
  });

  // sync all tables
  for (const table of tablesArray) {
    await table.model.sync({force});
  }

  try {
    Domains.model.create({
      name: 'test',
      data: '{"test": true}',
    });

    Users.model.create({
      domainId: 1,
      email: 'test@example.com',
      password: 'password123',
    });
  } catch (e) {
    console.error(e);
  }

  return {Domains, Users};
};

export {init};
