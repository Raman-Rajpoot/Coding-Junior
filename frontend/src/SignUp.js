import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Cookies from 'js-cookie';
import './SignUp.css';

function SignUp() {
  const [email, setEmail] = useState('');
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [fullName, setFullName] = useState('');
  const [error, setError] = useState('');
  const [fieldError, setFieldError] = useState({
    emailError: '',
    usernameError: '',
    passwordError: '',
    fullNameError: '',
  });

  const navigate = useNavigate();

  const fieldValidation = (value, field) => {
    switch (field) {
      case 'email':
        setEmail(value);
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        setFieldError((prev) => ({
          ...prev,
          emailError: emailRegex.test(value) ? '' : 'Invalid email format',
        }));
        break;

      case 'username':
        setUsername(value);
        setFieldError((prev) => ({
          ...prev,
          usernameError: value.length >= 3 ? '' : 'Username must be at least 3 characters',
        }));
        break;

      case 'password':
        setPassword(value);
        const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).{8,}$/;
        setFieldError((prev) => ({
          ...prev,
          passwordError: passwordRegex.test(value)
            ? ''
            : 'Password must be at least 8 characters, include an uppercase letter, a lowercase letter, and a number.',
        }));
        break;

      case 'fullName':
        setFullName(value);
        setFieldError((prev) => ({
          ...prev,
          fullNameError: value.length >= 3 ? '' : 'Full name must be at least 3 characters',
        }));
        break;

      default:
        break;
    }
  };

  const handleSignupSubmit = async (e) => {
    e.preventDefault();

    // Prevent submission if there are validation errors
    if (Object.values(fieldError).some((err) => err)) {
      setError('Please fix the errors before submitting');
      return;
    }

    try {
      const response = await fetch('http://localhost:8000/api/user/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, username, password, fullName }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        setError(errorData.message || 'Sign up failed');
        return;
      }

      // Automatically log in after successful signup
      setTimeout(async () => {
        try {
          const loginResponse = await fetch('http://localhost:8000/api/user/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, userName: username, password }),
          });

          if (!loginResponse.ok) {
            const loginError = await loginResponse.json();
            setError(loginError.message || 'Login failed after successful sign-up. Please log in manually.');
            setTimeout(() => navigate('/login'), 2000); // Redirect to login page after 2 seconds
            return;
          }

          const loginData = await loginResponse.json();
          Cookies.set('access_token', loginData.data.access_token, { expires: 7 });
          navigate('/profile'); // Redirect to profile page
        } catch (error) {
          setError('An error occurred during auto-login. Please log in manually.');
        }
      }, 2000);
    } catch (error) {
      setError('An error occurred while signing up.');
    }
  };

  return (
    <div className="container-signup">
      <button onClick={() => navigate('/profile')} className="switch-option-profile">
        Go To Profile
      </button>
      <div className="signup-container">
        <h2 className="form-title">Sign Up</h2>
        <form onSubmit={handleSignupSubmit}>
          {/* Email Field */}
          <label htmlFor="signup-email" className="label-text">Email:</label>
          <input
            type="email"
            id="signup-email"
            className="input-field"
            value={email}
            onChange={(e) => fieldValidation(e.target.value, 'email')}
          />
          {fieldError.emailError && <p className="error-message">{fieldError.emailError}</p>}

          {/* Username Field */}
          <label htmlFor="signup-username" className="label-text">Username:</label>
          <input
            type="text"
            id="signup-username"
            className="input-field"
            value={username}
            onChange={(e) => fieldValidation(e.target.value, 'username')}
          />
          {fieldError.usernameError && <p className="error-message">{fieldError.usernameError}</p>}

          {/* Full Name Field */}
          <label htmlFor="signup-fullName" className="label-text">Full Name:</label>
          <input
            type="text"
            id="signup-fullName"
            className="input-field"
            value={fullName}
            onChange={(e) => fieldValidation(e.target.value, 'fullName')}
          />
          {fieldError.fullNameError && <p className="error-message">{fieldError.fullNameError}</p>}

          {/* Password Field */}
          <label htmlFor="signup-password" className="label-text">Password:</label>
          <input
            type="password"
            id="signup-password"
            className="input-field"
            value={password}
            onChange={(e) => fieldValidation(e.target.value, 'password')}
          />
          {fieldError.passwordError && <p className="error-message">{fieldError.passwordError}</p>}

          {/* Error Message */}
          {error && <p className="error-message">{error}</p>}

          {/* Submit Button */}
          <button type="submit" className="submit-btn">Submit</button>

          {/* Switch to Login */}
          <div className="switch-option">
            <div>OR</div>
            <button onClick={() => navigate('/login')} className="switch-link">
              Login Instead
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default SignUp;
