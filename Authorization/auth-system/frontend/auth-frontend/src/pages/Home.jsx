import React from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

const Home = () => {
  const { isAuthenticated } = useAuth();

  return (
    <div className="home-page">
      <h1>Добро пожаловать в нашу систему</h1>
      <p>Это пример системы авторизации с использованием Go и React.</p>
      
      {!isAuthenticated ? (
        <div className="auth-buttons">
          <Link to="/login" className="btn btn-primary">Войти</Link>
          <Link to="/register" className="btn btn-secondary">Зарегистрироваться</Link>
        </div>
      ) : (
        <div className="welcome-back">
          <p>Вы уже авторизованы.</p>
          <Link to="/profile" className="btn btn-primary">Перейти в профиль</Link>
        </div>
      )}
    </div>
  );
};

export default Home;