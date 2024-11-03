import React, { useState } from 'react';
import { List, Rate, Button, Modal } from 'antd';
import { IBook } from '../../../types/bookTypes';
import { DeleteOutlined, EditOutlined } from '@ant-design/icons';

interface BookItemProps {
  book: IBook;
  onDelete: (id: number) => void;
  onRatingChange: (id: number, rating: number) => void;
  onTitleChange: (id: number, title: string) => void;
}

const BookItem: React.FC<BookItemProps> = ({ book, onDelete, onRatingChange, onTitleChange }) => {
  const [isEditing, setIsEditing] = useState<boolean>(false);
  const [title, setTitle] = useState<string>(book.title); // Local state for the title
  const [isHovered, setIsHovered] = useState<boolean>(false);

  const showDeleteConfirm = () => {
    Modal.confirm({
      title: 'Are you sure you want to delete this book?',
      content: `Book title: ${book.title}`,
      onOk() {
        onDelete(book.id);
      },
      onCancel() {},
    });
  };

  const handleEditToggle = () => {
    if (isEditing) {
      // When exiting edit mode, update the title
      onTitleChange(book.id, title);
    }
    setIsEditing(!isEditing);
  };

  return (
    <List.Item style={{ display: 'flex', justifyContent: 'space-between' }}>
      <div
        style={{ position: 'relative', display: 'flex', alignItems: 'center' }}
        onMouseEnter={() => setIsHovered(true)}
        onMouseLeave={() => setIsHovered(false)}
      >
        {isEditing ? (
          <input
            value={title}
            onChange={(e) => setTitle(e.target.value)} // Update local title state
            onBlur={handleEditToggle}
            style={{ border: '1px solid #d9d9d9', borderRadius: '4px', fontSize: '16px', marginRight: '10px' }}
            autoFocus
          />
        ) : (
          <>
            <span style={{ fontSize: '16px', marginRight: '5px' }}>{book.title}</span>
            {isHovered && (
              <EditOutlined
                style={{
                  cursor: 'pointer',
                  marginLeft: '5px',
                }}
                onClick={handleEditToggle}
              />
            )}
          </>
        )}
      </div>
      <div style={{ display: 'flex', alignItems: 'center' }}>
        <Rate
          allowHalf
          value={book.rating || 0}
          onChange={(value) => onRatingChange(book.id, value)}
        />
        <Button
          style={{ marginLeft: '10px' }}
          type="link"
          danger
          onClick={showDeleteConfirm}
          icon={<DeleteOutlined />}
        />
      </div>
    </List.Item>
  );
};

export default BookItem;
