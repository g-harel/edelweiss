import * as express from 'express';

import * as databse from '../database';
import IDomain from '../database/tables/domains';
import IUser from '../database/tables/users';

databse.init(true)
  .then(({Domains, Users}) => {
    const app = express();

    Domains.create({
      name: 'test',
      data: '{"test": true}',
    });

    Users.create({
      domainId: 1,
      email: 'test@example.com',
      password: 'password123',
    });

    app.get('/api/domain/:name', async (req, res) => {
      const domain = await Domains
        .find({
          attributes: ['name', 'data'],
          where: {
            name: req.params.name,
          },
        });
      if (domain) {
        res.setHeader('Content-Type', 'application/json');
        res.send(JSON.stringify(domain));
        return;
      }
      res.sendStatus(404);
    });

    app.delete('/api/domain/:name', async (req, res) => {
      const success = await Domains
        .destroy({
          where: {
            name: req.params.name,
          },
        });
      res.sendStatus(200);
    });

    app.listen(3000, () => {
      console.log('started server');
    });
  });
