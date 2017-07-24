import * as bcrypt from 'bcrypt';
import * as Sequelize from 'sequelize';

interface IUser {
  id?: number;
  domain: number;
  email: string;
  password: string;
}

const name = 'users';

const attributes = {
  id: {
    type: Sequelize.INTEGER,
    primaryKey: true,
    autoIncrement: true,
  },
    domain: {
    type: Sequelize.INTEGER,
    allowNull: false,
    references: {
      model: 'domains',
      key: 'id',
    },
  },
  email: {
    type: Sequelize.STRING,
    allowNull: false,
    unique: 'domain user',
    validate: {
      isEmail: true,
    },
  },
  password: {
    type: Sequelize.STRING(60),
    allowNull: false,
  },
};

const options = {
  hooks: {
    beforeCreate: async (entry: IUser) => {
      entry.password = await bcrypt.hash(entry.password, 10);
      return;
    },
  },
};

export default IUser;

export {name, attributes, options};
