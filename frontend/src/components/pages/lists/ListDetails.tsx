import React, {useEffect, useState} from 'react';
import {
  List,
  message,
  Typography,
  Select,
  Button,
  Form,
  Spin,
  Card,
  Row,
  Col,
} from 'antd';
import {DeleteOutlined} from '@ant-design/icons';
import api from '../../../api/api';
import {ListItem} from '../../../types/listTypes';
import {IAxiosError} from '../../../types/errorTypes';

const {Option} = Select;

interface ListDetailsProps {
  id: number;
}

interface IBook {
  id: number;
  title: string;
}

const ListDetails: React.FC<ListDetailsProps> = ({id}) => {
  const [items, setItems] = useState<ListItem[]>([]);
  const [bookSuggestions, setBookSuggestions] = useState<IBook[]>([]);
  const [selectedBook, setSelectedBook] = useState<number | null>(null);
  const [loading, setLoading] = useState(false);
  const [adding, setAdding] = useState(false);

  const fetchListItems = async (listId: number) => {
    setLoading(true);
    try {
      const response = await api.get<{list_items: ListItem[]}>('/list', {
        params: {list_id: listId},
      });
      setItems(response.data.list_items);
    } catch (error) {
      console.error('Error fetching list items:', error);
      message.error('Failed to fetch list items.');
    } finally {
      setLoading(false);
    }
  };

  const fetchBookSuggestions = async (title: string) => {
    if (!title) {
      setBookSuggestions([]);
      return;
    }

    setLoading(true);
    try {
      const response = await api.get(`/books/search`, {
        params: {limit: 10, title},
      });

      const books = response.data as IBook[];
      setBookSuggestions(books);
    } catch (error) {
      const axiosError = error as IAxiosError;
      message.error(
        'Failed to fetch book suggestions: ' +
          axiosError.response?.data.message,
      );
      setBookSuggestions([]);
    } finally {
      setLoading(false);
    }
  };

  const handleAddBook = async () => {
    if (!selectedBook) {
      message.warning('Please select a book!');
      return;
    }

    setAdding(true);
    try {
      await api.post('/list/item', {list_id: id, book_id: selectedBook});
      message.success('Book added to the list!');
      setSelectedBook(null);
      fetchListItems(id);
    } catch (error) {
      const axiosError = error as IAxiosError;
      if (axiosError.response?.status === 409) {
        message.warning('Book is already in the list.');
      } else {
        console.error('Error adding book to list:', error);
        message.error('Failed to add book to list.');
      }
    } finally {
      setAdding(false);
    }
  };

  const handleDeleteBook = async (itemId: number) => {
    try {
      await api.delete(`/list/${id}/item/${itemId}`);
      message.success('Book removed from the list!');
      fetchListItems(id);
    } catch (error) {
      console.error('Error deleting book from list:', error);
      message.error('Failed to remove book from list.');
    }
  };

  useEffect(() => {
    fetchListItems(id);
  }, [id]);

  return (
    <div style={{padding: '24px', maxWidth: '1200px', margin: 'auto'}}>
      <Row gutter={24}>
        <Col xs={24} md={16}>
          <Card
            title="Books in List"
            bordered={false}
            style={{
              borderRadius: '8px',
              boxShadow: '0 2px 8px rgba(0, 0, 0, 0.1)',
            }}>
            <List
              itemLayout="horizontal"
              dataSource={items}
              loading={loading}
              renderItem={item => (
                <List.Item
                  style={{
                    padding: '12px 16px',
                    borderRadius: '8px',
                    background: '#f9f9f9',
                    marginBottom: '8px',
                  }}
                  actions={[
                    <Button
                      key={item.id}
                      danger
                      icon={<DeleteOutlined />}
                      onClick={() => handleDeleteBook(item.id)}
                    />,
                  ]}>
                  <Typography.Text>{item.book_name}</Typography.Text>
                </List.Item>
              )}
            />
          </Card>
        </Col>

        <Col xs={24} md={8}>
          <Card
            title="Add a Book"
            bordered={false}
            style={{
              borderRadius: '8px',
              boxShadow: '0 2px 8px rgba(0, 0, 0, 0.1)',
            }}>
            <Form layout="vertical">
              <Form.Item label="Search for a Book" required>
                <Select
                  showSearch
                  placeholder="Type to search books"
                  onSearch={fetchBookSuggestions}
                  onChange={value => setSelectedBook(value)}
                  value={selectedBook}
                  notFoundContent={
                    loading ? <Spin size="small" /> : 'No results found'
                  }
                  filterOption={false}
                  style={{width: '100%'}}>
                  {bookSuggestions.map(book => (
                    <Option key={book.id} value={book.id}>
                      {book.title}
                    </Option>
                  ))}
                </Select>
              </Form.Item>
              <Form.Item>
                <Button
                  type="primary"
                  onClick={handleAddBook}
                  loading={adding}
                  disabled={!selectedBook}
                  block
                  style={{
                    borderRadius: '6px',
                    background: '#1890ff',
                    borderColor: '#1890ff',
                  }}>
                  Add Book
                </Button>
              </Form.Item>
            </Form>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default ListDetails;
