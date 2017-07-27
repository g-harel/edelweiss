import * as bodyParser from 'body-parser';
import * as express from 'express';

import * as database from '../database';
import IDomain from '../database/tables/domains';
import IUser from '../database/tables/users';

database.init(true)
  .then(({Domains, Users}) => {
    const app = express();

    app.use(bodyParser.json());

    app.post('/api/user/auth', async (req, res) => {
      const {domainId, email, password} = req.body;
      const success = await Users.authenticate({domainId, email, password});
      res.sendStatus(success ? 200 : 401);
    });

    app.get('/api/domain/:name', async (req, res) => {
      const domain = await Domains.get({name: req.params.name});
      res.setHeader('Content-Type', 'application/json');
      res.send(JSON.stringify(domain));
    });

    app.delete('/api/domain/:name', async (req, res) => {
      const success = await Domains.delete({name: req.params.name});
      res.sendStatus(success ? 200 : 401);
    });

    app.use((err, req, res, next) => {
      console.log(err, err.stack);
      res.sendStatus(500);
    });

    app.listen(3000, () => {
      console.log('started server');
    });
  });
