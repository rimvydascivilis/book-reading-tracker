import React from 'react';
import Logo from './Logo';
import {Layout, Dropdown, Avatar} from 'antd';
import {LogoutOutlined, SettingOutlined, UserOutlined} from '@ant-design/icons';
import {useAuth} from '../../context/AuthContext';
import {useNavigate} from 'react-router-dom';
import PathConstants from '../../routes/PathConstants';

const AntHeader = Layout.Header;

const Header: React.FC = () => {
  const {isAuthenticated, logout} = useAuth();
  const navigate = useNavigate();

  const menuItems = [
    {
      key: 'setGoal',
      label: 'Set Goal',
      icon: <SettingOutlined />,
      onClick: () => {
        navigate(PathConstants.GOAL);
      },
    },
    {
      key: 'logout',
      label: 'Logout',
      icon: <LogoutOutlined />,
      onClick: logout,
    },
  ];

  return (
    <AntHeader
      style={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        background: 'white',
        padding: '0 3vw 0 0',
        height: '8vh',
        minHeight: '30px',
      }}>
      <Logo />
      {isAuthenticated && (
        <Dropdown
          menu={{items: menuItems}}
          trigger={['hover']}
          overlayStyle={{minWidth: '10vw'}}>
          <div
            style={{
              cursor: 'pointer',
              marginLeft: 'auto',
              display: 'flex',
              alignItems: 'center',
              height: '100%',
            }}>
            <Avatar
              style={{
                backgroundColor: '#1677ff',
              }}
              icon={<UserOutlined />}
              size="large"
            />
          </div>
        </Dropdown>
      )}
    </AntHeader>
  );
};

export default Header;
