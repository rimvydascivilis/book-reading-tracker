import React, {useEffect, useState} from 'react';
import {Button, Card, List, message, Modal} from 'antd';
import api from '../../../api/api';
import {IAxiosError} from '../../../types/errorTypes';
import {INote} from '../../../types/note';

interface NotesListProps {
  bookId: string;
  refresh: boolean;
  onDelete: (noteId: string) => void;
}

const NotesList: React.FC<NotesListProps> = ({bookId, refresh, onDelete}) => {
  const [loading, setLoading] = useState(false);
  const [notes, setNotes] = useState<INote[]>([]);

  const fetchNotes = async (bookId: string) => {
    try {
      setLoading(true);
      const response = await api.get(`/notes/${bookId}`);
      setNotes(response.data);
    } catch (error) {
      const axiosError = error as IAxiosError;
      message.error(
        'Failed to fetch notes: ' +
          (axiosError.response?.data.message || 'Unknown error'),
      );
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchNotes(bookId);
  }, [bookId, refresh]);

  const confirmDelete = (noteId: string) => {
    Modal.confirm({
      title: 'Are you sure you want to delete this note?',
      content: 'Once deleted, this action cannot be undone.',
      onOk: () => {
        onDelete(noteId);
      },
      okText: 'Yes',
      cancelText: 'No',
    });
  };

  return (
    <List
      loading={loading}
      bordered
      dataSource={notes}
      locale={{emptyText: 'No notes available for this book.'}}
      renderItem={note => (
        <List.Item
          actions={[
            <Button danger onClick={() => confirmDelete(note.id)} key={note.id}>
              Delete
            </Button>,
          ]}>
          <Card title={`Page ${note.page_number}`} style={{width: '100%'}}>
            {note.content}
          </Card>
        </List.Item>
      )}
    />
  );
};

export default NotesList;
