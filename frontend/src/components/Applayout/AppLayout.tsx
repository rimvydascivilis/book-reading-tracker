import React, {ReactNode, Suspense} from 'react';
import {Outlet} from 'react-router-dom';
import Sidebar from './Sidebar';
import Loading from '../common/Loading';
import {Layout, theme} from 'antd';
import Header from './Header';

const {Content, Footer} = Layout;

interface AppLayoutProps {
  children?: ReactNode;
}

const AppLayout: React.FC<AppLayoutProps> = ({children}) => {
  const {
    token: {colorBgContainer, borderRadiusLG},
  } = theme.useToken();

  return (
    <Layout style={{height: '100vh'}}>
      <Header />
      <Layout>
        <Sidebar />
        <Layout>
          <Content
            style={{
              margin: '24px',
              textAlign: 'center',
              background: colorBgContainer,
              borderRadius: borderRadiusLG,
              overflowY: 'auto',
            }}>
            <Suspense fallback={<Loading />}>{children ?? <Outlet />}</Suspense>
          </Content>
          <Footer style={{
            background: colorBgContainer,
            maxHeight: '50px',
            height: '4vh',
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
          }}>
            Book tracker ©{new Date().getFullYear()} Created by Rimvydas
          </Footer>
        </Layout>
      </Layout>
    </Layout>
  );
};

export default AppLayout;
