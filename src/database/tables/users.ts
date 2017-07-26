import * as bcrypt from 'bcrypt';
import * as Sequelize from 'sequelize';

interface IUser {
  id?: number;
  domainId: number;
  email: string;
  password: string;
}

const init = async (database) => {
  const name = 'users';
  const attributes = {
    id: {
      type: Sequelize.INTEGER,
      primaryKey: true,
      autoIncrement: true,
    },
    domainId: {
      type: Sequelize.INTEGER,
      allowNull: false,
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
  return database.define(name, attributes, options);
};

const query  = (Model) => {
  const auth = async (callback, {domainId, email, password}) => {
    const user = await Model.findOne({
      attributes: ['password'],
      where: {email, domainId},
    });
    bcrypt.compare(password, user.password, callback);
  };

  return {auth};
};

export default IUser;

export {init, query};
