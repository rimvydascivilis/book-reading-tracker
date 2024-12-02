import React, {useState, useEffect} from 'react';
import {Modal, Input, message, Typography} from 'antd';
import ListSelector from './ListSelector';
import ListDetails from './ListDetails';
import api from '../../../api/api';
import {List} from '../../../types/listTypes';

const ListsPage: React.FC = () => {
  const [lists, setLists] = useState<List[]>([]);
  const [selectedList, setSelectedList] = useState<number | null>(null);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [newListTitle, setNewListTitle] = useState('');

  const fetchLists = async () => {
    try {
      const response = await api.get<List[]>('/lists');
      setLists(response.data);
    } catch (error) {
      console.error('Error fetching lists:', error);
      message.error('Failed to fetch lists.');
    }
  };

  const createList = async () => {
    if (!newListTitle.trim()) {
      message.error('Title cannot be empty.');
      return;
    }
    try {
      await api.post('/list', {title: newListTitle});
      message.success('List created successfully!');
      fetchLists();
      setIsModalVisible(false);
      setNewListTitle('');
    } catch (error) {
      console.error('Error creating list:', error);
      message.error('Failed to create list.');
    }
  };

  useEffect(() => {
    fetchLists();
  }, []);

  return (
    <div>
      <ListSelector
        lists={lists}
        selectedListId={selectedList}
        onSelect={listId => setSelectedList(listId)}
        onAddList={() => setIsModalVisible(true)}
      />
      {selectedList ? (
        <ListDetails id={selectedList} />
      ) : (
        <Typography.Text>Select a list to see its items.</Typography.Text>
      )}
      <Modal
        title="Create New List"
        open={isModalVisible}
        onOk={createList}
        onCancel={() => setIsModalVisible(false)}>
        <Input
          placeholder="Enter list title"
          value={newListTitle}
          onChange={e => setNewListTitle(e.target.value)}
        />
      </Modal>
    </div>
  );
};

export default ListsPage;
