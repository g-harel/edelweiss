import * as Sequelize from 'sequelize';

interface IDomain {
  id?: number;
  name: string;
  data: string;
}

const init = async (database) => {
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
  return database.define(name, attributes, options);
};

const query = (Model) => {
  const get = (callback, {name}) => {
    const domain = Model.findOne({
      attributes: ['name', 'data'],
      where: {name},
    })
      .then((res) => callback(null, res))
      .catch((reason) => callback(reason, null));
  };

  const del = (callback, {name}) => {
    Model.destroy({
      where: {name},
    })
      .then((res) => callback(null, res))
      .catch((reason) => callback(reason, null));
  };

  return {get, del};
};

export default IDomain;

export {init, query};
