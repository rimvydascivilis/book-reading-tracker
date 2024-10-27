import React, {ReactNode, Suspense} from "react";
import {Outlet} from "react-router-dom";
import Sidebar from "./Sidebar";
import Loading from "../common/Loading";
import Logo from "./Logo";
import {Layout, theme, Button} from "antd";
import {LogoutOutlined} from "@ant-design/icons";
import {useAuth} from "../../context/AuthContext";

const {Content, Footer, Header} = Layout;

interface AppLayoutProps {
  children?: ReactNode;
}

const AppLayout: React.FC<AppLayoutProps> = ({children}) => {
  const {
    token: {colorBgContainer, borderRadiusLG},
  } = theme.useToken();

  const {isAuthenticated, logout} = useAuth();

  return (
    <Layout style={{minHeight: "100vh"}}>
      <Header
        style={{
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          background: "white",
          padding: "0 24px 0 0",
        }}>
        <Logo />
        {isAuthenticated && (
          <Button onClick={logout} type="primary" style={{marginLeft: "auto"}}>
            <LogoutOutlined />
            Logout
          </Button>
        )}
      </Header>
      <Layout>
        <Sidebar />
        <Layout>
          <Content
            style={{
              margin: "24px",
              textAlign: "center",
              background: colorBgContainer,
              borderRadius: borderRadiusLG,
            }}>
            <Suspense fallback={<Loading />}>{children ?? <Outlet />}</Suspense>
          </Content>
          <Footer style={{textAlign: "center"}}>
            Book tracker ©{new Date().getFullYear()} Created by Rimvydas
          </Footer>
        </Layout>
      </Layout>
    </Layout>
  );
};

export default AppLayout;
