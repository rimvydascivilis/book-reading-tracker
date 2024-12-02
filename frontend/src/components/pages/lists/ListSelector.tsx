import React from 'react';
import {Button} from 'antd';

interface ListSelectorProps {
  lists: {id: number; title: string}[];
  selectedListId: number | null;
  onSelect: (listId: number) => void;
  onAddList: () => void;
}

const ListSelector: React.FC<ListSelectorProps> = ({
  lists,
  selectedListId,
  onSelect,
  onAddList,
}) => {
  return (
    <div
      style={{
        display: 'flex',
        gap: '16px',
        padding: '16px',
        backgroundColor: '#f9f9f9',
        borderRadius: '8px',
        boxShadow: '0 4px 8px rgba(0, 0, 0, 0.1)',
      }}>
      {lists.map(list => (
        <Button
          key={list.id}
          type={selectedListId === list.id ? 'primary' : 'default'}
          onClick={() => onSelect(list.id)}
          style={{
            flex: 1,
            textAlign: 'center',
            padding: '16px',
            borderRadius: '8px',
            fontSize: '16px',
            fontWeight: 'bold',
          }}>
          {list.title}
        </Button>
      ))}
      <Button
        type="dashed"
        onClick={onAddList}
        style={{
          flex: 1,
          textAlign: 'center',
          padding: '16px',
          borderRadius: '8px',
          fontSize: '16px',
        }}>
        +
      </Button>
    </div>
  );
};

export default ListSelector;
