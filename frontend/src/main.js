import React from 'react'
import App from './App'
import { Outlet, useNavigate } from 'react-router-dom';
function Main() {
    const navigate = useNavigate();
  return (
    <div>
        <div className="switch-option">
           <button onClick={() => navigate('/profile')} className="switch-link">Go To Profile</button>
        </div>
        <Outlet />
    </div>
  )
}

export default Main