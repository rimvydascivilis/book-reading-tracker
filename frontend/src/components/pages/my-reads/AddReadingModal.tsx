import React, {useState} from 'react';
import {Form, Input, InputNumber, Modal, Select, message} from 'antd';
import {RuleObject} from 'antd/lib/form';
import {IReadingFormValues} from '../../../types/readingTypes';
import {IBook} from '../../../types/bookTypes';
import {IAxiosError} from '../../../types/errorTypes';
import api from '../../../api/api';

const {Option} = Select;

interface AddReadingModalProps {
  open: boolean;
  onAddReading: (newReading: IReadingFormValues) => void;
  onCancel: () => void;
}

const AddReadingModal: React.FC<AddReadingModalProps> = ({
  open,
  onAddReading,
  onCancel,
}) => {
  const [form] = Form.useForm();
  const [bookSuggestions, setBookSuggestions] = useState<IBook[]>([]);
  const [loading, setLoading] = useState(false);

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
      if (axiosError.response?.status === 404) {
        setBookSuggestions([]);
        return;
      }

      message.error(
        'failed to fetch book suggestions: ' +
          axiosError.response?.data.message || 'unknown error',
      );
      setBookSuggestions([]);
    } finally {
      setLoading(false);
    }
  };

  const handleAddReading = (values: IReadingFormValues) => {
    values.book_title =
      bookSuggestions.find(book => book.id === values.book_id)?.title || '';
    onAddReading(values);
    form.resetFields();
  };

  const validateUrl = (_: RuleObject, value: string) =>
    value && !/^https?:\/\/[^\s$.?#].[^\s]*$/i.test(value)
      ? Promise.reject(
          'Please enter a valid URL starting with http:// or https://',
        )
      : Promise.resolve();

  return (
    <Modal
      title="Add New Reading"
      open={open}
      onCancel={onCancel}
      onOk={form.submit}
      destroyOnClose
      okText="Add Reading"
      cancelText="Cancel">
      <Form form={form} onFinish={handleAddReading} layout="vertical">
        <Form.Item
          label="Book Title"
          name="book_id"
          rules={[{required: true, message: 'Please select a book!'}]}>
          <Select
            showSearch
            placeholder="Search for a book"
            onSearch={fetchBookSuggestions}
            loading={loading}
            notFoundContent={loading ? 'Loading...' : 'No results found'}
            filterOption={false}
            style={{width: '100%'}}>
            {bookSuggestions.map(book => (
              <Option key={book.id} value={book.id}>
                {book.title}
              </Option>
            ))}
          </Select>
        </Form.Item>

        <Form.Item
          label="Total Pages"
          name="total_pages"
          rules={[{required: true, message: 'Please enter the total pages!'}]}>
          <InputNumber
            min={1}
            placeholder="Enter total pages"
            style={{width: '100%'}}
          />
        </Form.Item>

        <Form.Item
          label="Link (Optional)"
          name="link"
          rules={[{validator: validateUrl}]}>
          <Input placeholder="Enter a link to the book (optional)" />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default AddReadingModal;
