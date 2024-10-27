import React from "react";
import {Layout, Menu} from "antd";


const {Sider} = Layout;

const Sidebar: React.FC = () => {
  return (
    <Sider collapsible theme="light">
      <Menu theme="light" mode="inline" items={[]} />
    </Sider>
  );
};

export default Sidebar;
