import * as express from 'express';
import { Client } from 'pg';

const app = express();

const client = new Client({
  database: 'edelweiss',
  host: 'localhost',
  port: 5432,
  user: 'postgres',
});

const tableExists = async (name: string): Promise<boolean> => {
  const res = await client.query(`
    SELECT 1
    FROM pg_tables
    WHERE tablename = '${name}'
  `);
  return Boolean(res.rows && res.rows.length);
};

const makeSureTableExists = async (name: string): Promise<boolean> => {
  const exists = await tableExists(name);
  if (exists) {
    return true;
  }
  await client.query(`
    CREATE TABLE users (
      id     serial PRIMARY KEY,
      domain text   NOT NULL,
      email  text   NOT NULL,
    )
  `);
  return false;
};

client.connect(async (err) => {
  if (err) {
    throw new Error(err.message);
  }
  await makeSureTableExists('users');
  const res = await client.query('SELECT * from users');
  console.log(res.rows);
  process.exit(0);
});
