import React from 'react';
import { useAuth } from '../contexts/AuthContext';

const Profile = () => {
  const { user, logout } = useAuth();

  const handleLogout = () => {
    logout();
  };

  return (
    <div className="profile-container">
      <h2>Профиль пользователя</h2>
      
      <div className="profile-info">
        <div className="profile-avatar">
          {/* Генерируем аватарку из первой буквы имени пользователя */}
          <div className="avatar-placeholder">
            {user?.username?.charAt(0).toUpperCase()}
          </div>
        </div>
        
        <div className="profile-details">
          <h3>{user?.username}</h3>
          <p>ID: {user?.id}</p>
        </div>
      </div>
      
      <div className="profile-actions">
        <button onClick={handleLogout} className="btn btn-danger">
          Выйти из системы
        </button>
      </div>
    </div>
  );
};

export default Profile;