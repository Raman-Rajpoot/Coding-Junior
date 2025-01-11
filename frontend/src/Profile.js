import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom'; // For navigation
import Cookies from 'js-cookie';
import './Profile.css';

function Profile() {
  const [profileData, setProfileData] = useState(null);
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const clearAllCookies = () => {
    // Get all cookies
    const cookies = document.cookie.split(';');
  
    // Iterate over cookies and delete each one
    cookies.forEach((cookie) => {
      const name = cookie.split('=')[0].trim();
      document.cookie = `${name}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;`;
    });
  };

  useEffect(() => {
    const token = Cookies.get('access_token'); // Retrieve the token from cookies

    if (!token) {
      setError('You are not logged in.');
      clearAllCookies();
      setTimeout(() => navigate('/login'), 2000); // Redirect to login after 2 seconds
      return;
    }

    const fetchProfile = async () => {
        const access_token = Cookies.get('access_token');

       console.log( access_token)
      try {
        const response = await fetch('http://localhost:8000/api/user/profile', {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${access_token}`, // Attach the token in the Authorization header
          },
          credentials: 'include', // Send cookies with the request
        });

        if (!response.ok) {
          const errorData = await response.json();
          console.error('Error response:', errorData); // Log error details for debugging

          // Handle unauthorized access
          if (response.status === 401) {
            setError('Session expired. Redirecting to login.');
            clearAllCookies();
            setTimeout(() => navigate('/login'), 2000);
          } else {
            setError(errorData.message || 'Failed to fetch profile. Refresh it to try again.');
            setTimeout(() => navigate('/login'), 2000);
          }
          return;
        }

        // Parse and set profile data on success
        const data = await response.json();
        setProfileData(data.data);
      } catch (error) {
        console.error('Fetch error:', error);
        setError('An error occurred while fetching profile data. Refresh it to try again.');
        setTimeout(() => navigate('/login'), 2000);
      }
    };

    fetchProfile();
  }, [navigate]);

  // Render error message if present
  if (error) {
    return <div className="profile-error">{error}</div>;
  }

  // Show loading message while fetching profile data
  if (!profileData) {
    return <div className="loading-message">Loading profile...</div>;
  }

  // Render profile data
  return (
    <div className="profile-container">
      <h2>Profile</h2>
      <p><strong>Full Name:</strong> {profileData.fullName || 'N/A'}</p>
      <p><strong>Email:</strong> {profileData.email || 'N/A'}</p>
      <p><strong>Username:</strong> {profileData.userName || 'N/A'}</p>

      <button className='btn' onClick={() => navigate('/login')}>Go To Login</button>
    </div>
  );
}

export default Profile;
