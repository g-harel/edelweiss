interface IDomain {
  id: number;
  name: string;
  data: string;
}

interface IUser {
  id: number;
  email: string;
  password: string;
  domain: string;
}
