import { Layout, Menu, theme, MenuProps } from 'antd';
import { Outlet, useNavigate } from "react-router-dom";
import PathConstants from "../../routes/PathConstants";
import React from 'react';

const { Header, Content, Footer } = Layout;

const App: React.FC = () => {
  const navigate = useNavigate();

  const navItems: MenuProps["items"] = [
    {
      label: "Home",
      key: PathConstants.HOME,
    },
  ];

  const handleMenuClick: MenuProps['onClick'] = ({ key }) => {
    if (key) {
      navigate(key);
    }
  };

  const { token: { colorBgContainer, borderRadiusLG } } = theme.useToken();

  return (
    <Layout style={{minHeight:"100vh"}}>
      <Header style={{ display: 'flex', alignItems: 'center' }}>
        <div className="demo-logo" />
        <Menu
          theme="dark"
          mode="horizontal"
          style={{ flex: 1, minWidth: 0 }}
          defaultSelectedKeys={["/"]}
          items={navItems}
          onClick={handleMenuClick}
        />
      </Header>
      <Content style={{ padding: '0 48px' }}>
        <div
          style={{
            background: colorBgContainer,
            minHeight: 280,
            padding: 24,
            borderRadius: borderRadiusLG,
          }}
        >
          <Outlet />
        </div>
      </Content>
      <Footer style={{ textAlign: 'center' }}>
        Book reading tracker ©{new Date().getFullYear()} Created by Rimvydas
      </Footer>
    </Layout>
  );
};

export default App;