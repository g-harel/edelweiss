interface IDomain {
  uid: number;
  name: string;
  data: string;
}

interface IUser {
  uid: number;
  email: string;
  password: string;
  domain: string;
}
