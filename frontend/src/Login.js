import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Cookies from 'js-cookie';
import './Login.css';

function Login() {
  const [email, setEmail] = useState('');
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [fieldError, setFieldError] = useState({
    emailError: '',
    usernameError: '',
    passwordError: '',
  });

  const navigate = useNavigate();

  const fieldValidation = (e, field) => {
    const value = e.target.value;

    if (field === 'email') {
      setEmail(value);
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
      if (!emailRegex.test(value)) {
        setFieldError((prev) => ({
          ...prev,
          emailError: 'Invalid email format',
        }));
      } else {
        setFieldError((prev) => ({
          ...prev,
          emailError: '',
        }));
      }
    } else if (field === 'username') {
      setUsername(value);
      if (value.length < 3) {
        setFieldError((prev) => ({
          ...prev,
          usernameError: 'Username must be at least 3 characters',
        }));
      } else {
        setFieldError((prev) => ({
          ...prev,
          usernameError: '',
        }));
      }
    } else if (field === 'password') {
      setPassword(value);
      const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).{8,}$/;
      if (!passwordRegex.test(value)) {
        setFieldError((prev) => ({
          ...prev,
          passwordError:
            'Password must be at least 8 characters, include an uppercase letter, a lowercase letter, and a number.',
        }));
      } else {
        setFieldError((prev) => ({
          ...prev,
          passwordError: '',
        }));
      }
    }
  };

  const handleLoginSubmit = async (e) => {
    e.preventDefault();
    if (fieldError.emailError || fieldError.usernameError || fieldError.passwordError) {
      setError('Please fix the errors before submitting');
      return;
    }

    try {
      const response = await fetch('http://localhost:8000/api/user/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, userName: username, password }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        setError(errorData.message || 'Login failed');
        return;
      }

      const data = await response.json();
      Cookies.set('access_token', data.data.access_token, { expires: 7 });
      navigate('/profile'); // Redirect to profile page after successful login
    } catch (error) {
      setError('An error occurred while logging in.');
    }
  };

  return (
    <div className="container-login-signUp">
      
      <button onClick={() => navigate('/profile')} className="switch-option-profile">Go To Profile</button>
      
      <div className="login-container">
        <h2 className="form-title">Login</h2>
        <form onSubmit={handleLoginSubmit}>
          <label htmlFor="login-email" className="label-text">Email:</label>
          <input
            type="email"
            id="login-email"
            className="input-field"
            value={email}
            onChange={(e) => fieldValidation(e, 'email')}
          />
          {fieldError.emailError && <p className="error-message">{fieldError.emailError}</p>}

          <label htmlFor="login-username" className="label-text">Username:</label>
          <input
            type="text"
            id="login-username"
            className="input-field"
            value={username}
            onChange={(e) => fieldValidation(e, 'username')}
          />
          {fieldError.usernameError && <p className="error-message">{fieldError.usernameError}</p>}

          <label htmlFor="login-password" className="label-text">Password:</label>
          <input
            type="password"
            id="login-password"
            className="input-field"
            value={password}
            onChange={(e) => fieldValidation(e, 'password')}
          />
          {fieldError.passwordError && <p className="error-message">{fieldError.passwordError}</p>}

          {error && <p className="profile-error">{error}</p>}
          <button type="submit" className="submit-btn">Submit</button>
          <div className="switch-option">
            <div>OR </div>
            <button onClick={() => navigate('/register')} className="switch-link">Sign Up Instead</button>
          </div>
        </form>
      </div>        
    </div>
  );
}

export default Login;
