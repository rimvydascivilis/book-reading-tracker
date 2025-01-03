import React, {useEffect, useState} from 'react';
import {useNavigate, useLocation} from 'react-router-dom';
import PathConstants from '../../routes/PathConstants';
import {Layout, Menu} from 'antd';
import {
  BankOutlined as LibraryOutlined,
  BookOutlined,
  UnorderedListOutlined,
  EditOutlined as NotesOutlined,
  BarChartOutlined,
} from '@ant-design/icons';

const {Sider} = Layout;

const items = [
  {
    key: '1',
    icon: <BookOutlined />,
    label: 'My Reads',
    path: PathConstants.MY_READS,
  },
  {
    key: '2',
    icon: <UnorderedListOutlined />,
    label: 'Lists',
    path: PathConstants.LISTS,
  },
  {
    key: '3',
    icon: <NotesOutlined />,
    label: 'Notes',
    path: PathConstants.NOTES,
  },
  {
    key: '4',
    icon: <LibraryOutlined />,
    label: 'Library',
    path: PathConstants.LIBRARY,
  },
  {
    key: '5',
    icon: <BarChartOutlined />,
    label: 'Stats',
    path: PathConstants.STATS,
  },
];

const Sidebar: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const [selectedKey, setSelectedKey] = useState<string | null>(null);

  useEffect(() => {
    const currentItem = items.find(item => location.pathname === item.path);
    setSelectedKey(currentItem ? currentItem.key : null);
  }, [location.pathname]);

  const handleClick = (e: {key: string}) => {
    const clickedItem = items.find(item => item.key === e.key);
    if (clickedItem) {
      navigate(clickedItem.path);
    }
  };

  return (
    <Sider collapsible theme="light">
      <Menu
        theme="light"
        mode="inline"
        items={items}
        selectedKeys={selectedKey ? [selectedKey] : undefined}
        onClick={handleClick}
      />
    </Sider>
  );
};

export default Sidebar;
