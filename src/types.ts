interface Role {
    name: string,
    permission: number,
}

interface User {
    id: number,
    email: string,
    password: string,
    role: Role,
    domain: string,
};

interface DomainData {
    permissions: {
        read: number,
        write: number,
    }
    data: ['primitive', number | string] |
        ['object', DomainData] |
        ['array', Array<DomainData>],
}

interface Domain {
    name: string,
    roles: Array<Role>,
    data: null | DomainData,
};