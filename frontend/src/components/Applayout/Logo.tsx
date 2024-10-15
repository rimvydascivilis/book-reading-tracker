import React from 'react';
import logo from '../../assets/logo.png';

const Header: React.FC = () => {
  return (
    <div style={{ display: 'flex', alignItems: 'center', background: 'white', padding: '0' }}>
    <img src={logo} alt="logo" style={{ width: '50px', height: '50px', margin: '0 15px'}} />
    <span style={{ fontSize: '24px', fontWeight: 'bold' }}>
      Book reading tracker
    </span>
  </div>
  )
};

export default Header;
