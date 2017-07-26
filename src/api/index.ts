import * as bodyParser from 'body-parser';
import * as express from 'express';

import * as database from '../database';
import IDomain from '../database/tables/domains';
import IUser from '../database/tables/users';

database.init(true)
  .then(({domains, users}) => {
    const app = express();

    app.use(bodyParser.json());

    app.post('/api/user/auth', (req, res) => {
      const {domainId, email, password} = req.body;
      users.auth((err, auth) => {
        if (err) {
          return res.sendStatus(500);
        }
        res.sendStatus(auth ? 200 : 401);
      }, {domainId, email, password});
    });

    app.get('/api/domain/:name', (req, res) => {
      domains.get((err, domain) => {
        if (err) {
          return res.sendStatus(404);
        }
        res.setHeader('Content-Type', 'application/json');
        res.send(JSON.stringify(domain));
        return;
      }, {name: req.params.name});
    });

    app.delete('/api/domain/:name', async (req, res) => {
      domains.del((err, amount) => {
        if (err) {
          return res.sendStatus(500);
        }
        res.sendStatus(200);
      }, {name: req.params.name});
    });

    app.listen(3000, () => {
      console.log('started server');
    });
  });
