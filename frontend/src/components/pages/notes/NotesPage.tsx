import React, {useState} from 'react';
import {Button, Select, Space, message} from 'antd';
import {IAxiosError} from '../../../types/errorTypes';
import api from '../../../api/api';
import NotesList from './NotesList';
import NotesModal from './NotesModal';

const {Option} = Select;

interface Book {
  id: string;
  title: string;
}

const NotesPage: React.FC = () => {
  const [bookSuggestions, setBookSuggestions] = useState<Book[]>([]);
  const [selectedBookId, setSelectedBookId] = useState<string | null>(null);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [loading, setLoading] = useState(false);
  const [refreshList, setRefreshList] = useState(false);

  const fetchBookSuggestions = async (title: string) => {
    try {
      setLoading(true);
      const response = await api.get(`/books/search`, {
        params: {limit: 10, title},
      });
      setBookSuggestions(response.data);
    } catch (error) {
      const axiosError = error as IAxiosError;
      if (axiosError.response?.status === 404) {
        setBookSuggestions([]);
        return;
      } else {
        message.error(
          'Failed to fetch book suggestions: ' +
            axiosError.response?.data.message,
        );
      }
    } finally {
      setLoading(false);
    }
  };

  const handleCreateNote = async (note: {
    page_number: number;
    content: string;
  }) => {
    if (!selectedBookId) return;
    try {
      await api.post(`/notes/${selectedBookId}`, note);
      message.success('Note created successfully!');
      setIsModalVisible(false);
      setRefreshList(!refreshList);
    } catch (error) {
      const axiosError = error as IAxiosError;
      message.error(
        'Failed to create note: ' + axiosError.response?.data.message,
      );
    }
  };

  const handleDeleteNote = async (noteId: string) => {
    try {
      await api.delete(`/notes/${noteId}`);
      message.success('Note deleted successfully!');
      setRefreshList(!refreshList);
    } catch (error) {
      const axiosError = error as IAxiosError;
      message.error(
        'Failed to delete note: ' + axiosError.response?.data.message,
      );
    }
  };

  return (
    <div style={{padding: '24px'}}>
      <Space direction="vertical" style={{width: '100%'}}>
        <Select
          showSearch
          placeholder="Search for a book"
          onSearch={fetchBookSuggestions}
          loading={loading}
          notFoundContent={loading ? 'Loading...' : 'No results found'}
          filterOption={false}
          style={{width: '100%'}}
          onChange={value => setSelectedBookId(value)}>
          {bookSuggestions.map(book => (
            <Option key={book.id} value={book.id}>
              {book.title}
            </Option>
          ))}
        </Select>

        <Button
          type="primary"
          onClick={() => setIsModalVisible(true)}
          disabled={!selectedBookId}>
          Add Note
        </Button>

        <NotesList
          bookId={selectedBookId || '-1'}
          refresh={refreshList}
          onDelete={handleDeleteNote}
        />
      </Space>

      <NotesModal
        visible={isModalVisible}
        onClose={() => setIsModalVisible(false)}
        onCreate={handleCreateNote}
      />
    </div>
  );
};

export default NotesPage;
