import * as Sequelize from 'sequelize';

import Table from './';

interface IDomain {
  id?: number;
  name: string;
  data: string;
}

class DomainsTable extends Table {
  public getName() {
    return 'domains';
  }

  public getAttributes() {
    return {
      id: {
        type: Sequelize.INTEGER,
        primaryKey: true,
        autoIncrement: true,
      },
      name: {
        type: Sequelize.STRING,
        unique: true,
        allowNull: false,
      },
      data: {
        type: Sequelize.JSON,
        allowNull: false,
      },
    };
  }

  public getOptions() {
    return {};
  }

  public getConfig() {
    return {};
  }

  public async get({name}) {
    const domain = await this.model.findOne({
      attributes: ['name', 'data'],
      where: {name},
    }) as IDomain;
    return domain;
  }

  public async delete({name}) {
    const deleted = await this.model.destroy({
      where: {name},
    });
    return !!deleted;
  }
}

export default IDomain;

export {DomainsTable};
