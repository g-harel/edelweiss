import * as Sequelize from 'sequelize';

interface IConfig {
    belongsTo?: {name: string; options: object};
}

abstract class Table {
  public model: Sequelize.Model<{}, {}>;

  public abstract getName(): string;
  public abstract getAttributes(): object;
  public abstract getOptions(): object;
  public abstract getConfig(): IConfig;

  public async init(database) {
    this.model = await database.define(this.getName(), this.getAttributes(), this.getOptions());
    return this.model;
  }
}

export default Table;

export {IConfig};
