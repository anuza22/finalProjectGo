import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

const Navbar = () => {
  const { isAuthenticated, user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/');
  };

  return (
    <nav className="navbar">
      <div className="container">
        <Link to="/" className="navbar-brand">
          Auth System
        </Link>
        
        <div className="navbar-menu">
          <Link to="/" className="navbar-item">
            Главная
          </Link>
          
          {isAuthenticated ? (
            <>
              <Link to="/profile" className="navbar-item">
                Профиль
              </Link>
              
              <div className="navbar-item user-menu">
                <span className="welcome-text">
                  Привет, {user?.username}
                </span>
                <button onClick={handleLogout} className="btn btn-sm btn-outline">
                  Выйти
                </button>
              </div>
            </>
          ) : (
            <>
              <Link to="/login" className="navbar-item">
                Войти
              </Link>
              <Link to="/register" className="navbar-item btn btn-sm btn-primary">
                Регистрация
              </Link>
            </>
          )}
        </div>
      </div>
    </nav>
  );
};

export default Navbar;