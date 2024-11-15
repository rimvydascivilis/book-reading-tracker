import React from 'react';
import Logo from './Logo';
import {Layout, Button} from 'antd';
import {LogoutOutlined} from '@ant-design/icons';
import {useAuth} from '../../context/AuthContext';

const AntHeader = Layout.Header;

const Header: React.FC = () => {
  const {isAuthenticated, logout} = useAuth();

  return (
    <AntHeader
      style={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        background: 'white',
        padding: '0 24px 0 0',
      }}>
      <Logo />
      {isAuthenticated && (
        <Button onClick={logout} type="primary" style={{marginLeft: 'auto'}}>
          <LogoutOutlined />
          Logout
        </Button>
      )}
    </AntHeader>
  );
};

export default Header;
