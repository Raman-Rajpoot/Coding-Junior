import React, { useEffect, useState } from 'react';
import Cookies from 'js-cookie'; // For managing tokens
import { useNavigate } from 'react-router-dom';
import './Dashboard.css';

const Dashboard = () => {
  const navigate = useNavigate();
  const [userData, setUserData] = useState(null);

  useEffect(() => {
    // Fetch token and user details from cookies or localStorage
    const accessToken = Cookies.get('access_token');
    const refreshToken = Cookies.get('refresh_token');
    const userDetails = {
      userName: Cookies.get('userName'),
      email: Cookies.get('email'),
      fullName: Cookies.get('fullName'),
    };

    if (!accessToken || !userDetails.userName) {
      // Redirect to login if tokens or user details are missing
      navigate('/login');
      return;
    }

    // Set user details in state
    setUserData(userDetails);
  }, [navigate]);

  const handleLogout = () => {
    // Clear cookies on logout
    Cookies.remove('access_token');
    Cookies.remove('refresh_token');
    Cookies.remove('userName');
    Cookies.remove('email');
    Cookies.remove('fullName');

    // Navigate to login page
    navigate('/login');
  };

  if (!userData) {
    return <p>Loading dashboard...</p>;
  }

  return (
    <div className="dashboard-container">
      <h1>Welcome, {userData.fullName || userData.userName}!</h1>
      <p>Email: {userData.email}</p>

      <div className="dashboard-actions">
        <button onClick={() => navigate('/profile')} className="action-btn">Go to Profile</button>
        <button onClick={handleLogout} className="logout-btn">Logout</button>
      </div>
    </div>
  );
};

export default Dashboard;
