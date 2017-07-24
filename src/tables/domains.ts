import * as Sequelize from 'sequelize';

interface IDomain {
  id?: number;
  name: string;
  data: string;
}

const name = 'domains';

const attributes = {
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

const options = {};

export default IDomain;

export {name, attributes, options};
