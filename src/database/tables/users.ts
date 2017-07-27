import * as bcrypt from 'bcrypt';
import * as Sequelize from 'sequelize';

import Table from './';

interface IUser {
  id?: number;
  domainId: number;
  email: string;
  password: string;
}

class UsersTable extends Table {
  public getName() {
    return 'users';
  }

  public getAttributes() {
    return {
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
  }

  public getOptions() {
    return {
      hooks: {
        beforeCreate: async (entry: IUser) => {
          entry.password = await bcrypt.hash(entry.password, 10);
          return;
        },
      },
    };
  }

  public getConfig() {
    return {
      belongsTo: {
        name: 'domains',
        options: {
          onDelete: 'CASCADE',
          hooks: true,
        },
      },
    };
  }

  public async authenticate({domainId, email, password}) {
    const user = await this.model.findOne({
      attributes: ['password'],
      where: {email, domainId},
    }) as IUser;
    if (user == null) {
      return await false;
    }
    return await bcrypt.compare(password, user.password);
  }
}

export default IUser;

export {UsersTable};
