import React, { useEffect, useState } from 'react';
import { useHistory, Link } from 'react-router-dom';
import { useAuth0 } from '@auth0/auth0-react';

export function Nav(): JSX.Element {
  const { isAuthenticated, user, loginWithRedirect, logout } = useAuth0();
  const history = useHistory();
  const [pathname, setPathname] = useState(() => history.location.pathname);

  useEffect(() => {
    return history.listen(({ pathname }) => setPathname(pathname));
  }, [history]);

  return (
    <nav className="navbar navbar-expand-lg navbar-light bg-light">
      <span className="navbar-brand">GrocerySpend.io</span>
      <div className="collapse navbar-collapse">
        <div className="navbar-nav">
          <Link
            to="/"
            className={`nav-item nav-link${pathname === '/' ? ' active' : ''}`}
          >
            Home
          </Link>
          <Link
            to="/requests"
            className={`nav-item nav-link${
              pathname === '/requests' ? ' active' : ''
            }`}
          >
            Upload Receipt
          </Link>          
          <Link
            to="/receipts"
            className={`nav-item nav-link${
              pathname === '/receipts' ? ' active' : ''
            }`}
          >
            Receipts
          </Link>
          <Link
            to="/analytics"
            className={`nav-item nav-link${
              pathname === '/analytics' ? ' active' : ''
            }`}
          >
            Analytics
          </Link>
        </div>
      </div>

      {isAuthenticated ? (
        <div>
          <span id="hello">Hello, {user?.name}!</span>{' '}
          <button
            className="btn btn-outline-secondary"
            id="logout"
            onClick={() => logout({ returnTo: window.location.origin })}
          >
            logout
          </button>
        </div>
      ) : (
        <button
          className="btn btn-outline-success"
          id="login"
          onClick={() => loginWithRedirect()}
        >
          login
        </button>
      )}
    </nav>
  );
}