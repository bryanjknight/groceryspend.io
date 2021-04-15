import { useApi } from './use-api';
import React from 'react';
import { Loading } from './Loading';
import { Error } from './Error';

const PORT = 8080;

export function Users() {
  const { loading, error, data: users = [] } = useApi(
    `http://localhost:${PORT}/users/`,
    {
      audience: "https://bknight.dev.groceryspend.io",
      scope: 'read:users',
      mode: "cors",
      credentials: "include",
    }
  );

  if (loading) {
    return <Loading />;
  }

  if (error) {
    return <Error message={error.message} />;
  }

  return (
    <table className="table">
      <thead>
        <tr>
          <th scope="col">Name</th>
          <th scope="col">Email</th>
        </tr>
      </thead>
      <tbody>
        {users?.map(
          ({ name, email }: { name: string; email: string }, i: number) => (
            <tr key={i}>
              <td>{name}</td>
              <td>{email}</td>
            </tr>
          )
        )}
      </tbody>
    </table>
  );
}
